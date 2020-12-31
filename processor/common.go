package processor

import "time"

// Collection holds a collection of Hosts for further processing.
type Collection struct {
	Hosts []Host
}

// Host struct holds valuable data from a host (one domain).
type Host struct {
	Hostname string
	Source   string
	Date     time.Time
	JobType  string
}
