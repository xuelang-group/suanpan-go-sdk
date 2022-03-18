package storage

import (
	"io"

	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/util"
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

func New(argsMap map[string]string) Storage {
	envStorage := newEnvStorage(argsMap)
	switch envStorage.StorageType {
	case Minio:
		return newMinioStorage(argsMap)
	case Oss:
		return newOssStorage(argsMap)
	default:
		log.Errorf("Unsupported storage type: %s", envStorage.StorageType)
		return nil
	}
}

func newEnvStorage(argsMap map[string]string) *EnvStorage {
	return &EnvStorage{
		StorageType: util.MapDefault(argsMap, "--storage-type", "minio"),
	}
}