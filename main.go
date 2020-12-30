package main

import (
	"fmt"
	"log"
	"os"

	"github.com/robin-moser/bugspider/beanstalk"
	"github.com/robin-moser/bugspider/host"
	"github.com/robin-moser/bugspider/scraper"
)

// BsProducer scrapes the given source and puts the results to beanstalk
func BsProducer(source string, tube string) {
	bs := beanstalk.NewHandler("localhost:11300")
	err := bs.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer bs.Close()

	err = bs.UseTube(tube)
	if err != nil {
		log.Fatal(err)
	}

	// scrape the source, return a Host Collection
	hostCollection, err := scraper.Scrape(source)
	if err != nil {
		log.Fatal(err)
	}

	// loop through all recieved Hosts and store them one by one
	for _, host := range hostCollection.Hosts {
		bs.PutHost(&host)
	}
}

// BsWorker listens to the job queue and processes active jobs
func BsWorker(tubes ...string) {
	processor := host.NewProcessor("csv", "output/hostfile.txt")
	bs := beanstalk.NewHandler("localhost:11300")
	err := bs.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer bs.Close()

	bs.Watch(tubes...)

	for {
		bs.ProcessJob(*processor)
	}
}

func printUsage() {
	fmt.Printf("Usage: %v <command>\n\n", os.Args[0])
	fmt.Println("Commands:")
	fmt.Println("  scraper immuniweb|ssllabs")
	fmt.Println("  worker")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
	}

	if os.Args[1] == "worker" {

		if len(os.Args) >= 3 {
			tubes := os.Args[2:]
			BsWorker(tubes...)
		} else {
			BsWorker("domainGathering")
		}

	} else if os.Args[1] == "scraper" && len(os.Args) == 3 {
		BsProducer(os.Args[2], "domainGathering")

	} else {
		printUsage()
	}
}
