package helpers

import (
	"testing"
)

func TestRamdonStr(t *testing.T) {
	s, err := GenerateRandomString(20)
	if err != nil {
		t.Error(err)
	}
	t.Log(s)
}

func TestGenerateAppID(t *testing.T) {
	s := GenerateRandAppID(20)
	if len(s) == 0 {
		t.Error("can not generate random appid")
	}
	t.Log(s)
}
