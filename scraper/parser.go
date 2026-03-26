package scraper

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/PuerkitoBio/goquery"
)

type advertParameter struct {
	Label string `json:"name"`
	Code  string `json:"code"`
	Value string `json:"value"`
}

type advert struct {
	ID                   int               `json:"id"`
	Type                 string            `json:"type"`
	Url                  string            `json:"Url"`
	Name                 string            `json:"name"`
	Price                string            `json:"price"`
	IsMortgaged          bool              `json:"is_mortgaged"`
	Parameters           []advertParameter `json:"parameters"`
	DescriptionParaments []advertParameter `json:"description_paraments"`
	Description          string            `json:"description"`
}

func parseAdvert(client *http.Client, url string, html string) advert {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {
		slog.Error("Error parsing HTML document")
		panic(err)
	}

	ad := advert{
		Type:                 "ad",
		Parameters:           []advertParameter{},
		DescriptionParaments: []advertParameter{},
	}

	ad.Url = url

	urlParts := strings.Split(url, "/")
	ad.ID, _ = strconv.Atoi(urlParts[len(urlParts)-1])

	nameNode := doc.Find(".offer__advert-title > h1").First()
	ad.Name = strings.TrimSpace(nameNode.Text())

	priceNode := doc.Find(".offer__price").First()
	ad.Price = strings.TrimSpace(priceNode.Text())

	parametersNode := doc.Find(".offer__info-item")

	for i := range parametersNode.Nodes {
		parameterNode := parametersNode.Eq(i)

		parameterValueNode := parameterNode.Find(".offer__advert-short-info").First()
		codeAttr, _ := parameterNode.Attr("data-name")

		if parameterValueNode.HasClass("offer__location") {
			parameterValueNode = parametersNode.Find("span").First()
			codeAttr = "city"
		}

		ad.Parameters = append(ad.Parameters, advertParameter{
			Label: strings.TrimSpace(parameterNode.Find(".offer__info-title").First().Text()),
			Value: strings.TrimSpace(parameterValueNode.Text()),
			Code:  codeAttr,
		})
	}

	isMortgagedNode := doc.Find(".offer__parameters-mortgaged")
	ad.IsMortgaged = isMortgagedNode.Length() > 0

	descriptionParametersNode := doc.Find("div.offer__parameters dl")

	for i := range descriptionParametersNode.Nodes {
		descriptionParameterNode := descriptionParametersNode.Eq(i)
		labelNode := descriptionParameterNode.Find("dt").First()
		valueNode := descriptionParameterNode.Find("dd").First()
		codeAttr, _ := labelNode.Attr("data-name")

		ad.DescriptionParaments = append(ad.DescriptionParaments, advertParameter{
			Label: strings.TrimSpace(labelNode.Text()),
			Value: strings.TrimSpace(valueNode.Text()),
			Code:  codeAttr,
		})
	}

	descriptionNode := doc.Find(".offer__description > .text > div").First()
	descHtml, err := descriptionNode.Html()

	if err != nil {
		panic(err)
	}

	descMD, err := htmltomarkdown.ConvertString(descHtml)

	if err != nil {
		panic(err)
	}

	ad.Description = descMD

	// TODO: Make it work!
	//phoneResponse, err := load(client, "https://krisha.kz/a/ajaxPhones?id="+strconv.Itoa(ad.ID), withHeaders(http.Header{
	//	"Content-Type":     []string{"application/json; charset=utf-8"},
	//	"Accept":           []string{"application/json", "text/plain", "*/*"},
	//	"Referer":          []string{url},
	//	"X-Requested-With": []string{"XMLHttpRequest"},
	//	"User-Agent":       []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36"},
	//}))

	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println("Phone Response: " + phoneResponse)
	return ad
}
