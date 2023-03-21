package wingle

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"github.com/imroc/req/v3"
	"time"
)

const (
	ipApiHost = "http://ip-api.com"
	ipApiPath = "json"
)

var healthCheckClient = req.C().
	SetBaseURL(ipApiHost).
	SetCommonRetryCount(-1).
	SetCommonRetryCondition(func(resp *req.Response, err error) bool {
		return err != nil || resp.StatusCode >= 400
	}).
	SetCommonRetryFixedInterval(3 * time.Second)

type IpApiRes struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
}

func GetIp() string {
	var res IpApiRes
	_, err := healthCheckClient.R().
		SetSuccessResult(&res).
		Get(ipApiPath)

	if err != nil {
		return ""
	}

	return res.Query
}

func Base64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func Sha256(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))

	return hex.EncodeToString(hash.Sum(nil))
}
