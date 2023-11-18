package strings

import (
	"fmt"
	"testing"
)

func TestToSlicesInt(t *testing.T) {
	s := `1,2,3,4,5,`

	ss, err := JsonToSliceInt(s)
	if err != nil {
		t.Fatal(err)
	}

	if len(ss) != 6 {
		t.Error("len")
	}

	fmt.Println(ss)
}

func TestToSlicesString(t *testing.T) {
	s := `1,2,3,"4","5",""`

	ss, err := JsonToSlicesString(s)
	if err != nil {
		t.Fatal(err)
	}

	if len(ss) != 6 {
		t.Error("len")
	}

	fmt.Println(ss)
}
