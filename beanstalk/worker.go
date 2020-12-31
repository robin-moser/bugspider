package beanstalk

import (
	"encoding/json"
	"log"
	"time"

	"github.com/iwanbk/gobeanstalk"
	"github.com/robin-moser/bugspider/processor"
)

func (bs *Handler) Watch(tubes []string) error {
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

func (bs *Handler) ProcessJob() {
	job, err := bs.serverConnection.Reserve()
	if err != nil {
		log.Println(err)
		return
	}

	currentHost := &processor.Host{}
	err = json.Unmarshal(job.Body, &currentHost)
	if err != nil {
		bs.handleError(job, err)
		return
	}

	// process the current Host, the used processor is defines in the Host struct
	hostProc, err := processor.ProcessHost(currentHost, int(job.ID))
	if err != nil {
		bs.handleError(job, err)
		return
	}

	// send returned Host to further tubes for further processing
	for _, tube := range hostProc.Tubes {
		err = bs.UseTube(tube)
		hostProc.Host.JobType = tube
		if err != nil {
			bs.handleError(job, err)
			return
		}
		err = bs.PutHost(hostProc.Host, hostProc.Priority)
		if err != nil {
			bs.handleError(job, err)
			return
		}
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
