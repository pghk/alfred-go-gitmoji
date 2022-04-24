package internal

import (
	aw "github.com/deanishe/awgo"
	"io"
	"log"
	"net/http"
	"os"
)

var wf *aw.Workflow

func init() {
	wf = aw.New()
}

func loadIcon(gitmoji *Gitmoji, client *http.Client) error {
	fileName := gitmoji.iconName()
	if wf.Cache.Exists(fileName) {
		return nil
	}
	resp, err := client.Get(gitmoji.iconFile())
	if err != nil {
		log.Fatalln(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(resp.Body)

	file, err := os.Create(wf.CacheDir() + "/" + fileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	_, err = io.Copy(file, resp.Body)
	return err
}
