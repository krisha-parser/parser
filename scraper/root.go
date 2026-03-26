package scraper

import (
	"iter"
	"log/slog"
	"net/http"
	"time"
)

type scrapedPage struct {
	Url  string
	Html string
}

func ScrapeHtml(client *http.Client, rpm int) iter.Seq2[scrapedPage, error] {
	slog.Debug("Scraping HTML")
	return func(yield func(scrapedPage, error) bool) {
		count := 0
		minuteStart := time.Now()

		for adUrl := range getAdvertUrls(client) {
			slog.Debug("Processing ad: " + adUrl)
			if count == rpm {
				elapsed := time.Since(minuteStart)
				if elapsed < time.Minute {
					slog.Debug("Hitting RPM, waiting...")
					time.Sleep(time.Minute - elapsed)
				}
				count = 0
				minuteStart = time.Now()
			}

			count++

			slog.Debug("Loading ad: " + adUrl)
			response, err := load(client, adUrl)
			slog.Debug("Loaded ad: " + adUrl)

			if err != nil {
				slog.Error("Error loading ad: " + adUrl)
				panic(err)
			}

			if !yield(scrapedPage{
				Url:  adUrl,
				Html: response,
			}, err) {
				slog.Error("Error yielding: " + adUrl)
				return
			}
		}
	}
}

func Scrape(client *http.Client, rpm int) iter.Seq2[advert, error] {
	return func(yield func(advert, error) bool) {
		for page, err := range ScrapeHtml(client, rpm) {
			yield(parseAdvert(client, page.Url, page.Html), err)
		}
	}
}
