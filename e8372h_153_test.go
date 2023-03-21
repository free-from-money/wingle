package wingle

import (
	"log"
	"testing"
)

func TestChangeIp(t *testing.T) {
	ChangeIp()
}
func TestGetIp(t *testing.T) {
	log.Println(GetIp())
}
func TestReboot(t *testing.T) {
	RebootRouter()
}

func TestCheckStatus(t *testing.T) {
	t.Log(checkStatus())
}

func TestLogin(t *testing.T) {
	Login()
}

func Test_checkWingleIsConnected(t *testing.T) {
	t.Log(checkWingleIsConnected())
}
