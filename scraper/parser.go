package scraper

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type advertParameter struct {
	Label string `json:"name"`
	Value string `json:"value"`
}

type advert struct {
	ID         int               `json:"id"`
	Url        string            `json:"url"`
	Name       string            `json:"name"`
	Price      string            `json:"price"`
	Parameters []advertParameter `json:"parameters"`
}

func parseAdvert(html string) advert {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {
		log.Fatal(err)
	}

	ad := advert{
		Parameters: []advertParameter{},
	}

	urlMetaNode := doc.Find("link[rel=\"canonical\"]").First()
	ad.Url = urlMetaNode.AttrOr("href", "")

	nameNode := doc.Find(".offer__advert-title > h1").First()
	ad.Name = strings.TrimSpace(nameNode.Text())

	priceNode := doc.Find(".offer__price").First()
	ad.Price = strings.TrimSpace(priceNode.Text())

	parametersNode := doc.Find(".offer__info-item")

	for i := range parametersNode.Nodes {
		parameterNode := parametersNode.Eq(i)

		parameterValueNode := parameterNode.Find(".offer__advert-short-info").First()

		if parameterValueNode.HasClass("offer__location") {
			parameterValueNode = parametersNode.Find("span").First()
		}

		ad.Parameters = append(ad.Parameters, advertParameter{
			Label: strings.TrimSpace(parameterNode.Find(".offer__info-title").First().Text()),
			Value: strings.TrimSpace(parameterValueNode.Text()),
		})
	}

	return ad
}
