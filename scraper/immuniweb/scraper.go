package immuniweb

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/robin-moser/bugspider/processor"
	"github.com/robin-moser/bugspider/request"
)

type hostCollection struct {
	Hosts []processor.Host `json:"Results"`
}

const originURL string = "https://www.immuniweb.com/websec/api/v1/" +
	"latest/get_archived_results/get_archived_results.html"

// Scrape domains from the immuniweb provider
func Scrape() (*processor.Collection, error) {

	output, status, err := request.GetResponseBody(originURL, false)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, fmt.Errorf("Got a bad Status Code: %d", status)
	}

	res, err := decodeResponse(output)
	if err != nil {
		return nil, err
	}

	result := processor.Collection(*res)

	for i := 0; i < len(result.Hosts); i++ {
		result.Hosts[i].Date = time.Now()
		result.Hosts[i].Source = "immuniweb"
	}

	return &result, nil

}

func decodeResponse(encodedResponse []byte) (*hostCollection, error) {
	hostCollection := hostCollection{}
	err := json.Unmarshal(encodedResponse, &hostCollection)
	if err != nil {
		return nil, err
	}
	return &hostCollection, nil
}
