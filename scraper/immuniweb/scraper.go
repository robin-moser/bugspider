package immuniweb

import (
	"bugspider/host"
	"bugspider/request"
	"encoding/json"
	"fmt"
	"time"
)

type hostArray struct {
	Hosts []host.Host `json:"Results"`
}

func Scrape() *host.HostArray {

	originURL := "https://www.immuniweb.com/websec/api/v1/" +
		"latest/get_archived_results/get_archived_results.html"

	output, err := request.GetResponseBody(originURL)
	if err != nil {
		fmt.Printf("Error!\n%v", err)
	}

	res, err := decodeResponse(output)
	result := host.HostArray(*res)

	for i := 0; i < len(result.Hosts); i++ {
		result.Hosts[i].Date = time.Now()
		result.Hosts[i].Source = "immuniweb"
	}

	return &result

}

func decodeResponse(encodedResponse []byte) (*hostArray, error) {
	hostArray := hostArray{}
	err := json.Unmarshal(encodedResponse, &hostArray)
	if err != nil {
		return nil, err
	}
	return &hostArray, nil
}
