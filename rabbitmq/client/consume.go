package client

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"sync"
	"time"
)

// Consumer 监测 connection 状态，连接断开或异常，重新连接；
// 在 Consumer.StartConsume 中连接 connection，并启动 consumeChannel 相关任务，
// 之后 Consumer 只负责监测 connection 状态
//
// consumeChannel 监测 channel 状态，channel 连接断开或异常，重新连接 channel 并重启 consumeChannel 相关任务
//

var (
	consumerReConnectBackoff           = time.Second      // 重连避退时长
	consumerReConnectMaxBackoff        = time.Second * 30 // 重连最大避退时长
	defaultMessageHandlerConcurrent    = 1                // 默认并发数
	defaultMessageHandlerMaxConcurrent = 50               // 最大并发数，进行限制，防止出现并发数设置过大
)

// ConsumeQueue 消费队列的结构体（队列的设置以及消息处理器）
type ConsumeQueue struct {
	// queue 相关
	Name       string // 队列名称
	Durable    bool   // 是否持久化
	AutoDelete bool   // 是否自动删除
	Exclusive  bool   // 是否 exclusive

	// consume 相关
	PrefetchCount int // 同一时间处理的最大消息数
	PrefetchSize  int // 预取消息大小，通常为默认值 0 即可

	// 处理消息相关
	MaxMessageConcurrent int                   // 消息处理器最大并发数
	MessageHandler       ConsumeMessageHandler // 消息处理器
}

// ConsumeMessageHandler 消息处理器
// 返回 bool，true 表示消息处理成功，false 表示消息处理失败
// 不需要手动执行 amqp.Delivery 的 Ack/nAck 方法，
// Consumer 会统一处理（根据 ConsumeMessageHandler 返回的 bool 进行判断处理）
type ConsumeMessageHandler func(amqp.Delivery) bool

// ConsumerConfig 消费者配置
// rabbitmq 链接 amqp://[username]:[password]@[address]:[port]/[vhost]
// 示例：amqp://guest:guest@127.0.0.1:5762/guest
type ConsumerConfig struct {
	Dsn  string // rabbitmq 链接
	Name string // 消费者名称，便于区分、记录日志
}

// Consumer 消费者
// 一个消费者下可以有多个 channel，每个 channel 上可以绑定多个队列消费，
// 每个绑定的队列有队列自己的消息处理器 ConsumeMessageHandler
// 通过 HandleChannel 注册 channel、队列以及队列的消息处理器
type Consumer struct {
	connection
	name             string
	channels         map[string]*consumeChannel // Key：channel 的名称；
	shutDownChan     chan struct{}
	shutDownDoneChan chan struct{}
}

func NewConsumer(config ConsumerConfig) *Consumer {
	c := &Consumer{}
	c.dsn = config.Dsn
	c.name = config.Name

	c.channels = map[string]*consumeChannel{}
	c.shutDownChan = make(chan struct{})
	c.shutDownDoneChan = make(chan struct{})

	if config.Name == "" {
		c.name = fmt.Sprintf("%x", time.Now().UnixNano())
	}

	return c
}

// HandleChannel 先打开一个 channel，并在 channel 上绑定消费队列以及消费队列的消息处理器
// name：此消费者 channel 的名称，自定义名称，代表业务的名称，不要与其他业务重复
// []ConsumeQueue：消费队列（消费队列的队列的设置以及消息处理器）
func (c *Consumer) HandleChannel(name string, cqs []ConsumeQueue) error {
	if _, ok := c.channels[name]; ok {
		return errors.New(fmt.Sprintf("name[%s] of channel is exists", name))
	}

	c.channels[name] = newChannel(c, name, cqs)

	return nil
}

// ShutDown 关闭 Consumer
// 先发送 shutdown 通知，然后等待 shutdown 安全关闭
func (c *Consumer) ShutDown() {
	c.logf("shutdown ...")
	close(c.shutDownChan)
	<-c.shutDownDoneChan
	c.logf("shutdown done.")
}

// StartConsume 开启消息消费，
// 如果返回 error，表示开启失败，并且不会自动重新开启消费，需要检查 error
func (c *Consumer) StartConsume() error {
	err := c.dial()
	if err != nil {
		return err
	}

	for _, cc := range c.channels {
		err = cc.consume()
		if err != nil {
			return err
		}
	}

	go c.notify()

	return nil
}

func (c *Consumer) reConsume() {
	b := consumerReConnectBackoff

	for {
		select {
		case <-c.shutDownChan:
			c.logf("reConsumer shutDown")
			c.channelsShutDown()
			c.shutDownDone()

			return
		default:
			err := c.dial()
			if err != nil {
				c.logf("reconnect %v", err)
				b = backoff(b, consumerReConnectMaxBackoff)
			} else {
				goto LOOP
			}
		}
	}
LOOP:

	c.logf("reConsume")

	go c.notify()
}

func (c *Consumer) notify() {
	state := notifyStateContinue

	// 通过 shutdown 来关闭 Consumer
	// 通过 notify close chan 触发 Consumer 重连以及重启 channels
	select {
	case err := <-c.notifyCloseChan:
		c.logf("notify close %v", err)
	case <-c.shutDownChan:
		c.logf("notify shutDown")
		state = notifyStateShutdown

		// 优先关闭 channels，在关闭 Consumer(也就是 consumer)
		// 防止出现 channel 未关闭完全，Consumer 就提前关闭了
		c.channelsShutDown()

		if !c.isClosed() {
			err := c.close()
			if err != nil {
				c.logf("close %v", err)
			}
		}

		break
	}

	if isNotifyStateShutDown(state) {
		c.shutDownDone()

		return
	}

	c.reConsume()
}

func (c *Consumer) shutDownDone() {
	c.logf("shutdown end")
	close(c.shutDownDoneChan)
}

func (c *Consumer) channelsShutDown() {
	for _, cc := range c.channels {
		cc.shutDown()
	}
}

func (c *Consumer) logf(format string, a ...interface{}) {
	fmt.Printf("Consumer["+c.name+"] "+format+"\n", a...)
}

// consumeChannel 消费者下的channel
// 一个 channel 可以同时消费多个 queue
//
// 先通过 shutDownChan 通知关闭整个 channel，包括：取消 consume，关闭 channel
// shutdown 操作完成后，会通知 shutDownDoneChan，表示切底的关闭了 channel
type consumeChannel struct {
	sw               sync.WaitGroup
	name             string
	consumer         *Consumer
	amqpChannel      *amqp.Channel
	queues           []ConsumeQueue
	notifyCloseChan  chan *amqp.Error
	shutDownChan     chan struct{} // 通知 channel shutDown
	shutDownDoneChan chan struct{} // 彻底关闭的通知
}

func newChannel(consumer *Consumer, name string, cqs []ConsumeQueue) *consumeChannel {
	c := &consumeChannel{}
	c.consumer = consumer
	c.name = name
	c.queues = cqs

	c.shutDownChan = make(chan struct{})
	c.shutDownDoneChan = make(chan struct{})

	for k, q := range cqs {
		if q.MaxMessageConcurrent <= 0 {
			q.MaxMessageConcurrent = defaultMessageHandlerConcurrent
		} else if q.MaxMessageConcurrent > defaultMessageHandlerMaxConcurrent {
			q.MaxMessageConcurrent = defaultMessageHandlerMaxConcurrent
		}

		cqs[k] = q
	}

	return c
}

func (c *consumeChannel) consume() error {
	err := c.connect()
	if err != nil {
		return err
	}

	go c.notify()

	return nil
}

func (c *consumeChannel) reConsume() {
	b := consumerReConnectBackoff

	for {
		select {
		case <-c.shutDownChan:
			c.logf("reConsumer shutDown")
			c.shutDownDone()

			return
		default:
			err := c.connect()
			if err != nil {
				c.logf("reconnect %s", err)
				b = backoff(b, consumerReConnectMaxBackoff)
			} else {
				goto LOOP
			}
		}
	}
LOOP:

	c.logf("reConsume")

	go c.notify()
}

func (c *consumeChannel) notify() {
	state := notifyStateContinue

	select {
	// 通过 shutdown 关闭 channel
	// 通过 notify close chan 触发 channel 重启以及绑定在 channel 上的消费任务
	case err := <-c.notifyCloseChan:
		c.logf("notify close %v", err)

		break
	case <-c.shutDownChan:
		c.logf("notify shutdown")
		state = notifyStateShutdown

		err := c.cancel()
		if err != nil {
			c.logf("cancel %v", err)
		}

		err = c.amqpChannel.Close()
		if err != nil {
			c.logf("close %v", err)
		}

		break
	}

	// 等等消息处理器处理完消息
	c.sw.Wait()

	// shutdown 发送 shutDownDone 通知
	if isNotifyStateShutDown(state) {
		c.shutDownDone()

		return
	}

	c.reConsume()
}

func (c *consumeChannel) connect() error {
	amqpChannel, err := c.consumer.openChannel()
	if err != nil {
		return errors.New(fmt.Sprintf("consumeChannel[%s] open channel %s", c.name, err))
	}

	c.amqpChannel = amqpChannel
	c.notifyCloseChan = c.amqpChannel.NotifyClose(make(chan *amqp.Error))

	for _, q := range c.queues {
		_, err = c.amqpChannel.QueueDeclare(q.Name, q.Durable, q.AutoDelete, q.Exclusive, false, nil)
		if err != nil {
			return err
		}

		err = c.consumeQueue(q, q.MessageHandler)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *consumeChannel) consumeQueue(queue ConsumeQueue, handler ConsumeMessageHandler) error {
	err := c.amqpChannel.Qos(queue.PrefetchCount, queue.PrefetchSize, false)
	if err != nil {
		return err
	}

	deliveryChan, err := c.amqpChannel.Consume(queue.Name, c.consumerTag(queue.Name), false, false, false, true, nil)
	if err != nil {
		return err
	}

	var f func()
	f = func() {
		defer func() {
			// TODO 优化：补救措施，因一些意外情况，消息处理器 goroutine 退出
			if r := recover(); r != nil {
				c.logf("queue[%s] panic %v", queue.Name, r)
				f()
			}
		}()

		// TODO 优化：防止 handler 超时问题
		var err error
		for d := range deliveryChan {
			ok := handler(d)
			if ok {
				err = d.Ack(false)
				if err != nil {
					c.logf("queue[%s] ack: %v", queue.Name, err)
				}
			} else {
				err = d.Nack(false, true)
				if err != nil {
					c.logf("queue[%s] nack: %v", queue.Name, err)
				}
			}
		}

		c.logf("queue[%s] done", queue.Name)
		c.sw.Done()
	}

	for i := 0; i < queue.MaxMessageConcurrent; i++ {
		c.sw.Add(1)
		go f()
	}

	return nil
}

func (c *consumeChannel) consumerTag(queue string) string {
	return fmt.Sprintf("%s-%s", c.name, queue)
}

func (c *consumeChannel) cancel() error {
	var err error
	for _, q := range c.queues {
		e := c.amqpChannel.Cancel(c.consumerTag(q.Name), false)
		if e != nil {
			err = e
		}
	}

	return err
}

// shutdown 发送 shutdown 通知
// 等待所有 shutdown 完成，接受到 shutDownDoneChan 通知，就表示切底的安全关闭了
// 通过 close chan 防止发生 chan 阻塞
func (c *consumeChannel) shutDown() {
	c.logf("shutdown ...")
	close(c.shutDownChan)
	<-c.shutDownDoneChan
	c.logf("shutdown done")
}

// shutDownDone 所有需要 shutdown 的已经都完成，发送此状态通知
// 通过 close chan 防止发生 chan 阻塞
func (c *consumeChannel) shutDownDone() {
	c.logf("shutdown end")
	close(c.shutDownDoneChan)
}

func (c *consumeChannel) logf(format string, a ...interface{}) {
	fmt.Printf(fmt.Sprintf("ConsumeChannel[%s][%s] %s\n", c.consumer.name, c.name, format), a...)
}
