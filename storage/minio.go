package storage

import (
	"io"
	"net/url"

	"github.com/mcuadros/go-defaults"
	"github.com/minio/minio-go"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/util"
)

type MinioStorage struct {
	StorageMinioEndpoint   string `mapstructure:"--storage-minio-endpoint" default:"http://minio-service.default:9000"`
	StorageMinioBucketName string `mapstructure:"--storage-minio-bucket-name" default:"suanpan"`
	StorageMinioAccessKey  string `mapstructure:"--storage-minio-access-key"`
	StorageMinioSecretKey  string `mapstructure:"--storage-minio-secret-key"`
	StorageMinioTempStore  string `mapstructure:"--storage-minio-temp-store"`
}

func (m *MinioStorage) getClient() (*minio.Client, error) {
	var minioStorage MinioStorage
	defaults.SetDefaults(&minioStorage)
	u, err := url.Parse(m.StorageMinioEndpoint)
	if err != nil {
		log.Errorf("Parse StorageMinioEndpoint error: %w", err)
		return nil, err
	}

	cli, err := minio.New(
		u.Host, m.StorageMinioAccessKey,
		m.StorageMinioSecretKey, u.Scheme == "https")
	if err != nil {
		log.Errorf("Init minio client error: %w", err)
		return nil, err
	}

	return cli, nil
}

func (m *MinioStorage) FGetObject(objectName, filePath string) error {
	cli, err := m.getClient()
	if err != nil {
		return err
	}

	return cli.FGetObject(m.StorageMinioBucketName, objectName, filePath, minio.GetObjectOptions{})
}

func (m *MinioStorage) FPutObject(objectName, filePath string) error {
	cli, err := m.getClient()
	if err != nil {
		return err
	}

	n, err := cli.FPutObject(m.StorageMinioBucketName, objectName, filePath, minio.PutObjectOptions{})
	log.Infof("Uploaded %d bytes", n)
	return err
}

func (m *MinioStorage) PutObject(objectName string, reader io.Reader) error {
	cli, err := m.getClient()
	if err != nil {
		return err
	}

	n, err := cli.PutObject(m.StorageMinioBucketName, objectName, reader, -1, minio.PutObjectOptions{})
	log.Infof("Uploaded %d bytes", n)
	return err
}

func (m *MinioStorage) ListObjects(objectPrefix string, recursive bool, maxKeys int) ([]ObjectItem, error) {
	cli, err := m.getClient()
	if err != nil {
		return nil, err
	}

	doneCh := make(chan struct{})
	defer close(doneCh)

	objects := make([]ObjectItem, 0)
	for o := range cli.ListObjectsV2(m.StorageMinioBucketName, objectPrefix, recursive, doneCh) {
		if o.Err != nil {
			log.Errorf("List objects error: %w", o.Err)
			return nil, o.Err
		}
		objects = append(objects, ObjectItem{
			Name:         o.Key,
			LastModified: util.ISOString(o.LastModified),
		})
		if len(objects) >= maxKeys {
			doneCh <- struct{}{}
			break
		}
	}

	return objects, nil
}

func (m *MinioStorage) DeleteObject(objectName string) error {
	cli, err := m.getClient()
	if err != nil {
		return err
	}

	return cli.RemoveObject(m.StorageMinioBucketName, objectName)
}

func (m *MinioStorage) DeleteMultiObjects(objectNames []string) error {
	cli, err := m.getClient()
	if err != nil {
		return err
	}

	objectsCh := make(chan string, 0)
	defer close(objectsCh)
	for _, o := range objectNames {
		objectsCh <- o
	}

	go func() {
		for err := range cli.RemoveObjects(m.StorageMinioBucketName, objectsCh) {
			log.Errorf("Remove object error: %w", err)
		}
	}()

	return nil
}
