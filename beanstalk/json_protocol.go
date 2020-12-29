package beanstalk

import (
	"bugspider/host"
	"encoding/json"
)

type JsonHostProtocol struct {
}

func (protocol *JsonHostProtocol) Decode(encodedHost []byte) (*host.Host, error) {
	unCodedHost := host.Host{}
	err := json.Unmarshal(encodedHost, &unCodedHost)
	if err != nil {
		return nil, err
	}
	return &unCodedHost, nil
}

func (protocol *JsonHostProtocol) Encode(host *host.Host) ([]byte, error) {
	encodedHost, err := json.Marshal(host)
	if err != nil {
		return nil, err
	}
	return encodedHost, nil
}

func MakeJsonHostProtocol() *JsonHostProtocol {
	return &JsonHostProtocol{}
}
