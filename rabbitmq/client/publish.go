package client

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var (
	publisherReConnectBackoff          = time.Millisecond * 100 // 重连避退时长
	publisherReConnectMaxBackoff       = time.Second * 10       // 重连最大避退时长
	defaultPublishConfirmConcurrent    = 2                      // 默认并发数
	defaultPublishConfirmMaxConcurrent = 100                    // 最大并发数，进行限制，防止出现并发数设置过大

	errShutDowned = errors.New("publisher is shutdown")
)

// PublisherConfig 消息发布者配置
// rabbitmq Dsn 链接: amqp://[username]:[password]@[address]:[port]/[vhost]
// 示例：amqp://guest:guest@127.0.0.1:5762/guest
type PublisherConfig struct {
	Dsn                  string // rabbitmq 链接
	Name                 string // 发布者名称，便于区分，记录日志
	MaxPublishConcurrent int    // 最大发布消息发布并发数
	MaxConfirmConcurrent int    // 最大消息确认并发数
}

// PublishHandler 发送消息后，rabbitmq 服务会返回客户端消息是否接收处理成功
// 客户端根据返回的结果进行后续的业务处理
// bool：true 接收处理成功，false 失败
type PublishHandler func(bool)

// Message 消息结构体
type Message struct {
	Exchange       string          // 交换机名称
	Key            string          // 路由 key 或者 queue 名称
	Mandatory      bool            // mandatory
	Immediate      bool            // immediate
	Publishing     amqp.Publishing // 消息内容
	ConfirmHandler PublishHandler  // 消息是否发送成功的后续业务处理
}

// Publisher 消息发布者
// 使用 rabbitmq 的 confirm 模式发送消息
type Publisher struct {
	connection
	name             string
	channel          *publishChannel
	shutDownChan     chan struct{}
	shutDownDoneChan chan struct{}
	isShutDowned     uint32
}

func NewPublisher(config PublisherConfig) (*Publisher, error) {
	p := &Publisher{}
	p.dsn = config.Dsn
	p.shutDownChan = make(chan struct{})
	p.shutDownDoneChan = make(chan struct{})

	if config.Name == "" {
		p.name = fmt.Sprintf("%x", time.Now().UnixNano())
	}

	err := p.dial()
	if err != nil {
		return nil, err
	}

	p.channel, err = newPublishChannel(p, config.MaxConfirmConcurrent)
	if err != nil {
		return nil, err
	}

	go p.notify()

	return p, nil
}

// Publish 发布消息
func (p *Publisher) Publish(m *Message) error {
	return p.channel.publish(m)
}

// ShutDown 先发送关闭通知，等待安全的关闭
func (p *Publisher) ShutDown() {
	if !atomic.CompareAndSwapUint32(&p.isShutDowned, 0, 1) {
		return
	}

	p.logf("shutdown ...")
	close(p.shutDownChan)
	<-p.shutDownDoneChan
	p.logf("shutdown done")
}

func (p *Publisher) reConnect() {
	b := publisherReConnectBackoff

	for {
		select {
		case <-p.shutDownChan:
			p.logf("reConnect shutdown")
			p.channel.shutDown()
			p.shutDownDone()

			return
		default:
			err := p.dial()
			if err != nil {
				p.logf("reconnect %v", err)
				b = backoff(b, publisherReConnectMaxBackoff)
			} else {
				goto LOOP
			}
		}
	}
LOOP:

	go p.notify()
}

func (p *Publisher) notify() {
	select {
	case e := <-p.notifyCloseChan:
		p.logf("notify close %v", e)
		break
	case <-p.shutDownChan:
		p.logf("notify shutdown")
		p.channel.shutDown()

		err := p.close()
		if err != nil {
			p.logf("close %v", err)
		}

		p.shutDownDone()

		return
	}

	p.reConnect()
}

func (p *Publisher) shutDownDone() {
	p.logf("shutdown end")
	close(p.shutDownDoneChan)
}

func (p *Publisher) logf(format string, a ...interface{}) {
	fmt.Printf("Publisher["+p.name+"] "+format+"\n", a...)
}

// publishChannel 发布通道
type publishChannel struct {
	publisher         *Publisher
	amqpChannel       *amqp.Channel
	notifyCloseChan   chan *amqp.Error
	notifyPublishChan chan amqp.Confirmation
	confirms          *publishConfirms

	m                    sync.Mutex
	sw                   sync.WaitGroup
	maxConfirmConcurrent int
	shutDownChan         chan struct{}
	shutDownDoneChan     chan struct{}
	isShutDowned         bool
}

func newPublishChannel(p *Publisher, maxConfirm int) (*publishChannel, error) {
	c := &publishChannel{}
	c.publisher = p
	c.maxConfirmConcurrent = maxConfirm
	c.shutDownChan = make(chan struct{})
	c.shutDownDoneChan = make(chan struct{})

	if c.maxConfirmConcurrent <= 0 {
		c.maxConfirmConcurrent = defaultPublishConfirmConcurrent
	} else if c.maxConfirmConcurrent > defaultPublishConfirmMaxConcurrent {
		c.maxConfirmConcurrent = defaultPublishConfirmMaxConcurrent
	}

	err := c.connect()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *publishChannel) connect() error {
	amqpChannel, err := c.publisher.openChannel()
	if err != nil {
		return err
	}

	err = amqpChannel.Confirm(false)
	if err != nil {
		_ = amqpChannel.Close()
		return err
	}

	c.m.Lock()

	c.amqpChannel = amqpChannel
	c.confirms = &publishConfirms{
		published:  0,
		messageMap: map[uint64]*Message{},
	}

	c.notifyCloseChan = c.amqpChannel.NotifyClose(make(chan *amqp.Error))
	c.notifyPublishChan = c.amqpChannel.NotifyPublish(make(chan amqp.Confirmation, c.maxConfirmConcurrent))

	c.m.Unlock()

	c.notifyPublish()

	go c.notify()

	return nil
}

func (c *publishChannel) reConnect() {
	b := publisherReConnectBackoff

	for {
		select {
		case <-c.shutDownChan:
			c.logf("reConnect shutdown")
			c.shutDownDone()

			return
		default:
			err := c.connect()
			if err != nil {
				c.logf("reConnect %v", err)
				b = backoff(b, publisherReConnectMaxBackoff)
			} else {
				goto LOOP
			}
		}
	}
LOOP:

	c.logf("reConnect")
}

func (c *publishChannel) notify() {
	state := notifyStateContinue

	select {
	case err := <-c.notifyCloseChan:
		c.logf("notify close %v", err)

		break
	case <-c.shutDownChan:
		c.logf("notify shutdown")

		state = notifyStateShutdown

		c.waitNotifyPublishComplete()

		err := c.amqpChannel.Close()
		if err != nil {
			c.logf("close", err)
		}
	}

	c.sw.Wait()
	c.logf("notify publish done")

	if isNotifyStateShutDown(state) {
		c.shutDownDone()
		return
	}

	c.reConnect()
}

func (c *publishChannel) notifyPublish() {
	var f func()
	f = func() {
		defer func() {
			// TODO 待优化；补救措施，因一些意外情况，消息处理器 goroutine 退出
			if r := recover(); r != nil {
				c.logf("notify publish handler", r)
				f()
			}
		}()

		for confirmation := range c.notifyPublishChan {
			c.m.Lock()
			p, ok := c.confirms.confirm(confirmation)
			c.m.Unlock()

			// TODO 待优化，防止超时问题
			if ok {
				p.ConfirmHandler(confirmation.Ack)
			}
		}

		c.sw.Done()
	}

	for i := 0; i < c.maxConfirmConcurrent; i++ {
		c.sw.Add(1)
		go f()
	}
}

// waitNotifyPublishComplete 等待所有 Confirmations 到达并处理完成，才能安全的关闭 channel
func (c *publishChannel) waitNotifyPublishComplete() {
	deadline := time.NewTimer(time.Second * 30)
	t := time.NewTicker(time.Millisecond * 100)
	var ok bool

	for {
		select {
		case <-t.C:
			c.m.Lock()
			ok = c.confirms.complete()
			c.m.Unlock()

			if ok {
				t.Stop()
				return
			}
		case <-deadline.C:
			deadline.Stop()
			return
		}
	}
}

func (c *publishChannel) publish(m *Message) error {
	c.m.Lock()
	defer c.m.Unlock()

	if c.isShutDowned {
		return errors.New("publisher prepare to shutdown")
	}

	err := c.amqpChannel.Publish(m.Exchange, m.Key, m.Mandatory, m.Immediate, m.Publishing)
	if err != nil {
		return err
	}

	c.confirms.publish(m)

	return nil
}

func (c *publishChannel) shutDown() {
	if err := c.shutDowned(); err != nil {
		return
	}

	c.logf("shutdown ...")
	close(c.shutDownChan)
	<-c.shutDownDoneChan
	c.logf("shutdown done")
}

func (c *publishChannel) shutDownDone() {
	c.logf("shutdown end")
	close(c.shutDownDoneChan)
}

func (c *publishChannel) shutDowned() error {
	c.m.Lock()
	defer c.m.Unlock()

	if c.isShutDowned {
		return errShutDowned
	}

	c.isShutDowned = true

	return nil
}

func (c *publishChannel) logf(format string, a ...interface{}) {
	fmt.Printf(fmt.Sprintf("PublishChannel[%s] %s\n", c.publisher.name, format), a...)
}

// publishConfirms rabbitmq 的发送消息 confirm 机制
// 发送给服务端的消息按照发送序号保存，发送序号和服务器返回的序号一一对应
// 客户端收到服务发送的消息 amqp.Confirmation，根据 amqp.Confirmation 确定消息是否真正发送成功
type publishConfirms struct {
	messageMap map[uint64]*Message
	published  uint64
	confirmed  uint64
}

func (pc *publishConfirms) complete() bool {
	return pc.published == pc.confirmed
}

func (pc *publishConfirms) publish(m *Message) uint64 {
	pc.published++

	m.Publishing.MessageId = strconv.FormatInt(int64(pc.published), 10)

	pc.messageMap[pc.published] = m

	return pc.published
}

func (pc *publishConfirms) confirm(confirmation amqp.Confirmation) (*Message, bool) {
	pc.confirmed++

	p, found := pc.messageMap[confirmation.DeliveryTag]
	if found {
		pc.delete(confirmation.DeliveryTag)
	}

	return p, found
}

func (pc *publishConfirms) delete(published uint64) {
	delete(pc.messageMap, published)
}
