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

func Test_parseGitmoji(t *testing.T) {
	client, vcr := setup("get_gitmoji")
	_, err := parseGitmoji(getGitmojiList(client))
	if err != nil {
		t.Error(err)
	}
	teardown(vcr)
}

func TestGitmoji_iconName(t *testing.T) {
	type fields struct {
		Emoji string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Art", fields{Emoji: "🎨"}, "1f3a8.png"},
		{"Label", fields{Emoji: "🏷️"}, "1f3f7.png"},
		{"Sparkles", fields{Emoji: "✨"}, "2728.png"},
		{"Technologist", fields{Emoji: "🧑‍💻"}, "1f9d1-1f4bb.png"},
		{"Zap", fields{Emoji: "⚡️️"}, "26a1.png"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gi := &Gitmoji{
				Emoji: tt.fields.Emoji,
			}
			if got := gi.iconName(); got != tt.want {
				t.Errorf("iconName() = %v, want %v", got, tt.want)
			}
		})
	}
}
