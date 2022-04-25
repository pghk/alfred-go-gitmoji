package internal

import (
	"encoding/json"
	aw "github.com/deanishe/awgo"
	"github.com/kinbiko/jsonassert"
	"testing"
)

func Test_loadIcon(t *testing.T) {
	client, vcr := setup("load_icon")

	type args struct {
		gitmoji *Gitmoji
	}
	tests := []struct {
		name string
		args args
	}{
		{"Sparkles", args{&Gitmoji{Emoji: "✨"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wf := aw.New()
			if wf.Cache.Exists(tt.args.gitmoji.iconName()) {
				t.Errorf("Already loaded")
			}
			err := preloadIcon(tt.args.gitmoji, wf, client)
			if err != nil {
				t.Error(err)
			}
			if !wf.Cache.Exists(tt.args.gitmoji.iconName()) {
				t.Errorf("Did not load")
			}
			if err != nil {
				t.Error(err)
			}
			err = wf.Cache.Store(tt.args.gitmoji.iconName(), nil)
			if err != nil {
				t.Error(err)
			}
		})
	}

	teardown(vcr)
}

func Test_item(t *testing.T) {
	wf := aw.New()
	client, vcr := setup("item_test")
	ja := jsonassert.New(t)
	builtItem := BuildItem(&Gitmoji{
		"🐛",
		"&#x1f41b;",
		":bug:",
		"Fix a bug.",
		"bug",
		"patch",
	}, wf, client)
	got, err := json.Marshal(builtItem)
	if err != nil {
		t.Error(err)
	}
	ja.Assertf(string(got), `{
		"title": "Fix a bug.",
		"autocomplete": "Fix a bug.",
		"match": "bug Fix a bug.",
		"subtitle": ":bug:",
		"arg": "🐛",
		"icon": {"path": "%s"},
		"valid": true
	}`, wf.CacheDir() + "/1f41b.png")
	teardown(vcr)
}
