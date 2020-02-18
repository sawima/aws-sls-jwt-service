package main

import (
	"testing"
)

func TestDBAuth(t *testing.T) {
	_, result := dbAuth("kimademo", "demopwd")
	if result == false {
		t.Error("auth failed")
	}
	t.Log("&&&&")
}
