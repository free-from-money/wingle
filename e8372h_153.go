package wingle

import (
	"bytes"
	"github.com/secr3t/req/v3"
	"log"
	"strings"
	"sync"
	"time"
)

const (
	defaultHost = `192.168.8.1`

	// huawei
	sessionPath = `/api/webserver/SesTokInfo`
	rebootPath  = `/api/device/control`
	statusPath  = `/api/monitoring/status`
	loginPath   = `/api/user/login`

	tokenHeader        = `__RequestVerificationToken`
	rebootPayload      = `<?xml version="1.0" encoding="UTF-8"?><request><Control>1</Control></request>`
	defaultContentType = `text/xml`

	// oem
	oemRebootPath    = `/reqproc/proc_post`
	oemContentType   = `application/x-www-form-urlencoded; charset=UTF-8`
	oemRebootPayload = `goformId=REBOOT_DEVICE`
)

type Wingle struct {
	Host         string
	oemClient    *req.Client
	huaweiClient *req.Client
	once         *sync.Once
}

// NewWingle args defaultHost must 0 or 1
func NewWingle(host ...string) *Wingle {
	if len(host) > 1 {
		log.Fatal("defaultHost must be just one")
	}
	if len(host) == 0 {
		return newWingle(defaultHost)
	}
	return newWingle(host[0])
}

func newWingle(host string) *Wingle {
	w := &Wingle{
		Host: host,
		once: &sync.Once{},
	}

	w.oemClient = req.C().
		SetBaseURL(w.baseUrl()).
		SetTimeout(time.Second).
		SetCommonRetryCount(0)

	w.huaweiClient = req.C().
		SetBaseURL(w.baseUrl()).
		SetTimeout(time.Second).
		SetCommonRetryCount(-1)

	w.once.Do(func() {
		if !w.checkWingleIsConnected() {
			log.Fatal("Wingle is not connected")
		}
	})

	return w
}

func (w *Wingle) baseUrl() string {
	return "http://" + w.Host
}

func (w *Wingle) checkWingleIsConnected() bool {
	res := w.huaweiClient.Get().SetRetryCount(0).Do()
	return res.IsSuccessState()
}

func (w *Wingle) ChangeIp() {
	currentIp := GetIp()
	log.Printf("현재 IP={%s}", currentIp)
	for {
		w.rebootRouter()
		time.Sleep(5 * time.Second)
		changedIp := GetIp()
		log.Printf("변경 후 IP={%s}", changedIp)
		if currentIp == changedIp {
			continue
		}
		break
	}
}

func (w *Wingle) isHuawei() bool {
	get, err := w.huaweiClient.R().
		SetRetryCount(0).
		Get(statusPath)
	if err != nil {
		return false
	} else if get.StatusCode != 200 {
		return false
	}

	return true
}

func (w *Wingle) getSessionToken() SesTok {
	var st SesTok
	res, _ := w.huaweiClient.R().
		SetRetryCount(0).
		Get(sessionPath)

	res.Header.Set(`Content-Type`, defaultContentType)
	res.Unmarshal(&st)

	return st
}

func (w *Wingle) login() SesTok {
	st := w.getSessionToken()
	payload := NewLoginRequest(st.TokInfo)
	res := w.huaweiClient.R().
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

func (w *Wingle) rebootRouter() {
	if w.isHuawei() {
		w.rebootHuawei()
	} else {
		w.rebootOEM()
	}
}

func (w *Wingle) rebootOEM() {
	_, _ = w.oemClient.R().
		SetBodyString(oemRebootPayload).
		SetHeader(`Content-Type`, oemContentType).
		Post(oemRebootPath)
}

func (w *Wingle) rebootHuawei() {
	st := w.login()
	defaultHeaders := map[string]string{
		"Host":             defaultHost,
		"Origin":           w.baseUrl(),
		"Referer":          w.baseUrl() + `/html/reboot.html`,
		"X-Requested-With": `XMLHTTPRequest`,
	}
	res := w.huaweiClient.R().
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
}
