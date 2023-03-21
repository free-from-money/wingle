package wingle

import (
	"bytes"
	"github.com/imroc/req/v3"
	"log"
	"strings"
	"sync"
	"time"
)

const (
	host    = `192.168.8.1`
	baseUrl = `http://` + host

	sessionPath = `/api/webserver/SesTokInfo`
	rebootPath  = `/api/device/control`
	statusPath  = `/api/monitoring/status`
	loginPath   = `/api/user/login`

	tokenHeader        = `__RequestVerificationToken`
	rebootPayload      = `<?xml version="1.0" encoding="UTF-8"?><request><Control>1</Control></request>`
	defaultContentType = `text/xml`
)

var (
	client = req.C().
		SetBaseURL(baseUrl).
		SetTimeout(time.Second).
		SetCommonRetryCount(-1)

	once = sync.Once{}
)

func init() {
	once.Do(func() {
		if !checkWingleIsConnected() {
			log.Fatal("wingle is not connected")
		}
	})
}

func checkWingleIsConnected() bool {
	res := client.Get().SetRetryCount(0).Do()
	return res.IsSuccessState()
}

func ChangeIp() {
	currentIp := GetIp()
	log.Printf("현재 IP={%s}", currentIp)
	for {
		RebootRouter()
		time.Sleep(5 * time.Second)
		changedIp := GetIp()
		log.Printf("변경 후 IP={%s}", changedIp)
		if currentIp == changedIp {
			continue
		}
		break
	}
}

func checkStatus() bool {
	get, err := client.R().
		SetRetryCount(0).
		Get(statusPath)
	if err != nil {
		return false
	} else if get.StatusCode != 200 {
		return false
	}

	return true
}

func getSessionToken() SesTok {
	var st SesTok
	res, _ := client.R().
		SetRetryCount(0).
		Get(sessionPath)

	res.Header.Set(`Content-Type`, defaultContentType)
	res.Unmarshal(&st)

	return st
}

func Login() SesTok {
	st := getSessionToken()
	payload := NewLoginRequest(st.TokInfo)
	res := client.R().
		SetRetryCount(0).
		SetCookies(st.GetCookie()).
		SetHeader(tokenHeader, st.TokInfo).
		SetContentType(defaultContentType).
		SetBodyXmlMarshal(&payload).
		MustPost(loginPath)

	defer res.Body.Close()

	return SesTok{
		SesInfo: strings.Split(res.GetHeader(`Set-Cookie`), ";")[0],
		TokInfo: res.GetHeader(`__requestverificationtokenone`),
	}
}

func RebootRouter() {
	st := Login()
	defaultHeaders := map[string]string{
		"Host":             host,
		"Origin":           baseUrl,
		"Referer":          baseUrl + `/html/reboot.html`,
		"X-Requested-With": `XMLHTTPRequest`,
	}
	res := client.R().
		DisableAutoReadResponse().
		DisableTrace().
		SetCookies(st.GetCookie()).
		SetHeaders(defaultHeaders).
		SetHeader(tokenHeader, st.TokInfo).
		SetContentType(`application/x-www-form-urlencoded; charset=UTF-8`).
		SetBodyBytes(bytes.NewBufferString(rebootPayload).Bytes()).
		SetRetryCount(0).
		EnableCloseConnection().
		MustPost(rebootPath)

	res.Header.Set(`Content-Type`, defaultContentType)

	var result string
	res.Unmarshal(&result)

	log.Println(result)
}
