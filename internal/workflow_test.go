package internal

import (
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
			if wf.Cache.Exists(tt.args.gitmoji.iconName()) {
				t.Errorf("Already loaded")
			}
			err := loadIcon(tt.args.gitmoji, client)
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
