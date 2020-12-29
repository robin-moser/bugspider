package scraper

import (
	"bugspider/host"
	"bugspider/scraper/immuniweb"
	"bugspider/scraper/ssllabs"
	"errors"
)

// Scrape - Scrapes URLs from various Sources
func Scrape(source string) (*host.HostArray, error) {

	switch source {
	case "immuniweb":
		return immuniweb.Scrape(), nil
	case "ssllabs":
		return ssllabs.Scrape(), nil
	default:
		err := errors.New("Undefined Scrape Source")
		return nil, err
	}

}
