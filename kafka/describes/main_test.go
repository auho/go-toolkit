package describes

import (
	"os"
	"testing"
)

var _network = "tcp"
var _address = "127.0.0.1:9092"

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setUp() {

}

func tearDown() {

}
