package beanstalk

import (
	"encoding/json"
	"log"
	"regexp"

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
		bs.handleError(job, currentHost, err)
		return
	}

	// process the current Host, the used processor is defines in the Host struct
	hostProc, err := processor.ProcessHost(currentHost, int(job.ID))
	if err != nil {
		bs.handleError(job, currentHost, err)
		return
	}

	// send returned Host to further tubes for further processing
	for _, tube := range hostProc.Tubes {
		err = bs.UseTube(tube)
		hostProc.Host.JobType = tube
		if err != nil {
			bs.handleError(job, currentHost, err)
			return
		}
		err = bs.PutHost(hostProc.Host, hostProc.Priority, DefaultDelay)
		if err != nil {
			bs.handleError(job, currentHost, err)
			return
		}
	}

	bs.serverConnection.Delete(job.ID)

}

// handleError gets called, when a job cant finish, so it can be released
// to get processed at a later time
func (bs *Handler) handleError(job *gobeanstalk.Job, currentHost *processor.Host, err error) {

	currentHost.Retries++
	requestErr := regexp.MustCompile(`\(.*\)`).FindString(err.Error())

	log.Printf("Error: [%d] %s %s", currentHost.Retries, currentHost.Hostname, requestErr)

	bs.serverConnection.Delete(job.ID)

	if currentHost.Retries < MaxRetries {
		priority := uint32(20)
		delay := DelayOnError
		bs.UseTube(currentHost.JobType)
		bs.PutHost(currentHost, priority, delay)
	}
}
