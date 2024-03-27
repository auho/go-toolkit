package describes

import (
	"os"
	"testing"
)

var _address = "http://127.0.0.1:9200"

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
