package web

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"

	"github.com/golang/glog"
	"github.com/xuelang-group/suanpan-go-sdk/config"
)

const (
	Windows = "windows"
)

var (
	httpServerUrl string
	httpOnce      sync.Once
)

type Credentials struct {
	SecurityToken   string `json:"SecurityToken"`
	AccessKeyId     string `json:"AccessKeyId"`
	AccessKeySecret string `json:"AccessKeySecret"`
	Expiration      string `json:"Expiration"`
}

type StsTokenResp struct {
	RequestId   string      `json:"RequestId"`
	Credentials Credentials `json:"Credentials"`
}

func getHttpServerUrl() string {
	httpOnce.Do(func() {
		httpServerUrl = buildHttpServerUrl()
	})

	return httpServerUrl
}

func buildHttpServerUrl() string {
	b, err := strconv.ParseBool(config.GetEnv().SpHostTls)
	if err != nil {
		glog.Warningf("SpHostTls is not a valid bool value: %s", config.GetEnv().SpHostTls)
		b = false
	}
	protocol := `http`
	if b {
		protocol = `https`
	}
	host := config.GetEnv().SpHost
	if config.GetEnv().SpOs == Windows {
		host = host + `:` + config.GetEnv().SpPort
	}

	return protocol + `://` + host
}

func GetStsTokenResp() (*StsTokenResp, error) {
	path := `/oss/token`
	req, err := http.NewRequest("GET", getHttpServerUrl()+path, nil)
	req.Header = GetHeaders()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		glog.Errorf("Request sts token error: %w", err)
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf("Read response body error: %w", err)
		return nil, err
	}

	var stsTokenResp *StsTokenResp
	err = json.Unmarshal(data, &stsTokenResp)
	if err != nil {
		glog.Errorf("Unmarshal json format error: %w", err)
	}

	return stsTokenResp, nil
}
