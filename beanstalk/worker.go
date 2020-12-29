package beanstalk

import (
	"bugspider/host"
	"fmt"
	"os"
	"time"

	"github.com/iwanbk/gobeanstalk"
)

type HostWorker struct {
	MainBeanstalk
	protocol  host.HostProtocol
	processor host.HostProcessor
}

func (worker *HostWorker) ProcessJob() {
	job, err := worker.serverConnection.Reserve()
	if err != nil {
		fmt.Println(err)
		return
	}
	host, err := worker.protocol.Decode(job.Body)
	if err != nil {
		worker.handleError(job, err)
		return
	}

	inserted, err := worker.processor.DoProcess(host)
	if err != nil {
		worker.handleError(job, err)
		return
	} else if inserted {
		fmt.Printf("processed Job ID %v: inserted %v\n", job.ID, host.Hostname)
	}
	worker.serverConnection.Delete(job.ID)
}

func (worker *HostWorker) watch() error {
	watching, err := worker.serverConnection.Watch("default")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("watching ", watching, " tubes")
	return nil
}

func (worker *HostWorker) Connect() {
	worker.MainBeanstalk.Connect()
	worker.watch()
}

func (worker *HostWorker) handleError(job *gobeanstalk.Job, err error) {
	fmt.Println(err)
	priority := uint32(5)
	delay := 0 * time.Second
	worker.serverConnection.Release(job.ID, priority, delay)
}

func MakeNewWorker(serverAddress string, protocol host.HostProtocol, processor host.HostProcessor) *HostWorker {
	worker := HostWorker{protocol: protocol, processor: processor}
	worker.ServerAddress = serverAddress
	return &worker
}
