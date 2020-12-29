package beanstalk

import (
	"fmt"
	"os"

	"github.com/iwanbk/gobeanstalk"
)

type MainBeanstalk struct {
	ServerAddress    string
	serverConnection *gobeanstalk.Conn
}

func (bs *MainBeanstalk) Connect() {
	beanstalkConnection, err := gobeanstalk.Dial(bs.ServerAddress)
	if err != nil {
		// TODO: do retries
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("connection established")
	bs.serverConnection = beanstalkConnection
}

func (bs *MainBeanstalk) Close() {
	if bs.serverConnection != nil {
		bs.serverConnection.Quit()
	}
	fmt.Println("connection closed")
}
