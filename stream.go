package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type Subscription struct {
	Type   string `json:"type"`
	Stream string `json:"stream"`
	Tag    string `json:"tag"`
}

type Stream struct {
	schema, host, token string
	subscriptions       []*Subscription
	conn                *websocket.Conn
	handler             *Haiken
}

func NewStream(s, h, t, ss string) (*Stream, error) {
	u := url.URL{
		Scheme: s,
		Host:   h,
		Path:   "/api/v1/streaming",
	}
	log.Println("listening to: ", u.String())

	params := url.Values{}
	params.Set("access_token", t)
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

	var subs []*Subscription
	if err := json.Unmarshal([]byte(ss), &subs); err != nil {
		return nil, errors.Wrapf(err, "unmarshal err '%s'", ss)
	}

	return &Stream{
		schema:        s,
		subscriptions: subs,
		host:          h,
		token:         t,
		conn:          c,
	}, nil
}

func (m *Stream) Destroy() {
	m.conn.Close()
}

func (m *Stream) SetHandler(h *Haiken) {
	m.handler = h
}

func (s *Stream) subscribe() error {
	for _, subs := range s.subscriptions {
		log.Printf("subs: '%#v'", subs)
		if err := s.conn.WriteJSON(subs); err != nil {
			return errors.Wrap(err, "write err")
		}
	}
	return nil
}

func (m *Stream) Stream() error {
	if err := m.subscribe(); err != nil {
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
