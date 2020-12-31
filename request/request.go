package request

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"time"
)

// GetResponseBody makes a GET request to the specified url.
// The wrapper sets a generic UserAgent, a default timeout
// and allows SSL Validation to be skipped.
func GetResponseBody(origin string, skipssl bool) ([]byte, int, error) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipssl},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	request, err := http.NewRequest("GET", origin, nil)
	if err != nil {
		return nil, 0, err
	}

	request.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) "+
			"AppleWebKit/537.36 (KHTML, like Gecko) "+
			"Chrome/42.0.2311.135 "+
			"Safari/537.36 Edge/12.246",
	)

	response, err := client.Do(request)
	if err != nil {
		return nil, 0, err
	}

	defer response.Body.Close()

	pageContent, err := ioutil.ReadAll(response.Body)
	return pageContent, response.StatusCode, err
}
