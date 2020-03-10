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
