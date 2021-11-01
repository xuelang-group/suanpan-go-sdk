package storage

import (
	"io"

	"github.com/xuelang-group/suanpan-go-sdk/config"
	"github.com/xuelang-group/suanpan-go-sdk/storage"
)

func FGetObject(objectName, filePath string) error {
	s := storage.New(config.GetArgs())
	return s.FGetObject(objectName, filePath)
}

func FPutObject(objectName, filePath string) error {
	s := storage.New(config.GetArgs())
	return s.FPutObject(objectName, filePath)
}

func PutObject(objectName string, reader io.Reader) error {
	s := storage.New(config.GetArgs())
	return s.PutObject(objectName, reader)
}

func ListObjects(objectPrefix string, recursive bool, maxKeys int) ([]storage.ObjectItem, error) {
	s := storage.New(config.GetArgs())
	return s.ListObjects(objectPrefix, recursive, maxKeys)
}

func DeleteObject(objectName string) error {
	s := storage.New(config.GetArgs())
	return s.DeleteObject(objectName)
}

func DeleteMultiObjects(objectNames []string) error {
	s := storage.New(config.GetArgs())
	return s.DeleteMultiObjects(objectNames)
}