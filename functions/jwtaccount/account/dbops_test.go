package account

import (
	"testing"
)

func TestInitDefaultAccount(t *testing.T) {
	err := CheckDefaultAppAccount()
	if err != nil {
		t.Error(err)
	}
	t.Log("success init default account")
}

func TestUpdateDefaultHashedPasswd(t *testing.T) {
	_, passwd, err := UpdateDefaultSecurityKey()
	if err != nil {
		t.Error(err)
	}
	t.Log("success update default password", passwd)
}
