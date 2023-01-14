package client

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

// Consumer 使用示例
func ExampleConsumer() {
	var dsn = fmt.Sprintf("amqp://%s:%s@%s:%d/%s", "guest", "guest", "127.0.0.1", 5672, "guest")

	// 队列 queue 的消息处理器
	var handler = func(delivery amqp.Delivery) bool {
		log.Println("msg is ", string(delivery.Body), delivery.DeliveryTag)
		return true
	}

	// 初始化 Consumer
	c := NewConsumer(ConsumerConfig{Dsn: dsn, Name: "example"})

	// 定义一个 channel， 并注入队列 queue 以及 queue 的处理器
	err := c.HandleChannel("one", []ConsumeQueue{
		{
			Name:           "testOne",
			Durable:        true,
			MessageHandler: handler,
			PrefetchCount:  200,
		},
		{
			Name:           "testTwo",
			Durable:        true,
			MessageHandler: handler,
			PrefetchCount:  100,
		},
		{
			Name:           "testThree",
			Durable:        true,
			MessageHandler: handler,
			PrefetchCount:  100,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// 开始消费队列 queue 的消息
	// 如果返回 error，需要检查 error
	err = c.StartConsume()
	if err != nil {
		log.Fatal(err)
	}

	// 关闭 Consumer，停止消费消息
	c.ShutDown()
}

// Publisher 示例
func ExamplePublisher_Publish() {
	var dsn = fmt.Sprintf("amqp://%s:%s@%s:%d/%s", "guest", "guest", "127.0.0.1", 5672, "guest")

	// 初始化 Publisher
	p, err := NewPublisher(PublisherConfig{Dsn: dsn, Name: "publisher"})
	if err != nil {
		log.Fatal(err)
	}

	// 发布消息
	err = p.Publish(&Message{
		Exchange: "",
		Key:      "route key",
		Publishing: amqp.Publishing{
			ContentType: ContentTypePlain,
			Body:        []byte(fmt.Sprintf("message body")),
		},
		// 服务器是否接受并处理消息，业务根据结果进行后续处理
		ConfirmHandler: func(b bool) {
			if b {
				// 发送成功的处理
			} else {
				// 发送失败的处理
			}
		},
	})

	// 因为使用的是 confirm 模式，err 表示发送动作执行是否成功
	// 消息是否被服务其接受并处理，需要在 ConfirmHandler 中判断并处理
	if err != nil {
		fmt.Println("publish:", err)
	}

	// 关闭 Publisher，确保消息都发送出去并被服务器接受处理，以及接受到服务器返回的处理结果
	p.ShutDown()
}
