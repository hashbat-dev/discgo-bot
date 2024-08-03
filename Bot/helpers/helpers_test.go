package helpers

import (
	"testing"
)

func TestGetRandomInt(t *testing.T) {
	inputSlice := []int{1, 2, 3, 4, 5}
	n := GetRandomInt(inputSlice)
	var found bool
	for _, v := range inputSlice {
		if n == v {
			found = true
		}
	}
	if found != true {
		t.Errorf("expected %v, got %v", true, found)
	}
}
