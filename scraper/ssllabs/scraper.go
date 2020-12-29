package ssllabs

import (
	"bugspider/host"
	"bugspider/request"
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func Scrape() *host.HostArray {

	originURL := "https://www.ssllabs.com/ssltest"

	output, err := request.GetResponseBody(originURL)
	if err != nil {
		fmt.Printf("Error!\n%v", err)
	}

	bodyReader := bytes.NewReader(output)
	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		fmt.Printf("Error!\n%v", err)
	}

	result := host.HostArray{}

	doc.Find(".boxContent").Each(func(i int, s *goquery.Selection) {

		sourceType := (s.Find(".boxHead").First().Text())
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

	return &result

}
