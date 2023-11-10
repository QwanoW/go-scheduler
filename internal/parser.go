package internal

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

const (
	ParseTarget = "mtli_doc"
)

var viewerTemplate string = "https://view.officeapps.live.com/op/view.aspx?src="

func FindLinks(url string, parseTarget string) (map[string]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err := resp.Body.Close()
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	body, err := html.Parse(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}

	return filterHTML(body, parseTarget), nil
}

func filterHTML(n *html.Node, parseTarget string) map[string]string {
	// создаем новый документ из узла n
	doc := goquery.NewDocumentFromNode(n)
	// создаем пустой слайс для хранения найденных элементов
	result := make(map[string]string)
	// ищем все элементы с атрибутом data-foo, равным bar
	doc.Find("a." + parseTarget).Each(func(i int, s *goquery.Selection) {
		// добавляем href в слайс
		text := s.Text()
		attr, _ := s.Attr("href")
		result[text] = viewerTemplate + attr
	})
	// возвращаем слайс
	return result
}
