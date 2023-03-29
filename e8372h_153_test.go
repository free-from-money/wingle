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
	rebootHuawei()
}

func TestCheckStatus(t *testing.T) {
	t.Log(isHuawei())
}

func TestLogin(t *testing.T) {
	Login()
}

func TestRebootRouter(t *testing.T) {
	RebootRouter()
}

func Test_checkWingleIsConnected(t *testing.T) {
	t.Log(checkWingleIsConnected())
}
