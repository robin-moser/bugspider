package beanstalk

import (
	"encoding/json"

	"github.com/robin-moser/bugspider/host"
)

type JSONHostProtocol struct {
}

func (protocol *JSONHostProtocol) Decode(encodedHost []byte) (*host.Host, error) {
	unCodedHost := host.Host{}
	err := json.Unmarshal(encodedHost, &unCodedHost)
	if err != nil {
		return nil, err
	}
	return &unCodedHost, nil
}

func (protocol *JSONHostProtocol) Encode(host *host.Host) ([]byte, error) {
	encodedHost, err := json.Marshal(host)
	if err != nil {
		return nil, err
	}
	return encodedHost, nil
}

func MakeJSONHostProtocol() *JSONHostProtocol {
	return &JSONHostProtocol{}
}
