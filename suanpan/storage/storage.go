package storage

import (
	"io"

	"github.com/xuelang-group/suanpan-go-sdk/storage"
)

func FGetObject(objectName, filePath string) error {
	s := storage.GetStorage()
	return s.FGetObject(objectName, filePath)
}

func FPutObject(objectName, filePath string) error {
	s := storage.GetStorage()
	return s.FPutObject(objectName, filePath)
}

func PutObject(objectName string, reader io.Reader) error {
	s := storage.GetStorage()
	return s.PutObject(objectName, reader)
}

func ListObjects(objectPrefix string, recursive bool, maxKeys int) ([]storage.ObjectItem, error) {
	s := storage.GetStorage()
	return s.ListObjects(objectPrefix, recursive, maxKeys)
}

func DeleteObject(objectName string) error {
	s := storage.GetStorage()
	return s.DeleteObject(objectName)
}

func DeleteMultiObjects(objectNames []string) error {
	s := storage.GetStorage()
	return s.DeleteMultiObjects(objectNames)
}