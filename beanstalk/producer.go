package beanstalk

import (
	"encoding/json"
	"time"

	"github.com/robin-moser/bugspider/host"
)

func (bs *Handler) UseTube(tube string) error {
	err := bs.serverConnection.Use(tube)
	if err != nil {
		return err
	}
	return nil
}

// PutHost pushes the given Host to beanstalk
func (bs *Handler) PutHost(currentHost *host.Host) error {
	body, err := json.Marshal(currentHost)
	if err != nil {
		return err
	}
	priority := uint32(10)
	delay := 0 * time.Second
	timeToRun := 20 * time.Second
	_, err = bs.serverConnection.Put(body, priority, delay, timeToRun)
	if err != nil {
		return err
	}
	return nil
}
