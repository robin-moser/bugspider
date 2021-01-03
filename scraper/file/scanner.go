package file

import (
	"bufio"
	"os"
	"time"

	"github.com/robin-moser/bugspider/processor"
)

// Scrape domains from the immuniweb provider
func Scrape(filepath string) (*processor.Collection, error) {

	f, err := os.OpenFile(filepath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// define the generic HostArray struct
	result := processor.Collection{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {

		host := processor.Host{
			Hostname: scanner.Text(),
			Source:   "file." + filepath,
			Date:     time.Now(),
		}

		result.Hosts = append(result.Hosts, host)

	}

	return &result, nil

}
