package main

import (
	"bugspider/beanstalk"
	"bugspider/host"
	"bugspider/scraper"
	"fmt"
	"os"
)

func BsProducer(source string) {
	protocol := beanstalk.MakeJsonHostProtocol()
	producer := beanstalk.MakeNewProducer("localhost:11300", protocol)
	producer.Connect()
	defer producer.Close()

	hosts, err := scraper.Scrape(source)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, host := range hosts.Hosts {
		producer.PutHost(&host)
	}
}

func BsWorker() {
	protocol := beanstalk.MakeJsonHostProtocol()
	processor := host.MakeNewHostProcessor("output/hostfile.txt")
	worker := beanstalk.MakeNewWorker("localhost:11300", protocol, processor)
	worker.Connect()
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
