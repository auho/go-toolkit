package source

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/redis/client"
	"github.com/go-redis/redis/v8"
)

var _redisOptions = redis.Options{
	Network:  "tcp",
	Addr:     "127.0.0.1:6379",
	Password: "123456",
}

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setUp() {
	log.Println("set up")
	rand.Seed(time.Now().UnixNano())
}

func tearDown() {
	log.Println("tear down")
}

func _testKey[E storage.Entry](
	t *testing.T,
	key string,
	bFunc func(config Config) (*key[E], error),
	lFunc func(ctx context.Context, c *client.Redis) (int64, error),
) {
	ctx := context.Background()

	k, err := bFunc(Config{
		Concurrency: 1,
		Amount:      0,
		PageSize:    0,
		Key:         key,
		Options:     &_redisOptions,
	})

	if err != nil {
		t.Fatal(err)
	}

	if err != nil {
		t.Fatal("new", err)
	}

	err = k.Scan()
	if err != nil {
		t.Fatal("scan", err)
	}

	amount := 0
	for items := range k.ReceiveChan() {
		l := len(items)
		amount = amount + l
	}

	fmt.Println(k.Summary())
	fmt.Println(k.State())

	if k.total != k.state.Amount() || k.state.Amount() != int64(amount) {
		t.Error(fmt.Sprintf("total != amount != actual %d != %d != %d", k.total, k.state.Amount(), amount))
	}

	dbAmount, err := lFunc(ctx, k.GetClient())
	if err != nil {
		t.Error("db amount ", err)
	}

	if k.total != dbAmount {
		t.Error(fmt.Sprintf("total != db amount %d != %d", k.total, dbAmount))
	}
}

func _randAmount() int {
	i := int(10e3)
	i += rand.Intn(1000)
	return i
}
