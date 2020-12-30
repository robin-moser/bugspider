package beanstalk

import (
	"time"

	"github.com/robin-moser/bugspider/host"
)

type Producer struct {
	MainBeanstalk
	protocol host.Protocol
}

// PutHost pushes the given Host to beanstalk
func (producer *Producer) PutHost(host *host.Host) error {
	body, err := producer.protocol.Encode(host)
	if err != nil {
		return err
	}
	priority := uint32(10)
	delay := 0 * time.Second
	timeToRun := 20 * time.Second
	_, err = producer.serverConnection.Put(body, priority, delay, timeToRun)
	if err != nil {
		return err
	}
	return nil
}

// MakeNewProducer returns a Beanstalk Producer
func MakeNewProducer(serverAdress string, protocol host.Protocol) *Producer {
	producer := Producer{protocol: protocol}
	producer.ServerAddress = serverAdress
	return &producer
}
