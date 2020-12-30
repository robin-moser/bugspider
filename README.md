# bugspider

This Application scrapes domains from public scanning tools like ssllabs.com and saves them for further scanning. It uses beanstalkd for job queuing and saves the gathered domains to a csv file.

## Current Providers

- [ssllabs.com](https://www.ssllabs.com) (Recently Seen / Recent Best / Recent Workst)
- [immuniweb.com](https://www.immuniweb.com/websec/#latest) (Recent Scans)

## Requirements

- a running beanstalkd instance

## Usage

```
./bugspider <command>

COMMANDS:
    worker
       Starts a listener to check for duplicates and save the domain.
       The worker runs until manually stopped and waits for new jobs
       from beanstalk.

   scraper (ssllabs|immuniweb)
       Starts scraping one of the given providers. The scraped domains
       will be sent to beanstalk for further processing. The scraper
       quits after a single provider fetch.
```
