package scraper

import (
	"errors"
	"io"
	"net/http"
)

func load(client *http.Client, url string) (string, error) {
	response, err := client.Get(url)

	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		return "", errors.New(response.Status)
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}
