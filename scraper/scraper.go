package scraper

import (
	"errors"

	"bugspider/host"
	"bugspider/scraper/immuniweb"
	"bugspider/scraper/ssllabs"
)

// Scrape domains from various Sources.
// A predefined source needs to be given.
func Scrape(source string) (*host.Collection, error) {

	switch source {
	case "immuniweb":
		hostarray, err := immuniweb.Scrape()
		if err != nil {
			return nil, err
		} else if len(hostarray.Hosts) == 0 {
			return nil, errors.New("Error: no domains found")
		}

		return hostarray, nil

	case "ssllabs":
		hostarray, err := ssllabs.Scrape()
		if err != nil {
			return nil, err
		} else if len(hostarray.Hosts) == 0 {
			return nil, errors.New("Error: no domains found")
		}

		return hostarray, nil

	default:
		err := errors.New("Undefined Scrape Source")
		return nil, err
	}

}
