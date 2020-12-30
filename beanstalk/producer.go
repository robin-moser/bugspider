package beanstalk

import (
	"encoding/json"
	"time"

	"github.com/robin-moser/bugspider/host"
)

type Producer struct {
	MainBeanstalk
}

// PutHost pushes the given Host to beanstalk
func (producer *Producer) PutHost(currentHost *host.Host) error {
	body, err := json.Marshal(currentHost)
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
func MakeNewProducer(serverAdress string) *Producer {
	producer := Producer{}
	producer.ServerAddress = serverAdress
	return &producer
}
