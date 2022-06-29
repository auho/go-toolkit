package client

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"sync"
	"testing"
	"time"
)

func (c *consumeChannel) publish(queue string) {
	if c.consumer.isClosed() {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	err := c.amqpChannel.Publish("", queue, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(c.name + " " + queue),
	})

	if err != nil {
		log.Printf("test publish error: %s %v \n", c.name, err)
	}
}

func (c *consumeChannel) errorPublish() {
	if c.consumer.isClosed() {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	err := c.amqpChannel.Publish("abc", c.name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(c.name),
	})

	if err != nil {
		log.Printf("test publish error: %s %v \n", c.name, err)
	}
}

var handler = func(delivery amqp.Delivery) bool {
	log.Println("test msg", string(delivery.Body), delivery.DeliveryTag)
	return true
}

func Test_consume(t *testing.T) {
	c := NewConsumer(ConsumerConfig{Dsn: dsn})

	err := c.HandleChannel("one", []ConsumeQueue{
		{
			Name:           "testOne",
			Durable:        true,
			MessageHandler: handler,
			PrefetchCount:  2000,
		},
		{
			Name:           "testTwo",
			Durable:        true,
			MessageHandler: handler,
			PrefetchCount:  1000,
		},
		{
			Name:           "testThree",
			Durable:        true,
			MessageHandler: handler,
			PrefetchCount:  1000,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	err = c.StartConsume()
	if err != nil {
		t.Fatal(err)
	}

	var ticker1, ticker2, ticker3 *time.Ticker

	// 每一段时间，发送一个合法的消息进入队列
	go func() {
		ticker1 = time.NewTicker(time.Second * 1)
		for tt := range ticker1.C {
			log.Println("ticker1 :: ", tt.Unix())

			for _, cc := range c.channels {
				for _, q := range cc.queues {
					cc.publish(q.Name)
				}
			}

			log.Println()
		}
	}()

	// 每一段时间，发送一个不合法的消息进入队列，触发队列 notify close
	go func() {
		ticker2 = time.NewTicker(time.Second * 4)
		for tt := range ticker2.C {
			log.Println("ticker2 :: ", tt.Unix())
			for _, cc := range c.channels {
				cc.errorPublish()
			}
		}
	}()

	go func() {
		// 每一段时间，手动关闭连接，触发 Consumer 的 notify close
		ticker3 = time.NewTicker(time.Second * 7)
		for range ticker3.C {
			log.Println("================close consumer")
			err := c.close()
			if err != nil {
				log.Println("consumer:", err)
			}
		}
	}()

	var sw sync.WaitGroup
	sw.Add(1)
	go func() {
		time.AfterFunc(time.Second*100, func() {
			ticker1.Stop()
			ticker2.Stop()
			ticker3.Stop()

			time.Sleep(time.Second * 10)
			log.Println("================ shutDown")
			c.ShutDown()

			sw.Done()
		})
	}()

	sw.Wait()

	log.Println("all done.")
	time.Sleep(time.Second * 1)
}
