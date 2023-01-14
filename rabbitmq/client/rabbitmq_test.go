package client

import "fmt"

var dsn = fmt.Sprintf("amqp://%s:%s@%s:%d/%s", "guest", "guest", "127.0.0.1", 5672, "guest")
