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

type Haiken struct {
	reviewer   *ikku.Reviewer
	account    *Account
	rest       *RestClient
	allowedTag string
}

func NewHaiken(r *ikku.Reviewer, a *Account, rc *RestClient, t string) *Haiken {
	return &Haiken{
		reviewer:   r,
		account:    a,
		rest:       rc,
		allowedTag: t,
	}
}

func (h *Haiken) Handle(message []byte) error {
	var m Message
	if err := json.Unmarshal(message, &m); err != nil {
		return err
	}
	switch m.Event {
	case "update":
		var s *Status
		if err := json.Unmarshal([]byte(m.Payload), &s); err != nil {
			return err
		}
		if s.Account.ID == h.account.ID {
			log.Printf("skip own post")
			return nil
		}
		if s.Reblog != nil {
			log.Printf("skip reblog")
			return nil
		}
		if s.Visibility != "public" && s.Visibility != "unlisted" {
			log.Printf("skip private post")
			return nil
		}
		for _, m := range s.Mentions {
			if m.ID == h.account.ID {
				return h.handleMention(s)
			}
		}
		if err := h.review(s, false); err != nil {
			return errors.Wrap(err, "review err")
		}
	case "notification":
		var p *Notification
		if err := json.Unmarshal([]byte(m.Payload), &p); err != nil {
			return err
		}
		if p.Type == "follow" {
			if err := h.rest.Follow(p.Account.ID, true); err != nil {
				return errors.Wrap(err, "follow err")
			}
		}
	case "delete":
	default:
		return errors.Errorf("unknown event: %s", string(message))
	}
	return nil
}

func (h *Haiken) review(s *Status, force bool) error {
	con := strip.StripTags(s.Content)
	nodes, songs, err := h.reviewer.Search(con)
	if err != nil {
		return errors.Wrap(err, "failed to review")
	}
	if force && len(songs) < 1 {
		if err := h.sendDetail(nodes, s.Account.Acct, s.ID); err != nil {
			return errors.Wrap(err, "sendDetail err")
		}
	}
	if len(songs) > 0 {
		resBody, err := h.sendReport(nodes, songs, s)
		if err != nil {
			return errors.Wrap(err, "sendReport err")
		}
		var resStat *Status
		if err := json.Unmarshal(resBody, &resStat); err != nil {
			return errors.Wrap(err, "resBody unmarshal err")
		}
		if err := h.sendDetail(nodes, "", resStat.ID); err != nil {
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

	var tags string
	for _, tag := range s.Tags {
		if tag.Name == h.allowedTag {
			tags = " #" + tag.Name
			break
		}
	}

	report := fmt.Sprintf("『%s』%s", strings.Join(sSongs, "』\n\n『"), tags)

	if s.SpoilerText != "" {
		return h.rest.Post(
			fmt.Sprintf("@%s\n%s", s.Account.Acct, report),
			s.ID,
			"俳句を発見しました！",
			s.Visibility)
	} else {
		return h.rest.Post(
			fmt.Sprintf("@%s 俳句を発見しました！\n%s", s.Account.Acct, report),
			s.ID,
			"",
			s.Visibility)
	}
}

func (h *Haiken) sendDetail(nodes []*ikku.Node, acct, id string) error {
	var ds []string
	for _, node := range nodes {
		ds = append(ds, fmt.Sprintf("[%s:%d]", node.Pronunciation(), node.PronunciationLength()))
	}
	details := strings.Join(ds, ",")
	if acct != "" {
		details = fmt.Sprintf("@%s %s", acct, details)
	}
	_, err := h.rest.Post(details, id, "俳句解析結果", "unlisted")
	return err
}

func (h *Haiken) handleMention(s *Status) error {
	if strings.Contains(s.Content, "俳句検出を停止してください") {
		return h.rest.Follow(s.Account.ID, false)
	}
	return h.review(s, true)
}
