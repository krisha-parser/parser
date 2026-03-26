package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/krisha-parser/parser/scraper"
)

func main() {
	var args struct {
		ProxyUrl string `arg:"--proxy-url"`
		Format   string `arg:"--format" default:"json"`
	}

	arg.MustParse(&args)

	var transport = &http.Transport{}

	if args.ProxyUrl != "" {
		proxyUrl, err := url.Parse(args.ProxyUrl)

		if err != nil {
			panic(err)
		}

		transport.Proxy = http.ProxyURL(proxyUrl)
		transport.MaxIdleConnsPerHost = 100
		transport.MaxIdleConns = 100
		transport.MaxConnsPerHost = 100
	}

	client := &http.Client{
		Transport: transport,
	}

	if args.Format == "html" {
		for item := range scraper.ScrapeHtml(client) {
			_, _ = os.Stdout.Write([]byte(strings.ReplaceAll(item, "\n", " ") + "\n"))
		}
	} else {
		for item := range scraper.Scrape(client) {
			jsonString, err := json.Marshal(item)

			if err != nil {
				panic(err)
			}

			_, _ = os.Stdout.Write([]byte(string(jsonString) + "\n"))
		}
	}
}
