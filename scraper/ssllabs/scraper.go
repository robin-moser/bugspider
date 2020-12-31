package ssllabs

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/robin-moser/bugspider/processor"
	"github.com/robin-moser/bugspider/request"

	"github.com/PuerkitoBio/goquery"
)

const originURL string = "https://www.ssllabs.com/ssltest"

// Scrape domains from the sslllabs provider
func Scrape() (*processor.Collection, error) {

	output, status, err := request.GetResponseBody(originURL, false)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, fmt.Errorf("Got a bad Status Code: %d", status)
	}

	bodyReader := bytes.NewReader(output)
	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		return nil, err
	}

	// define the generic HostArray struct
	result := processor.Collection{}

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
			host := processor.Host{
				Hostname: hostname,
				Source:   "ssllabs." + sourceType,
				Date:     time.Now(),
			}

			result.Hosts = append(result.Hosts, host)
		})
	})

	return &result, nil

}
