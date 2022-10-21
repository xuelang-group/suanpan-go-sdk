package web

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/xuelang-group/suanpan-go-sdk/config"
	"github.com/xuelang-group/suanpan-go-sdk/util"
)

const (
	Windows = "windows"
)

var (
	httpServerUrl string
	httpOnce      sync.Once
)

func getHttpServerUrl() string {
	httpOnce.Do(func() {
		httpServerUrl = buildHttpServerUrl()
	})

	return httpServerUrl
}

func buildHttpServerUrl() string {
	b, err := strconv.ParseBool(config.GetEnv().SpHostTls)
	if err != nil {
		logrus.Warnf("SpHostTls is not a valid bool value: %s", config.GetEnv().SpHostTls)
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
		logrus.Errorf("Request sts token error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Read response body error: %v", err)
		return nil, err
	}

	var stsTokenResp *StsTokenResp
	err = json.Unmarshal(data, &stsTokenResp)
	if err != nil {
		logrus.Errorf("Unmarshal json format error: %v", err)
	}

	return stsTokenResp, nil
}

func GetGraph() (*Graph, error) {
	graphPath := path.Join(`/appcontroller/graph`, config.GetEnv().SpUserId, config.GetEnv().SpAppId)
	req, err := http.NewRequest("GET", getHttpServerUrl()+graphPath, nil)
	req.Header = GetHeaders()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Errorf("Request graph error: %v", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Read response body error: %v", err)
		return nil, err
	}

	var graphResp *GraphResp
	err = json.Unmarshal(data, &graphResp)
	if err != nil {
		logrus.Errorf("Unmarshal graphResp json format error: %v", err)
		return nil, err
	}

	var graph *Graph
	err = json.Unmarshal([]byte(graphResp.GraphJson.GraphJsonStr), &graph)
	if err != nil {
		logrus.Errorf("Unmarshal graph json format error: %v", err)
		return nil, err
	}

	return graph, nil
}

func RegisterFreePort(nodePort string) error {
	port, err := util.GetFreePort()
	if err != nil {
		logrus.Errorf("Get free port error: %v", err)
		return err
	}

	return RegisterPort(nodePort, strconv.Itoa(port))
}

func RegisterPort(nodePort string, port string) error {
	path := "/app/service/register"
	body := map[string]string{
		"appId": config.GetEnv().SpAppId,
		"nodeId": config.GetEnv().SpNodeId,
		"userId": config.GetEnv().SpUserId,
		"nodePort": nodePort,
		"port": port,
	}
	bodyByte, err := json.Marshal(body)
	if err != nil {
		logrus.Errorf("Unmarshal json format error: %v", err)
	}
	req, err := http.NewRequest("POST", getHttpServerUrl()+path, bytes.NewBuffer(bodyByte))
	req.Header = GetHeaders()
	req.Header.Set("Content-Type","application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Errorf("Regiter port error: %v", err)
		return err
	}
	defer resp.Body.Close()

	return nil
}