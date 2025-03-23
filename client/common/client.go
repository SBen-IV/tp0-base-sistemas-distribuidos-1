package common

import (
	"bufio"
	"net"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
	BetsFilePath  string
}

// 
type Client interface {
	Connect() error
	Send(buffer []byte, size int) (int, error)
	Recv(buffer []byte, size int) (int, error)
	Stop()
}


// client Entity that encapsulates connection to the server.
type client struct {
	config ClientConfig
	conn   net.Conn
}

// CreateClient Initializes a new client receiving the configuration
// as a parameter
func CreateClient(config ClientConfig) *client {
	client := &client{
		config: config,
	}

	return client
}

// CreateClientSocket Initializes client socket (conn). In case of
// failure, error is printed in stdout/stderr and error is returned
func (c *client) Connect() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return err
	}
	c.conn = conn
	return nil
}

func (c *client) Send(buffer []byte, size int) (int, error) {
	var total_bytes_sent int = 0

	for total_bytes_sent < size {
		writer := bufio.NewWriterSize(c.conn, (size - total_bytes_sent))
		bytes_sent, err := writer.Write(buffer[total_bytes_sent:])
		
		if err != nil {
			return total_bytes_sent, err
		}

		if err := writer.Flush(); err != nil {
			return total_bytes_sent, err
		}

		total_bytes_sent += bytes_sent
	}

	return total_bytes_sent, nil
}

func (c *client) Recv(buffer []byte, size int) (int, error) {
	var total_bytes_recv int = 0 

	for total_bytes_recv < size {
		bytes_recv, err := bufio.NewReaderSize(c.conn, (size - total_bytes_recv)).Read(buffer[total_bytes_recv:])

		if err != nil {
			return total_bytes_recv, err
		}

		total_bytes_recv += bytes_recv
	}

	return total_bytes_recv, nil
}

// Handle close connection
func (c *client) Stop() {
	if c.conn != nil {
		c.conn.Close()
		log.Info("Client socket closed")
	}
}
