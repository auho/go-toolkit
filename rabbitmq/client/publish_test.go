package client

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
	"testing"
	"time"
)

var queues = []string{"testOne", "testTwo", "testThree"}
var p *Publisher

func setup() {
	var err error
	p, err = NewPublisher(PublisherConfig{Dsn: dsn})
	if err != nil {
		log.Fatal(err)
	}
}

func tearDown() {
	p.ShutDown()

	if p.channel.confirms.published != p.channel.confirms.confirmed {
		log.Println(p.channel.confirms.messageMap)

		log.Fatalf("acutal %d expect %d", p.channel.confirms.published, p.channel.confirms.confirmed)
	}

	time.Sleep(time.Second * 1)
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func Test_publish(t *testing.T) {
	var err error
	var ticker1, ticker2 *time.Ticker
	closeChan := make(chan struct{})

	go func() {
		ticker1 = time.NewTicker(time.Second)

		for {
			select {
			case <-ticker1.C:
				publish()
			case <-closeChan:
				ticker1.Stop()
				log.Println("ticker 1 close")
				return
			}
		}
	}()

	go func() {
		ticker2 = time.NewTicker(time.Second * 5)

		for {
			select {
			case <-ticker2.C:
				err = p.Publish(&Message{
					Exchange: "abc",
					Key:      "testFour",
					Publishing: amqp.Publishing{
						ContentType: ContentTypePlain,
						Body:        []byte(fmt.Sprintf("from %s publish %d", "", time.Now().Unix())),
					},
					ConfirmHandler: func(b bool) {
						if !b {
							log.Println("ack:", b)
						}
					},
				})

				if err != nil {
					t.Log(err)
				}

			case <-closeChan:
				ticker2.Stop()
				log.Println("ticker 2 close")
				goto LOOP
			}
		}
	LOOP:
	}()

	time.Sleep(time.Second * 20)
	close(closeChan)
	time.Sleep(time.Second * 5)

	err = p.close()
	if err != nil {
		log.Println("publisher consumer:", err)
	}
}

func Benchmark_publishBench(b *testing.B) {
	for i := 0; i < b.N; i++ {
		publish()
	}
}

func publish() {
	var err error
	for _, q := range queues {
		err = p.Publish(&Message{
			Exchange: "",
			Key:      q,
			Publishing: amqp.Publishing{
				ContentType: ContentTypePlain,
				Body:        []byte(fmt.Sprintf("from %s publish %d", q, time.Now().Unix())),
			},
			ConfirmHandler: func(b bool) {
				if !b {
					log.Println("ack:", b)
				}
			},
		})

		if err != nil {
			fmt.Println("publish:", err)
		}
	}
}
