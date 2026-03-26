package main

import (
	"encoding/json"
	"log/slog"
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
		RPM      int    `arg:"--rpm" default:"60"`
		LogLevel string `arg:"--log-level" default:"info"`
	}

	arg.MustParse(&args)

	var level slog.Level
	err := level.UnmarshalText([]byte(args.LogLevel))

	if err != nil {
		panic(err)
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	handlerAttrs := []slog.Attr{slog.String("type", "log")}
	handler = handler.WithAttrs(handlerAttrs).(*slog.JSONHandler)
	slog.SetDefault(slog.New(handler))

	slog.Debug("Starting parser...")

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

	rpm := args.RPM

	if args.Format == "html" {
		for item, err := range scraper.ScrapeHtml(client, rpm) {

			if err != nil {
				panic(err)
			}

			_, _ = os.Stdout.Write([]byte(strings.ReplaceAll(item.Html, "\n", " ") + "\n"))
		}
	} else {
		slog.Debug("Scraping to JSON")

		for item := range scraper.Scrape(client, rpm) {
			jsonString, err := json.Marshal(item)

			if err != nil {
				panic(err)
			}

			_, _ = os.Stdout.Write([]byte(string(jsonString) + "\n"))
		}
	}
}
