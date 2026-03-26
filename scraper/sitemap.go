package scraper

import (
	"encoding/xml"
	"iter"
	"log/slog"
	"net/http"
	"strings"
)

type sitemapIndexItem struct {
	Loc string `xml:"loc"`
}

type sitemapIndex struct {
	Sitemaps []sitemapIndexItem `xml:"sitemap"`
}

type advertSitemapItem struct {
	Loc string `xml:"loc"`
}

type advertSitemap struct {
	Items []advertSitemapItem `xml:"url"`
}

// Returns list of adverts related sitemap file urls
func getAdvertSitemapUrls(client *http.Client) iter.Seq[string] {
	slog.Debug("getting advertSitemapUrls")
	return func(yield func(string) bool) {
		body, err := load(client, "https://krisha.kz/sitemap.xml")

		if err != nil {
			panic(err)
		}

		sitemap := &sitemapIndex{}
		err = xml.Unmarshal([]byte(body), sitemap)

		if err != nil {
			panic(err)
		}

		for _, item := range sitemap.Sitemaps {
			if isCorrectAdvertSitemapUrl(item.Loc) {
				yield(item.Loc)
			}
		}
	}
}

func isCorrectAdvertSitemapUrl(url string) bool {
	if url == "https://krisha.kz/sitemap/frontend/advert.xml" {
		return true
	}

	if strings.HasPrefix(url, "https://krisha.kz/sitemap/frontend/advert_") &&
		strings.HasSuffix(url, ".xml") {
		return true
	}

	return false
}

func getAdvertUrls(client *http.Client) iter.Seq[string] {
	return func(yield func(string) bool) {
		slog.Debug("getting advertUrls")
		for sitemapUrl := range getAdvertSitemapUrls(client) {

			slog.Debug("Fetching sitemap: " + sitemapUrl)
			response, err := load(client, sitemapUrl)

			if err != nil {
				slog.Error("Error fetching sitemap: " + sitemapUrl)
				panic(err)
			}

			slog.Debug("Fetched sitemap:" + sitemapUrl)

			sitemap := &advertSitemap{}
			err = xml.Unmarshal([]byte(response), sitemap)

			if err != nil {
				slog.Error("Error parsing sitemap: " + sitemapUrl)
				panic(err)
			}

			slog.Debug("Sitemap " + sitemapUrl + " has " + string(len(sitemap.Items)) + " items")

			for _, item := range sitemap.Items {
				if strings.HasPrefix(item.Loc, "https://krisha.kz/a/show/") {
					yield(item.Loc)
				}
			}
		}
	}
}
