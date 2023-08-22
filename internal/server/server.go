package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/modaniru/streamer-notifier-telegram/internal/entity"
	"github.com/modaniru/streamer-notifier-telegram/internal/telegram"
)

// TODO TLS
type server struct {
	mux           *http.ServeMux
	telegramBot *telegram.TelegramBot
}

func NewServer(s *http.ServeMux, telegram *telegram.TelegramBot) *server {
	return &server{mux: s, telegramBot: telegram}
}

func (s *server) Start(port string, channel chan int) {
	s.mux.HandleFunc("/stream-online", s.StreamOnline)
	log.Println("server start in " + port + " port")
	err := http.ListenAndServe(":"+port, s.mux)
	channel <- 1
	log.Fatal(err.Error())
}

func (s *server) StreamOnline(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Twitch-Eventsub-Message-Type") == "webhook_callback_verification"{
		var v entity.Verify
		b, err := io.ReadAll(r.Body)
		if err != nil{
			log.Println(err.Error())
		}
		err = json.Unmarshal(b, &v)
		if err != nil{
			log.Println(err.Error())
		}
		w.Write([]byte(v.Challenge))
	} else {
		b, _ := io.ReadAll(r.Body)
		var response entity.StreamOnlineNotification
		err := json.Unmarshal(b, &response)
		log.Println(response, err)
		s.telegramBot.SendNotification(response, 31312)
	}
	
}
