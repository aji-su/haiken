package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aji-su/haiken/ikku-go"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/pkg/errors"
)

const (
	FindMessage   = "俳句を発見しました！"
	DetailMessage = "俳句解析結果"
	StopMessage   = "俳句検出を停止してください"
)

type Haiken struct {
	reviewer     *ikku.Reviewer
	account      *Account
	rest         *RestClient
	homeStreamID string
	mainStreamID string
}

func NewHaiken(r *ikku.Reviewer, a *Account, rc *RestClient, h, m string) *Haiken {
	return &Haiken{
		reviewer:     r,
		account:      a,
		rest:         rc,
		homeStreamID: h,
		mainStreamID: m,
	}
}

func (h *Haiken) Handle(message []byte) error {
	var m Message
	if err := json.Unmarshal(message, &m); err != nil {
		log.Printf("failed to unmarshal message: %s", message)
		return err
	}
	switch m.Body.ID {
	case h.homeStreamID:
		var s *Status
		if err := json.Unmarshal([]byte(m.Body.Body), &s); err != nil {
			return err
		}
		if s.Account.ID == h.account.ID {
			log.Printf("skip own post")
			return nil
		}
		if s.Renote != nil {
			log.Printf("skip reblog")
			return nil
		}
		for _, id := range s.Mentions {
			if id == h.account.ID {
				if err := h.handleMention(s); err != nil {
					log.Printf("failed to handle mention: %v, message=%#v", err, m)
				}
				return nil
			}
		}
		if s.Visibility != "public" && s.Visibility != "home" {
			log.Printf("skip private post")
			return nil
		}
		if err := h.review(s, false); err != nil {
			log.Printf("review err: %v, message=%#v", err, m)
			return nil
		}
	case h.mainStreamID:
		if m.Body.Type != "followed" {
			log.Printf("Ignoring unknown type: %s", m.Body.Type)
			return nil
		}
		var a *Account
		if err := json.Unmarshal([]byte(m.Body.Body), &a); err != nil {
			log.Printf("failed to unmarshal account: %s, message=%#v", m.Body.Body, m)
			return err
		}
		log.Printf("Followed by the account: %#v", a)
		if err := h.rest.Follow(a.ID, true); err != nil {
			log.Printf("failed to follow account: %s, message=%#v", m.Body.Body, m)
			return nil
		}
	default:
		log.Printf("Unknown streamID: %s, message=%#v", m.Body.ID, m)
		return nil
	}
	return nil
}

func (h *Haiken) review(s *Status, force bool) error {
	log.Printf("reviewing: %#v", s)
	con := strip.StripTags(s.Text)
	nodes, songs, err := h.reviewer.Search(con)
	if err != nil {
		return errors.Wrap(err, "failed to review")
	}
	if force && len(songs) < 1 {
		if err := h.sendDetail(nodes, s.Account.Username, s.ID, s.Visibility, s.LocalOnly); err != nil {
			return errors.Wrap(err, "sendDetail err")
		}
	}
	if len(songs) > 0 {
		resBody, err := h.sendReport(nodes, songs, s)
		if err != nil {
			return errors.Wrap(err, "sendReport err")
		}
		var resp struct {
			CreatedNote *Status `json:"createdNote"`
		}
		if err := json.Unmarshal(resBody, &resp); err != nil {
			return errors.Wrap(err, "resBody unmarshal err")
		}
		log.Printf("result id: %v", resp.CreatedNote.ID)
		if err := h.sendDetail(nodes, s.Account.Username, s.ID, s.Visibility, s.LocalOnly); err != nil {
			return errors.Wrap(err, "sendDetail err")
		}
	}
	return nil
}

func (h *Haiken) sendReport(nodes []*ikku.Node, songs []*ikku.Song, s *Status) ([]byte, error) {
	var sSongs []string
	for _, song := range songs {
		var s []string
		for _, phrases := range song.Phrases() {
			var sPhrase []string
			for _, phrase := range phrases {
				sPhrase = append(sPhrase, phrase.Surface())
			}
			s = append(s, strings.Join(sPhrase, ""))
		}
		sSongs = append(sSongs, strings.Join(s, " "))
	}

	report := fmt.Sprintf("『%s』", strings.Join(sSongs, "』\n\n『"))

	if s.Cw != nil {
		return h.rest.Post(
			message(s.Account.Username, s.Account.Host, report),
			stringP(s.ID),
			stringP(FindMessage),
			s.Visibility,
			s.LocalOnly,
		)
	} else {
		return h.rest.Post(
			message(s.Account.Username, s.Account.Host, fmt.Sprintf("%s\n%s", FindMessage, report)),
			stringP(s.ID),
			nil,
			s.Visibility,
			s.LocalOnly,
		)
	}
}

func message(username string, host *string, message string) string {
	if host == nil {
		return fmt.Sprintf("@%s\n%s", username, message)
	}
	return fmt.Sprintf("@%s@%s\n%s", username, *host, message)
}

func (h *Haiken) sendDetail(nodes []*ikku.Node, username string, replyID string, visibility string, localOnly bool) error {
	var ds []string
	for _, node := range nodes {
		ds = append(ds, fmt.Sprintf("[%s:%d]", node.Pronunciation(), node.PronunciationLength()))
	}
	details := strings.Join(ds, ",")
	if username != "" {
		details = fmt.Sprintf("@%s %s", username, details)
	}
	_, err := h.rest.Post(details, stringP(replyID), stringP(DetailMessage), visibility, localOnly)
	return err
}

func (h *Haiken) handleMention(s *Status) error {
	if strings.Contains(s.Text, StopMessage) {
		return h.rest.Follow(s.Account.ID, false)
	}
	return h.review(s, true)
}

func stringP(s string) *string {
	return &s
}
