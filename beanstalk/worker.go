package beanstalk

import (
	"log"
	"time"

	"github.com/iwanbk/gobeanstalk"
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

// handleError gets called, when a job cant finish, so it can be released
// to get processed at a later time
func (bs *Handler) handleError(job *gobeanstalk.Job, err error) {
	log.Println(err)
	priority := uint32(5)
	delay := 20 * time.Second
	bs.serverConnection.Release(job.ID, priority, delay)
}
