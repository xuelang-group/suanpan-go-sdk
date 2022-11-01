package microservice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/xuelang-group/suanpan-go-sdk/config"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
)

func Lookup(portKey string) (string, error) {
	appId := config.GetEnv().SpAppId
	nodeId := config.GetEnv().SpNodeId

	url := fmt.Sprintf("http://app-%s:8001/internal/microservice/lookup", appId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errorf("New lookup Request error: %v", err)
		return "", err
	}
	q := req.URL.Query()
	q.Add("nodeId", nodeId)
	q.Add("portKey", portKey)
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("Do lookup request error: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Read lookup response body error: %v", err)
		return "", err
	}

	var lookupResp *LookupResp
	err = json.Unmarshal(body, &lookupResp)
	if err != nil {
		log.Errorf("Unmarshal json format error: %v", err)
		return "", err
	}

	return lookupResp.Data.Address, nil
}
