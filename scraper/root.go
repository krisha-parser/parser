package scraper

import (
	"iter"
	"net/http"
)

func ScrapeHtml(client *http.Client) iter.Seq[string] {
	return func(yield func(string) bool) {
		for adUrl := range getAdvertUrls(client) {
			response, err := load(client, adUrl)

			if err != nil {
				panic(err)
			}

			yield(response)
		}
	}
}

func Scrape(client *http.Client) iter.Seq[advert] {
	return func(yield func(advert) bool) {
		for html := range ScrapeHtml(client) {
			yield(parseAdvert(html))
		}
	}
}
