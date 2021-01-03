# bugspider

`bugspider` is a commandline utility for domain gathering and automated vulnerability scanning.

## How it works

- **Domain Gathering**
  the tool scrapes domains from public scanning tools like `ssllabs.com` and queues them for further scanning using a beanstalk deamon. It saves the domains in a csv file for easy manual analysis.
- **Vulnerability scanning**
  the tool has differnet processors to check the found domains for differnet vunlerabilities. See the list below for available processors.

## Current Domain Providers

- [ssllabs.com](https://www.ssllabs.com) (Recently Seen / Recent Best / Recent Workst)
- [immuniweb.com](https://www.immuniweb.com/websec/#latest) (Recent Scans)

## Current Processors

- `deduplication`: test, if the domain was processed before and push it to further tubes, if not processed already

- `opengit`: test the domain for a publicly available git repository. Read the [writeup](https://en.internetwache.org/dont-publicly-expose-git-or-how-we-downloaded-your-websites-sourcecode-an-analysis-of-alexas-1m-28-07-2015/) for additional information on this vunlerability.

## Installation

### Docker

There a two docker images provided:

- a bundled image (Tag: `bundle`) including a local beanstalk daemon and an entrypoint for spawning multiple producers and workers. This image is ideal for a quick start launching only one container, but lacks the scaleability to spawn more producers/workers when needed.

- a standalone image (Tag: `latest`) containing only the commandline tool as an entrypoint. This image is ideal in a deployment cluster like Docker Swarm or Kubernetes. The image needs an external beanstalk daemon.

## Usage

```
./bugspider <command>

COMMANDS:
    worker [tubes_to_watch...]
        Starts a listener to check for duplicates and save the domain.
        The worker runs until manually stopped and waits for new jobs
        from beanstalk. Defaults to all tubes.

    scraper (ssllabs|immuniweb|file:<filepath>)
        Starts scraping one of the given providers. The scraped domains
        will be sent to beanstalk for further processing. The scraper
        runs until manually stopped and repeats the process every 10 seconds.
        If the scraper uses a file, it scans all lines from the file
        and exists, when finished.

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

## Configuration

To configure the tubes, which the worker should be listen to, you have to use command line arguments (see Usage).

To configure the beanstalk host (and port), you can set the environment Variable `BENASTALK_HOST`

## Examples

Using `bugspider` as a bundled, ready to use, version:

```sh
docker run -d \
    --name=bugspider \
    -v $(pwd)/output:/app/output \
    --net=container:openvpn robinmoser/bugspider:bundle
```

Using multiple  `bugspider`  instances with the standalone image:

```sh
# create a shared network
docker network create bugspider-net

# start a beanstalk daemon
docker run -d \
    --network bugspider-web \
    --name=beanstalk schickling/beanstalkd

# start multiple workers
docker run -d \
    --network bugspider-web \
    -e BEANSTALK_HOST=beanstalk:11300 \
    robinmoser.de/bugspider:latest worker deduplication

docker run -d \
    --network bugspider-web \
    -e BEANSTALK_HOST=beanstalk:11300 \
    robinmoser.de/bugspider:latest worker opengit

docker run -d \
    --network bugspider-web \
    -e BEANSTALK_HOST=beanstalk:11300 \
    robinmoser.de/bugspider:latest worker opengit

# start both producers
docker run -d \
    --network bugspider-web \
    -e BEANSTALK_HOST=beanstalk:11300 \
    robinmoser.de/bugspider:latest scraper ssllabs

docker run -d \
    --network bugspider-web \
    -e BEANSTALK_HOST=beanstalk:11300 \
    robinmoser.de/bugspider:latest scraper immuniweb
```
