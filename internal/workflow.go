package internal

import (
	aw "github.com/deanishe/awgo"
	"io"
	"log"
	"net/http"
	"os"
)

func BuildItem(gi *Gitmoji, wf *aw.Workflow, client *http.Client) *aw.Item {
	err := preloadIcon(gi, wf, client)
	if err != nil {
		return nil
	}
	return wf.NewItem(gi.Description).
		Autocomplete(gi.Description).
		Match(gi.Name + " " + gi.Description).
		Subtitle(gi.Code).
		Arg(gi.Emoji).
		Icon(&aw.Icon{
			Value: wf.CacheDir() + "/" + gi.iconName(),
		}).
		Valid(true)
}

func preloadIcon(gitmoji *Gitmoji, wf *aw.Workflow, client *http.Client) error {
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
			log.Fatalln(err)
		}
	}(file)

	_, err = io.Copy(file, resp.Body)
	return err
}
