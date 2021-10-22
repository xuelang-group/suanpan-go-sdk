package storage

import (
	"io"

	"github.com/golang/glog"
	"github.com/mcuadros/go-defaults"
	"github.com/mitchellh/mapstructure"
	"github.com/xuelang-group/suanpan-go-sdk/config"
)

type Storage interface {
	FGetObject(objectName, filePath string) error
	FPutObject(objectName, filePath string) error
	PutObject(objectName string, reader io.Reader) error
	ListObjects(objectPrefix string, recursive bool, maxKeys int) ([]ObjectItem, error)
	DeleteObject(objectName string) error
	DeleteMultiObjects(objectNames []string) error
}

type EnvStorage struct {
	StorageType	string	`mapstructure:"--storage-type" default:"minio"`
}

const (
	Minio = "minio"
	Oss = "oss"
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
	case Oss:
		var ossStorage OssStorage
		mapstructure.Decode(argsMap, &ossStorage)
		defaults.SetDefaults(&ossStorage)
		return &ossStorage
	default:
		glog.Errorf("Unsupported storage type: %s", envStorage.StorageType)
		return nil
	}
}