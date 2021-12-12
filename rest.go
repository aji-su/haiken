package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type RestClient struct {
	scheme, host, token string
}

func NewRestClient(s, h, t string) *RestClient {
	return &RestClient{
		scheme: s,
		host:   h,
		token:  t,
	}
}

func (r *RestClient) VerifyCredentials() (*Account, error) {
	u := url.URL{
		Scheme: r.scheme,
		Host:   r.host,
		Path:   "/api/i",
	}
	resp, err := http.Post(
		u.String(),
		"application/json",
		bytes.NewBuffer([]byte(fmt.Sprintf(`{"i":"%s"}`, r.token))))
	if err != nil {
		return nil, errors.Wrap(err, "get 'i' err")
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "io err")
	}
	var act *Account
	if err := json.Unmarshal(b, &act); err != nil {
		return nil, errors.Wrap(err, "act unmarshal err")
	}
	return act, nil
}

func (r *RestClient) Post(text string, replyID *string, cw *string, visibility string, localOnly bool) ([]byte, error) {
	u := url.URL{
		Scheme: r.scheme,
		Host:   r.host,
		Path:   "/api/notes/create",
	}
	payload := map[string]interface{}{
		"i":          r.token,
		"localOnly":  localOnly,
		"text":       text,
		"visibility": visibility,
		"replyId":    replyID,
		"cw":         cw,
	}
	p, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(
		u.String(),
		"application/json",
		bytes.NewReader(p),
	)
	if err != nil {
		return nil, errors.Wrap(err, "post err")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "io err")
	}
	log.Printf("Response to post: %s", string(body))
	return body, err
}

func (r *RestClient) Follow(id string, follow bool) error {
	var path string
	if follow {
		path = "/api/following/create"
	} else {
		path = "/api/following/delete"
	}
	u := url.URL{
		Scheme: r.scheme,
		Host:   r.host,
		Path:   path,
	}
	payload := map[string]interface{}{
		"i":      r.token,
		"userId": id,
	}

	p, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		u.String(),
		"application/json",
		bytes.NewReader(p),
	)
	if err != nil {
		return errors.Wrap(err, "post err")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "io err")
	}
	log.Printf("Response to follow: %s", string(body))
	return err
}
