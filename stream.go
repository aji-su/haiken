package main

import (
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type Stream struct {
	schema, host, token string
	conn                *websocket.Conn
	handler             *Haiken
}

func NewStream(s, h, t string) (*Stream, error) {
	u := url.URL{
		Scheme: s,
		Host:   h,
		Path:   "/streaming",
	}
	log.Println("listening to: ", u.String())

	params := url.Values{}
	params.Set("i", t)
	u.RawQuery = params.Encode()

	c, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "dial err")
	}

	defer resp.Body.Close()
	rBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "io err")
	}
	log.Printf("response: %s, %s", resp.Status, rBody)

	return &Stream{
		schema: s,
		host:   h,
		token:  t,
		conn:   c,
	}, nil
}

func (m *Stream) Destroy() {
	m.conn.Close()
}

func (m *Stream) SetHandler(h *Haiken) {
	m.handler = h
}

func (s *Stream) subscribe(channel, streamID string) error {
	return s.conn.WriteJSON(map[string]interface{}{
		"type": "connect",
		"body": map[string]string{
			"channel": channel,
			"id":      streamID,
		},
	})
}

func (m *Stream) Stream(homeStreamID, mainStreamID string) error {
	if err := m.subscribe("homeTimeline", homeStreamID); err != nil {
		return errors.Wrap(err, "failed to subscribe")
	}
	if err := m.subscribe("main", mainStreamID); err != nil {
		return errors.Wrap(err, "failed to subscribe")
	}

	fatalErr := make(chan error)

	go func() {
		defer close(fatalErr)
		for {
			_, message, err := m.conn.ReadMessage()
			if err != nil {
				fatalErr <- errors.Wrap(err, "read err")
				return
			}
			log.Printf("received: %s", message)
			if err := m.handler.Handle(message); err != nil {
				fatalErr <- errors.Wrap(err, "review err")
				return
			}
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		select {
		case err := <-fatalErr:
			return err
		case <-interrupt:
			log.Println("interrupt")
			return nil
		}
	}
}
