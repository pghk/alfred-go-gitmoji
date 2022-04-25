package main

import (
	"alfred-go-gitmoji/internal"
	"flag"
	aw "github.com/deanishe/awgo"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

var (
	cacheName   = "gitmoji_index.json"
	maxCacheAge = 1 * time.Minute

	inRefreshMode bool
	query         string

	wf     *aw.Workflow
	client *http.Client
)

func init() {
	flag.BoolVar(&inRefreshMode, "refresh", false, "Populate local gitmojidex from GitHub repository")

	wf = aw.New()
	client = &http.Client{}
}

func main() {
	wf.Run(run)
}

func run() {
	setArgs()

	if inRefreshMode {
		fetchData()
		return
	}

	log.Printf("[main] query=%s", query)

	cacheData := map[string][]internal.Gitmoji{}
	if wf.Cache.Exists(cacheName) {
		if err := wf.Cache.LoadJSON(cacheName, &cacheData); err != nil {
			wf.FatalError(err)
		}
	}

	if wf.Cache.Expired(cacheName, maxCacheAge) {
		// Loop in the background until cache is fresh
		wf.Rerun(0.3)
		if !wf.IsRunning("refresh") {
			startBackgroundJob()
		} else {
			log.Printf("download in progress.")
		}
		// Add placeholder item if cache is empty.
		if len(cacheData) == 0 {
			wf.NewItem("Downloading gitmoji index…").
				Icon(aw.IconInfo)
			wf.SendFeedback()
			return
		}
	}

	gitmojidex, cacheOk := cacheData["gitmojis"]
	if !cacheOk {
		wf.Fatal("Error parsing cached gitmoji index")
	}

	for _, gitmoji := range gitmojidex {
		internal.BuildItem(&gitmoji, wf, client)
	}

	if query != "" {
		wf.Filter(query)
	}

	wf.WarnEmpty("No matching gitmoji found", "Try a different query?")

	wf.SendFeedback()
}

func startBackgroundJob() {
	cmd := exec.Command(os.Args[0], "-refresh")
	if err := wf.RunInBackground("refresh", cmd); err != nil {
		wf.FatalError(err)
	}
}

func fetchData() {
	wf.Configure(aw.TextErrors(true))
	log.Printf("[main] fetching gitmoji index...")
	gitmoji, err := internal.GetGitmoji(client)
	if err != nil {
		wf.FatalError(err)
	}
	if err := wf.Cache.StoreJSON(cacheName, gitmoji); err != nil {
		wf.FatalError(err)
	}
	log.Printf("[main] downloaded gitmoji index")
	return
}

func setArgs() {
	wf.Args()
	flag.Parse()

	if args := flag.Args(); len(args) > 0 {
		query = args[0]
	}
}
