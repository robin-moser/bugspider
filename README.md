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
    worker [tubes_to_watch...]
        Starts a listener to check for duplicates and save the domain.
        The worker runs until manually stopped and waits for new jobs
        from beanstalk. Defaults to all tubes.

    scraper (ssllabs|immuniweb)
        Starts scraping one of the given providers. The scraped domains
        will be sent to beanstalk for further processing. The scraper
        quits after a single provider fetch.

TUBES:
    deduplication
        Checks from a CSV Log File, if a host was already processed.
        If not already processed, the Host will be pushed to other tubes
        for further processing.

    opengit
        Checks a host for public accessable git repositories (domain.com/.git).
        If a Host ist vulnerable, it will be logged, the accessable git config
        will be stored in a separate 'config' dir.

```
