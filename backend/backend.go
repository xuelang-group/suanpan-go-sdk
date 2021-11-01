package backend

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/golang/glog"
	"github.com/xuelang-group/suanpan-go-sdk/config"
)

const (
	Windows = "windows"
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

func GetStsTokenResp() (*StsTokenResp, error) {
	b, err := strconv.ParseBool(config.GetEnv().SpHostTls)
	if err != nil {
		glog.Warningf("SpHostTls is not a valid bool value: %s", config.GetEnv().SpHostTls)
		b = false
	}
	protocol := `http`
	if b {
		protocol = `https`
	}
	path := `/oss/token`
	host := config.GetEnv().SpHost
	if config.GetEnv().SpOs == Windows {
		host = host + `:` + config.GetEnv().SpPort
	}

	req, err := http.NewRequest("GET", protocol+`://`+host+path, nil)
	for k, v := range GetHeaders() {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		glog.Errorf("Request sts token error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf("Read response body error: %v", err)
		return nil, err
	}

	var stsTokenResp *StsTokenResp
	err = json.Unmarshal(data, &stsTokenResp)
	if err != nil {
		glog.Errorf("Unmarshal json format error: %v", err)
	}

	return stsTokenResp, nil
}