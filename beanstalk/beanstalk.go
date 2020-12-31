package beanstalk

import (
	"errors"
	"log"
	"time"

	"github.com/iwanbk/gobeanstalk"
)

const (
	DefaultDelay time.Duration = 0
	DelayOnError time.Duration = 12 * time.Hour
	MaxRetries   uint16        = 3
)

type Handler struct {
	ServerAddress    string
	serverConnection *gobeanstalk.Conn
}

// Connect to the Beanstalk instance
func (bs *Handler) Connect() error {

	// try the connection three times before aborting
	for i := 1; i <= 3; i++ {
		beanstalkConnection, err := gobeanstalk.Dial(bs.ServerAddress)
		if err != nil {
			log.Printf("%v (Retry %d from %d)\n", err, i, 3)
			time.Sleep(time.Second * 5)
		} else {
			bs.serverConnection = beanstalkConnection
			return nil
		}
	}

	return errors.New("connection could not be established")

}

// Close the open Beanstalk connection
func (bs *Handler) Close() {
	if bs.serverConnection != nil {
		bs.serverConnection.Quit()
	}
}

func GetDefaultDelay() time.Duration {
	return DefaultDelay
}

// NewProducer returns a Beanstalk Producer
func NewHandler(serverAdress string) *Handler {
	bs := Handler{ServerAddress: serverAdress}
	return &bs
}
