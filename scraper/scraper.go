package scraper

import (
	"errors"

	"github.com/robin-moser/bugspider/processor"
	"github.com/robin-moser/bugspider/scraper/immuniweb"
	"github.com/robin-moser/bugspider/scraper/ssllabs"
)

// Scrape domains from various Sources.
// A predefined source needs to be given.
func Scrape(source string) (*processor.Collection, error) {

	switch source {
	case "immuniweb":
		return populateHosts(immuniweb.Scrape())

	case "ssllabs":
		return populateHosts(ssllabs.Scrape())

	default:
		err := errors.New("Undefined Scrape Source")
		return nil, err
	}

}

func populateHosts(hostarray *processor.Collection, err error) (*processor.Collection, error) {
	if err != nil {
		return nil, err
	}

	if len(hostarray.Hosts) == 0 {
		return nil, errors.New("Error: no domains found")
	}

	for i := 0; i < len(hostarray.Hosts); i++ {
		hostarray.Hosts[i].JobType = "deduplication"
		hostarray.Hosts[i].Retries = 0
	}

	return hostarray, nil
}
