# krisha.kz парсер

Парсер объявлений сайта krisha.kz с прокси в JSON

## Использование через терминал

```shell
go build -o parser .
./parser --proxy-url "http://username:password@rotating-datacenter.geonode.io:9000"

# HTML - пропускает парсинг в json и отдает сырой HTML 
./parser --format html
```

Выводит объявления построчно

Пример вывода:
```json
{"id":761098335,"url":"https://krisha.kz/a/show/761098335","name":"3-комнатная квартира · 105.44 м², Е-899 7","price":"от 82 981 280 〒","parameters":[{"name":"Город","value":"Астана, Нура р-н"},{"name":"Тип дома","value":"монолитный"},{"name":"Жилой комплекс","value":"GreenLine. Garden"},{"name":"Год постройки","value":"2026"},{"name":"Площадь","value":"105.44 м²"},{"name":"Санузел","value":"2 с/у и более"},{"name":"Высота потолков","value":"3 м"}]}
```

## Использование через код

```go
package main

import (
	"net/http"

	"github.com/krisha-parser/parser"
)

func main() {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}

	for ad := range scraper.Scrape(client) {
		// do smth with ad
	}
}

```