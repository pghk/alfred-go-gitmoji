package internal

import (
	"github.com/dnaeon/go-vcr/v2/recorder"
	"github.com/kinbiko/jsonassert"
	"log"
	"net/http"
	"testing"
)

func setup(name string) (*http.Client, *recorder.Recorder) {
	vcr, err := recorder.New("../test/_fixtures/" + name)
	if err != nil {
		log.Fatal(err)
	}
	return &http.Client{Transport: vcr}, vcr
}

func teardown(vcr *recorder.Recorder) {
	err := vcr.Stop()
	if err != nil {
		log.Fatal(err)
	}
}

func Test_getGitmojiList(t *testing.T) {
	client, vcr := setup("get_gitmoji")
	ja := jsonassert.New(t)

	ja.Assertf(getGitmojiList(client), `{"gitmojis": "<<PRESENCE>>"}`)

	teardown(vcr)
}
