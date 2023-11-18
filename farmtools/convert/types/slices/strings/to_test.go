package strings

import (
	"fmt"
	"testing"
)

func TestSpiltToSliceString(t *testing.T) {
	s := `1,2,3,"4","5",""`

	ss, err := SpiltToSliceString(s, ",")
	if err != nil {
		t.Fatal(err)
	}

	if len(ss) != 6 {
		t.Error("len")
	}

	fmt.Println(ss)
}

func TestSpiltToSliceInt(t *testing.T) {
	s := `1,2,3,4,5,`

	ss, err := SpiltToSliceInt(s, ",")
	if err != nil {
		t.Fatal(err)
	}

	if len(ss) != 6 {
		t.Error("len")
	}

	fmt.Println(ss)
}
