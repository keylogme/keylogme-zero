package keylogme

import (
	"fmt"
	"log"
	"log/slog"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

type Sender struct {
	origin_endpoint string
	url_ws          string
	ws              *websocket.Conn
	max_retries     int64
	retry_duration  time.Duration
}

func MustGetNewSender(origin, apikey string) *Sender {
	if origin == "" {
		log.Fatal("Origin endpoint is empty string")
	}
	if apikey == "" {
		log.Fatal("ApiKey endpoint is empty string")
	}

	trimmedOrigin := strings.TrimPrefix(origin, "http")
	url_ws := fmt.Sprintf("ws%s?apikey=%s", trimmedOrigin, apikey)

	ws, err := websocket.Dial(url_ws, "", origin)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &Sender{
		origin_endpoint: origin,
		url_ws:          url_ws,
		ws:              ws,
		max_retries:     10000,
		retry_duration:  1 * time.Second,
	}
}

func (s *Sender) Send(p []byte) error {
	_, err := s.ws.Write(p)
	if err != nil {
		slog.Error(err.Error())
		err = s.reconnect()
		if err != nil {
			return err
		}
		// retry
		s.Send(p)
	}
	return nil
}

func (s *Sender) reconnect() error {
	for i := range s.max_retries {
		slog.Info("Waiting for reconnecting...")
		time.Sleep(s.retry_duration)
		ws_reconnect, err := websocket.Dial(s.url_ws, "", s.origin_endpoint)
		if err != nil {
			continue
		}
		slog.Info(fmt.Sprintf("Reconnected after %d retries\n", i+1))
		s.ws = ws_reconnect
		return nil
	}
	return fmt.Errorf("Maximum retries excedeed\n")
}

func (s *Sender) Close() error {
	return s.ws.Close()
}
