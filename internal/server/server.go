package server

import (
	"encoding/json"
	"fmt"
	"io"
	log "log/slog"
	"net/http"

	"github.com/modaniru/streamer-notifier-telegram/internal/entity"
	"github.com/modaniru/streamer-notifier-telegram/internal/service"
	"github.com/modaniru/streamer-notifier-telegram/internal/telegram"
)

type server struct {
	mux           *http.ServeMux
	telegramBot *telegram.TelegramBot
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
	if r.Header.Get("Twitch-Eventsub-Message-Type") == "webhook_callback_verification"{
		var v entity.Verify
		b, err := io.ReadAll(r.Body)
		if err != nil{
			log.Error(err.Error())
		}
		err = json.Unmarshal(b, &v)
		if err != nil{
			log.Error(err.Error())
		}
		w.Write([]byte(v.Challenge))
	} else {
		b, _ := io.ReadAll(r.Body)
		var response entity.StreamOnlineNotification
		err := json.Unmarshal(b, &response)
		if err != nil{
			log.Error("unmarshal error", log.String("error", err.Error()))
			return
		}
		userList, err := s.followService.GetStreamerUsers(response.Event.BroadcasterUserId)
		if err != nil{
			log.Error("get users that followed to streamer error", log.String("error", err.Error()))
			return 
		}
		for _, id := range userList{
			s.telegramBot.SendNotification(response, int64(id))
		}
	}
	
}
