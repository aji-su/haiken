package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type Visibility = string

const (
	Public   Visibility = "public"
	Unlisted Visibility = "unlisted"
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
		Path:   "/api/v1/accounts/verify_credentials",
	}
	params := url.Values{}
	params.Set("access_token", r.token)
	u.RawQuery = params.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, errors.Wrap(err, "verify_credentials err")
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

func (r *RestClient) Post(status, inReplyToID, spoiler string, visibility Visibility) ([]byte, error) {
	u := url.URL{
		Scheme: r.scheme,
		Host:   r.host,
		Path:   "/api/v1/statuses",
	}
	params := url.Values{}
	params.Set("access_token", r.token)
	params.Set("status", status)
	params.Set("in_reply_to_id", inReplyToID)
	params.Set("spoiler_text", spoiler)
	params.Set("visibility", visibility)

	resp, err := http.Post(
		u.String(),
		"application/x-www-form-urlencoded",
		strings.NewReader(params.Encode()),
	)
	if err != nil {
		return nil, errors.Wrap(err, "post err")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "io err")
	}
	log.Printf("resp: %s", string(body))
	return body, err
}

func (r *RestClient) Follow(id string, follow bool) error {
	var path string
	if follow {
		path = fmt.Sprintf("/api/v1/accounts/%s/follow", id)
	} else {
		path = fmt.Sprintf("/api/v1/accounts/%s/unfollow", id)
	}
	u := url.URL{
		Scheme: r.scheme,
		Host:   r.host,
		Path:   path,
	}
	params := url.Values{}
	params.Set("access_token", r.token)

	resp, err := http.Post(
		u.String(),
		"application/x-www-form-urlencoded",
		strings.NewReader(params.Encode()),
	)
	if err != nil {
		return errors.Wrap(err, "post err")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "io err")
	}
	log.Printf("resp: %s", string(body))
	return err
}
