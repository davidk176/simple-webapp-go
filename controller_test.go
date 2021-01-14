package main

import "testing"

func TestSum(t *testing.T) {
	s := 1 + 1
	if s != 2 {
		t.Error("Test failed")
	}
}
