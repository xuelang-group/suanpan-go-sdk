package storage

import (
	"io"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/util"
	"github.com/xuelang-group/suanpan-go-sdk/web"
)

type OssStorage struct {
	StorageOssEndpoint    string `mapstructure:"--storage-oss-endpoint" default:"http://oss-cn-beijing.aliyuncs.com"`
	StorageOssBucketName  string `mapstructure:"--storage-oss-bucket-name" default:"suanpan"`
	StorageOssAccessId    string `mapstructure:"--storage-oss-access-id"`
	StorageOssAccessKey   string `mapstructure:"--storage-oss-access-key"`
	StorageOssTempStore   string `mapstructure:"--storage-oss-temp-store"`
	StorageOssGlobalStore string `mapstructure:"--storage-oss-global-store"`
}

func newOssStorage(argsMap map[string]string) *OssStorage {
	return &OssStorage{
		StorageOssEndpoint: util.MapDefault(argsMap, "--storage-oss-endpoint", "http://oss-cn-beijing.aliyuncs.com"),
		StorageOssBucketName: util.MapDefault(argsMap, "--storage-oss-bucket-name", "suanpan"),
		StorageOssAccessId: argsMap["--storage-oss-access-id"],
		StorageOssAccessKey: argsMap["--storage-oss-access-key"],
		StorageOssTempStore: argsMap["--storage-oss-temp-store"],
		StorageOssGlobalStore: argsMap["--storage-oss-global-store"],
	}
}

func (o *OssStorage) getBucket() (*oss.Bucket, error) {
	resp, err := web.GetStsTokenResp()
	if err != nil {
		return nil, err
	}

	cli, err := oss.New(o.StorageOssEndpoint, resp.Credentials.AccessKeyId,
		resp.Credentials.AccessKeySecret, oss.SecurityToken(resp.Credentials.SecurityToken))
	if err != nil {
		return nil, err
	}

	bucket, err := cli.Bucket(o.StorageOssBucketName)
	if err != nil {
		log.Errorf("Get oss bucket error: %v", err)
		return nil, err
	}

	return bucket, nil
}

func (o *OssStorage) FGetObject(objectName, filePath string) error {
	bucket, err := o.getBucket()
	if err != nil {
		return err
	}

	return bucket.GetObjectToFile(objectName, filePath)
}

func (o *OssStorage) FPutObject(objectName, filePath string) error {
	bucket, err := o.getBucket()
	if err != nil {
		return err
	}

	return bucket.PutObjectFromFile(objectName, filePath)
}

func (o *OssStorage) PutObject(objectName string, reader io.Reader) error {
	bucket, err := o.getBucket()
	if err != nil {
		return err
	}

	return bucket.PutObject(objectName, reader)
}

func (o *OssStorage) ListObjects(objectPrefix string, recursive bool, maxKeys int) ([]ObjectItem, error) {
	bucket, err := o.getBucket()
	if err != nil {
		return nil, err
	}

	delimiter := oss.Delimiter(`/`)
	if recursive {
		delimiter = oss.Delimiter(``)
	}

	res, err := bucket.ListObjectsV2(oss.Prefix(objectPrefix), oss.MaxKeys(maxKeys), delimiter)
	if err != nil {
		log.Errorf("List oss objects error: %v", err)
		return nil, err
	}

	objects := make([]ObjectItem, 0)
	for _, o := range res.Objects {
		objects = append(objects, ObjectItem{
			Name:         o.Key,
			LastModified: util.ISOString(o.LastModified),
		})
	}

	return objects, nil
}

func (o *OssStorage) DeleteObject(objectName string) error {
	bucket, err := o.getBucket()
	if err != nil {
		return err
	}

	return bucket.DeleteObject(objectName)
}

func (o *OssStorage) DeleteMultiObjects(objectNames []string) error {
	bucket, err := o.getBucket()
	if err != nil {
		return err
	}

	_, err = bucket.DeleteObjects(objectNames)
	return err
}