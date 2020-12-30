package host

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
}

// Protocol defines Methods for decoding and encoding Hosts between type Host an stings.
// This is needed as Beanstalk stores jobs as strings.
type Protocol interface {
	Decode(encodedHost []byte) (*Host, error)
	Encode(host *Host) ([]byte, error)
}

// Processor defines Methods for further processing Hosts, for example stroring Hosts.
type Processor interface {
	DoProcess(host *Host) (bool, error)
}
