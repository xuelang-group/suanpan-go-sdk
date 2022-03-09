package storage

import (
	"path"

	"github.com/xuelang-group/suanpan-go-sdk/config"
)

const (
	pathPrefix     = `studio`
	dataPath       = `share`
	configsPath    = `configs`
	tmpPath        = `tmp`
	logsPath       = `logs`
	componentsPath = `component`
)

func GetKeyInAppStore(paths ...string) string {
	return path.Join(appStoreKey(), path.Join(paths...))
}

func GetKeyInAppDataStore(paths ...string) string {
	return path.Join(appDataStoreKey(), path.Join(paths...))
}

func GetKeyInAppConfigsStore(paths ...string) string {
	return path.Join(appConfigsStoreKey(), path.Join(paths...))
}

func GetKeyInAppTmpStore(paths ...string) string {
	return path.Join(appTmpStoreKey(), path.Join(paths...))
}

func GetKeyInAppLogsStore(paths ...string) string {
	return path.Join(appLogsStoreKey(), path.Join(paths...))
}

func GetKeyInNodeStore(paths ...string) string {
	return path.Join(nodeStoreKey(), path.Join(paths...))
}

func GetKeyInNodeDataStore(paths ...string) string {
	return path.Join(nodeDataStoreKey(), path.Join(paths...))
}

func GetKeyInNodeConfigsStore(paths ...string) string {
	return path.Join(nodeConfigsStoreKey(), path.Join(paths...))
}

func GetKeyInNodeTmpStore(paths ...string) string {
	return path.Join(nodeTmpStoreKey(), path.Join(paths...))
}

func GetKeyInNodeLogsStore(paths ...string) string {
	return path.Join(nodeLogsStoreKey(), path.Join(paths...))
}

func GetKeyInComponentsStoreKey(paths ...string) string {
	return path.Join(componentsStoreKey(), path.Join(paths...))
}

func appStoreKey() string {
	return path.Join(pathPrefix, config.GetEnv().SpUserId, config.GetEnv().SpAppId)
}

func appDataStoreKey() string {
	return path.Join(pathPrefix, config.GetEnv().SpUserId, dataPath, config.GetEnv().SpAppId)
}

func appConfigsStoreKey() string {
	return path.Join(pathPrefix, config.GetEnv().SpUserId, configsPath, config.GetEnv().SpAppId)
}

func appTmpStoreKey() string {
	return path.Join(pathPrefix, config.GetEnv().SpUserId, tmpPath, config.GetEnv().SpAppId)
}

func appLogsStoreKey() string {
	return path.Join(pathPrefix, config.GetEnv().SpUserId, logsPath, config.GetEnv().SpAppId)
}

func nodeStoreKey() string {
	return path.Join(pathPrefix, config.GetEnv().SpUserId, config.GetEnv().SpAppId, config.GetEnv().SpNodeId)
}

func nodeDataStoreKey() string {
	return path.Join(pathPrefix, config.GetEnv().SpUserId, dataPath, config.GetEnv().SpAppId, config.GetEnv().SpNodeId)
}

func nodeConfigsStoreKey() string {
	return path.Join(pathPrefix, config.GetEnv().SpUserId, configsPath, config.GetEnv().SpAppId, config.GetEnv().SpNodeId)
}

func nodeTmpStoreKey() string {
	return path.Join(pathPrefix, config.GetEnv().SpUserId, tmpPath, config.GetEnv().SpAppId, config.GetEnv().SpNodeId)
}

func nodeLogsStoreKey() string {
	return path.Join(pathPrefix, config.GetEnv().SpUserId, logsPath, config.GetEnv().SpAppId, config.GetEnv().SpNodeId)
}

func componentsStoreKey() string {
	return path.Join(pathPrefix, config.GetEnv().SpUserId, componentsPath)
}
