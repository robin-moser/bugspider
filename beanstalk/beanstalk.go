package beanstalk

import (
	"errors"
	"log"
	"time"

	"github.com/iwanbk/gobeanstalk"
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

// NewProducer returns a Beanstalk Producer
func NewHandler(serverAdress string) *Handler {
	bs := Handler{ServerAddress: serverAdress}
	return &bs
}
