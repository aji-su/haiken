package main

import (
	"log"
	"os"

	"github.com/aji-su/haiken/ikku-go"
)

func main() {
	hscheme := os.Getenv("MASTODON_HTTP_SCHEME")
	wsscheme := os.Getenv("MASTODON_WS_SCHEME")
	hhost := os.Getenv("MASTODON_HTTP_HOST")
	wshost := os.Getenv("MASTODON_WS_HOST")
	token := os.Getenv("MASTODON_ACCESS_TOKEN")
	subscriptions := os.Getenv("MASTODON_WS_SUBSCRIPTIONS")

	parser, err := ikku.NewParser()
	if err != nil {
		log.Fatalf("parser err: %s", err)
	}
	defer parser.Destroy()
	reviewer := ikku.NewReviewer(parser, []int{5, 7, 5})

	stream, err := NewStream(wsscheme, wshost, token, subscriptions)
	if err != nil {
		log.Fatalf("NewStream err: %s", err)
	}
	defer stream.Destroy()

	rest := NewRestClient(hscheme, hhost, token)
	act, err := rest.VerifyCredentials()
	if err != nil {
		log.Fatalf("VerifyCredentials err: %s", err)
	}
	haiken := NewHaiken(reviewer, act, rest, os.Getenv("MASTODON_ALLOWED_TAGS"))
	stream.SetHandler(haiken)

	if err != nil {
		log.Fatalf("receiver err: %s", err)
	}
	log.Fatal("stream err: ", stream.Stream())
}
