package client

import (
	"github.com/streadway/amqp"
	"sync"
)

type connection struct {
	dsn             string
	amqpConnection  *amqp.Connection
	notifyCloseChan chan *amqp.Error
	amqpMutex       sync.Mutex
}

func (c *connection) dial() error {
	amqpConnection, err := amqp.Dial(c.dsn)
	if err != nil {
		return err
	}

	// 有可能会和 connection.openChannel 产生 data race
	c.amqpMutex.Lock()
	defer c.amqpMutex.Unlock()

	c.amqpConnection = amqpConnection
	c.notifyCloseChan = c.amqpConnection.NotifyClose(make(chan *amqp.Error))

	return nil
}

func (c *connection) openChannel() (*amqp.Channel, error) {
	// 有可能会和 connection.dial 产生 data race
	c.amqpMutex.Lock()
	defer c.amqpMutex.Unlock()

	return c.amqpConnection.Channel()
}

func (c *connection) isClosed() bool {
	return c.amqpConnection.IsClosed()
}

func (c *connection) close() error {
	return c.amqpConnection.Close()
}
