package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const GithubUrl string = "https://raw.githubusercontent.com"
const GitmojiIndex string = "carloscuesta/gitmoji/master/src/data/gitmojis.json"

func getGitmojiList(client *http.Client) string {
	resp, err := client.Get(GithubUrl + "/" + GitmojiIndex)
	if err != nil {
		log.Fatalln(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(body)
}

type Gitmoji struct {
	Emoji       string      `json:"emoji"`
	Entity      string      `json:"entity"`
	Code        string      `json:"code"`
	Description string      `json:"description"`
	Name        string      `json:"name"`
	Semver      interface{} `json:"semver"`
}

func parseGitmoji(input string) (map[string][]Gitmoji, error) {
	test := map[string][]Gitmoji{}
	err := json.Unmarshal([]byte(input), &test)
	return test, err
}

// iconName derives PNG filenames from emoji, according to the convention used by JoyPixels in their repository
func (g *Gitmoji) iconName() string {
	var invisible = map[string]string{
		"200d": "Zero Width Joiner",
		"fe0f": "Variation Selector-16",
	}
	var hexes []string
	for _, r := range g.Emoji {
		codepoint := fmt.Sprintf("%x", r)
		_, found := invisible[codepoint]
		include := !found
		if include {
			hexes = append(hexes, codepoint)
		}
	}
	return strings.Join(hexes, "-") + ".png"
}

