package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	log "log/slog"
	"net/http"
	"os"

	"github.com/modaniru/streamer-notifier-telegram/internal/entity"
	"github.com/modaniru/streamer-notifier-telegram/internal/service"
	"github.com/modaniru/streamer-notifier-telegram/internal/telegram"
)

type server struct {
	mux           *http.ServeMux
	telegramBot   *telegram.TelegramBot
	followService service.StreamerService
}

func NewServer(s *http.ServeMux, telegram *telegram.TelegramBot, followService service.StreamerService) *server {
	return &server{mux: s, telegramBot: telegram, followService: followService}
}

func (s *server) Start(port string, channel chan int) {
	s.mux.HandleFunc("/stream-online", s.StreamOnline)
	log.Info(fmt.Sprintf("server %s port", port))
	err := http.ListenAndServe(":"+port, s.mux)
	log.Error("serve http server error", log.String("error", err.Error()))
	channel <- 1
}

// TODO check request sender
func (s *server) StreamOnline(w http.ResponseWriter, r *http.Request) {
	if ok, err := isTwitch(r); !ok || err != nil {
		fmt.Println(ok, err.Error())
		if err != nil {
			http.Error(w, err.Error(), 404)
			return
		}
		http.Error(w, "no valid key", 404)
	}
	if r.Header.Get("Twitch-Eventsub-Message-Type") == "webhook_callback_verification" {
		log.Info("webhook verification")
		var v entity.Verify
		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error(err.Error())
		}
		err = json.Unmarshal(b, &v)
		if err != nil {
			log.Error(err.Error())
		}
		w.Write([]byte(v.Challenge))
	} else {
		b, _ := io.ReadAll(r.Body)
		var response entity.StreamOnlineNotification
		err := json.Unmarshal(b, &response)
		if err != nil {
			log.Error("unmarshal error", log.String("error", err.Error()))
			return
		}
		log.Info("streamer online " + response.Event.BroadcasterUserLogin)
		userList, err := s.followService.GetStreamerUsers(response.Event.BroadcasterUserId)
		if err != nil {
			log.Error("get users that followed to streamer error", log.String("error", err.Error()))
			return
		}
		log.Info("userList", log.Any("arr", userList))
		for _, id := range userList {
			s.telegramBot.SendNotification(response, int64(id))
		}
	}

}

func isTwitch(r *http.Request) (bool, error) {
	body := r.Body
	b, err := io.ReadAll(body)
	if err != nil {
		return false, err
	}
	message := r.Header.Get("TWITCH_MESSAGE_ID") + r.Header.Get("TWITCH_MESSAGE_TIMESTAMP") + string(b)
	fmt.Println(message)
	sig := hmac.New(sha256.New, []byte(os.Getenv("SECRET")))
	sig.Write([]byte(message))
	h := "sha256=" + hex.EncodeToString(sig.Sum(nil))
	return hmac.Equal([]byte(h), []byte(r.Header.Get("Twitch-Eventsub-Message-Signature"))), nil
}
