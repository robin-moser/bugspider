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
func BsProducer(source string) {
	protocol := beanstalk.MakeJSONHostProtocol()
	producer := beanstalk.MakeNewProducer("localhost:11300", protocol)
	err := producer.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	// scrape the source, return a Host Collection
	hostCollection, err := scraper.Scrape(source)
	if err != nil {
		log.Fatal(err)
	}

	// loop through all recieved Hosts and store them one by one
	for _, host := range hostCollection.Hosts {
		producer.PutHost(&host)
	}
}

// BsWorker listens to the job queue and processes active jobs
func BsWorker() {
	protocol := beanstalk.MakeJSONHostProtocol()
	processor := host.MakeNewHostProcessor("output/hostfile.txt")
	worker := beanstalk.MakeNewWorker("localhost:11300", protocol, processor)
	err := worker.Connect()
	if err != nil {
		log.Fatal(err)
	}

	defer worker.Close()

	for {
		worker.ProcessJob()
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
		BsWorker()
	} else if os.Args[1] == "scraper" && len(os.Args) == 3 {
		BsProducer(os.Args[2])
	} else {
		printUsage()
	}
}
