package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/thoas/go-funk"
	"github.com/xuelang-group/suanpan-go-sdk/config"
	"github.com/xuelang-group/suanpan-go-sdk/web"
)

func Lookup(key string) (string, error) {
	graph, err := web.GetGraph()
	if err != nil {
		return "", err
	}
	nodeId := config.GetEnv().SpNodeId
	process := graph.Processes[nodeId]
	port, err := findPortFromProcess(process, key)
	if err != nil {
		return "", err
	}
	connections := funk.Filter(graph.Connections, func(c web.Connection) bool {
		return c.Src.Process == nodeId && c.Src.Port == port
	}).([]web.Connection)
	if len(connections) > 0 {
		return getServiceName(connections[0].Tgt.Process), nil
	}
	return "", errors.New(fmt.Sprintf("node %s not found in graph", nodeId))
}

func findPortFromProcess(process web.Process, key string) (string, error) {
	for _, port := range process.Metadata.Def.Ports {
		if port.UUID == key || port.Description.ZH_CN == key || port.Description.EN_US == key {
			return port.UUID, nil
		}
	}
	return "", errors.New(fmt.Sprintf("port %s not found", key))
}

func getServiceName(nodeId string) string {
	return strings.Join([]string{"app", config.GetEnv().SpAppId, nodeId}, "-")
}