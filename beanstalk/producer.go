package beanstalk

import (
	"bugspider/host"
	"time"
)

type Producer struct {
	MainBeanstalk
	protocol host.HostProtocol
}

func (producer *Producer) PutHost(host *host.Host) error {
	body, err := producer.protocol.Encode(host)
	if err != nil {
		return err
	}
	priority := uint32(10)
	delay := 0 * time.Second
	time_to_run := 20 * time.Second
	_, err = producer.serverConnection.Put(body, priority, delay, time_to_run)
	if err != nil {
		return err
	}
	return nil
}

func MakeNewProducer(serverAdress string, protocol host.HostProtocol) *Producer {
	producer := Producer{protocol: protocol}
	producer.ServerAddress = serverAdress
	return &producer
}
