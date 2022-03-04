package main

import (
	"log"
	"os"

	"github.com/aji-su/haiken/ikku-go"
	"github.com/google/uuid"
)

func main() {
	hscheme := os.Getenv("MISSKEY_HTTP_SCHEME")
	wsscheme := os.Getenv("MISSKEY_WS_SCHEME")
	hhost := os.Getenv("MISSKEY_HTTP_HOST")
	wshost := os.Getenv("MISSKEY_WS_HOST")
	token := os.Getenv("MISSKEY_ACCESS_TOKEN")

	parser, err := ikku.NewParser()
	if err != nil {
		log.Fatalf("parser err: %s", err)
	}
	defer parser.Destroy()
	reviewer := ikku.NewReviewer(parser, []int{5, 7, 5})

	stream, err := NewStream(wsscheme, wshost, token)
	if err != nil {
		log.Fatalf("NewStream err: %s", err)
	}
	defer stream.Destroy()

	rest := NewRestClient(hscheme, hhost, token)
	act, err := rest.VerifyCredentials()
	if err != nil {
		log.Fatalf("VerifyCredentials err: %s", err)
	}
	log.Printf("VerifyCredentials: %#v", act)

	var u uuid.UUID
	u, err = uuid.NewRandom()
	if err != nil {
		log.Fatalf("NewRandom err: %s", err)
	}
	homeStreamID := u.String()
	u, err = uuid.NewRandom()
	if err != nil {
		log.Fatalf("NewRandom err: %s", err)
	}
	mainStreamID := u.String()

	haiken := NewHaiken(reviewer, act, rest, homeStreamID, mainStreamID)
	stream.SetHandler(haiken)
	log.Fatal("FATAL stream err: ", stream.Stream(homeStreamID, mainStreamID))
}
