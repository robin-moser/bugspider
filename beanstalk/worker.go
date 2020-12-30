package beanstalk

import (
	"log"
	"time"

	"github.com/robin-moser/bugspider/host"
	"github.com/iwanbk/gobeanstalk"
)

type HostWorker struct {
	MainBeanstalk
	protocol  host.Protocol
	processor host.Processor
}

func (worker *HostWorker) Connect() error {
	err := worker.MainBeanstalk.Connect()
	if err != nil {
		return err
	}

	return worker.watch()
}

func (worker *HostWorker) watch() error {
	watching, err := worker.serverConnection.Watch("default")
	if err != nil {
		return err
	}
	log.Println("watching", watching, "tubes")
	return nil
}

func (worker *HostWorker) ProcessJob() {
	job, err := worker.serverConnection.Reserve()
	if err != nil {
		log.Println(err)
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
	}
	if inserted {
		log.Printf("processed Job ID %v: saved %v\n", job.ID, host.Hostname)
	}
	worker.serverConnection.Delete(job.ID)
}

// handleError gets called, when a job cant finish, so it can be released
// to get processed at a later time
func (worker *HostWorker) handleError(job *gobeanstalk.Job, err error) {
	log.Println(err)
	priority := uint32(5)
	delay := 20 * time.Second
	worker.serverConnection.Release(job.ID, priority, delay)
}

func MakeNewWorker(serverAddress string, protocol host.Protocol, processor host.Processor) *HostWorker {
	worker := HostWorker{protocol: protocol, processor: processor}
	worker.ServerAddress = serverAddress
	return &worker
}
