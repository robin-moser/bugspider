package ssllabs

import (
	"bytes"
	"net/url"
	"strings"
	"time"

	"bugspider/host"
	"bugspider/request"

	"github.com/PuerkitoBio/goquery"
)

const originURL string = "https://www.ssllabs.com/ssltest"

// Scrape domains from the sslllabs provider
func Scrape() (*host.Collection, error) {

	output, err := request.GetResponseBody(originURL, false)
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewReader(output)
	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		return nil, err
	}

	// define the generic HostArray struct
	result := host.Collection{}

	doc.Find(".boxContent").Each(func(i int, s *goquery.Selection) {

		var sourceType string
		sourceType = (s.Find(".boxHead").First().Text())
		sourceType = strings.Replace(sourceType, " ", "", -1)

		s.Find("a").Each(func(j int, s *goquery.Selection) {

			hostname, exists := s.Attr("href")
			if !exists {
				return
			}

			hostname, err := url.QueryUnescape(hostname)
			if err != nil {
				return
			}

			hostname = strings.TrimPrefix(hostname, "analyze.html?d=")
			host := host.Host{
				Hostname: hostname,
				Source:   "ssllabs." + sourceType,
				Date:     time.Now(),
			}

			result.Hosts = append(result.Hosts, host)
		})
	})

	return &result, nil

}
