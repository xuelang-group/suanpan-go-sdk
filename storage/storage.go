package storage

import (
	"github.com/golang/glog"
	"github.com/mcuadros/go-defaults"
	"github.com/mitchellh/mapstructure"
	"github.com/xuelang-group/suanpan-go-sdk/config"
)

type Storage interface {

}

type EnvStorage struct {
	StorageType	string	`mapstructure:"--storage-type" default:"minio"`
}

const (
	Minio = "minio"
)

func GetStorage() Storage {
	argsMap := config.GetArgs()
	var envStorage EnvStorage
	mapstructure.Decode(argsMap, &envStorage)
	defaults.SetDefaults(&envStorage)
	switch envStorage.StorageType {
	case Minio:
		var minioStorage MinioStorage
		mapstructure.Decode(argsMap, &minioStorage)
		defaults.SetDefaults(&minioStorage)
		return &minioStorage
	default:
		glog.Errorf("Unsupported storage type: %s", envStorage.StorageType)
		return nil
	}
}