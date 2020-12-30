package beanstalk

import (
	"encoding/json"
	"log"
	"time"

	"github.com/iwanbk/gobeanstalk"
	"github.com/robin-moser/bugspider/host"
)

func (bs *Handler) Watch(tubes ...string) error {

	var watching int
	var err error

	for _, tube := range tubes {
		watching, err = bs.serverConnection.Watch(tube)
		if err != nil {
			return err
		}
	}
	log.Println("watching", watching, "tubes")
	return nil
}

func (bs *Handler) ProcessJob(processor host.Processor) {
	job, err := bs.serverConnection.Reserve()
	if err != nil {
		log.Println(err)
		return
	}

	currentHost := host.Host{}
	err = json.Unmarshal(job.Body, &currentHost)
	if err != nil {
		bs.handleError(job, err)
		return
	}

	inserted, err := processor.Process(&currentHost)
	if err != nil {
		bs.handleError(job, err)
		return
	}
	if inserted {
		log.Printf("processed Job ID %v: saved %v\n", job.ID, currentHost.Hostname)
	}
	bs.serverConnection.Delete(job.ID)
}

// handleError gets called, when a job cant finish, so it can be released
// to get processed at a later time
func (bs *Handler) handleError(job *gobeanstalk.Job, err error) {
	log.Println(err)
	priority := uint32(5)
	delay := 20 * time.Second
	bs.serverConnection.Release(job.ID, priority, delay)
}
