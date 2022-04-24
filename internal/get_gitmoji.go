package internal

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
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