package host

import "time"

type HostArray struct {
	Hosts []Host
}

type Host struct {
	Hostname string
	Source   string
	Date     time.Time
}

type HostProtocol interface {
	Decode(encodedHost []byte) (*Host, error)
	Encode(host *Host) ([]byte, error)
}

type HostProcessor interface {
	DoProcess(host *Host) (bool, error)
}
