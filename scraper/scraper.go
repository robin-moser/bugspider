package scraper

import (
	"bugspider/host"
	"errors"
)

// Scrape - Scrapes URLs from various Sources
func Scrape(source string) (*host.HostArray, error) {

	switch source {
	case "immuniweb":
		return scrapeImmuniweb(), nil
	default:
		err := errors.New("Undefined Scrape Source")
		return nil, err
	}

}
