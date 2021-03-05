package processor

import (
	"database/sql"
	"log"
	"time"
)

// Collection holds a collection of Hosts for further processing.
type Collection struct {
	Hosts []Host
}

type HostProcess struct {
	Host     *Host
	Priority uint32
	Tubes    []string
}

// Host struct holds valuable data from a host (one domain).
type Host struct {
	Hostname string
	Source   string
	Date     time.Time
	JobType  string
	Retries  uint16
}

func ProcessHost(currentHost *Host, jobID int, db *sql.DB) (HostProcess, error) {

	hostProc := HostProcess{
		currentHost,
		10,
		[]string{},
	}

	switch currentHost.JobType {
	case "deduplication":
		inserted, err := ProcessDeduplication(currentHost)
		if err != nil {
			return hostProc, err
		}
		if inserted {
			log.Printf("[dedup] processed %v: %v\n", jobID, currentHost.Hostname)
			hostProc.Priority = 15
			hostProc.Tubes = []string{"opengit"}
			return hostProc, nil

		}

	case "opengit":
		_, _, err := ProcessOpenGit(currentHost)
		if err != nil {
			return hostProc, err
		}
		log.Printf("[opengit] processed %v: %v\n", jobID, currentHost.Hostname)
	}

	return hostProc, nil

}
