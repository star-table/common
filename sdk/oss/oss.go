package oss

import (
	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/consts"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
)

func signUrl(fileName string, expiredInSec int64, options ...oss.Option) (string, error) {
	oc := config.GetOSSConfig()
	if oc == nil {
		panic(consts.OssConfigMissingSentence)
	}

	client, err := oss.New(oc.EndPoint, oc.AccessKeyId, oc.AccessKeySecret)
	if err != nil {
		return "", err
	}

	bucket, err := client.Bucket(oc.BucketName)
	if err != nil {
		return "", err
	}

	signedURL, err := bucket.SignURL(fileName, oss.HTTPPut, expiredInSec, options...)
	if err != nil {
		return "", err
	}
	return signedURL, nil
}

func SignUrlWithStream(fileName string, fileSuffix string, expiredInSec int64) (string, error) {
	var options = []oss.Option{
		oss.ACL(oss.ACLPublicRead),
		oss.ContentType("image/" + fileSuffix),
	}

	return signUrl(fileName+"."+fileSuffix, expiredInSec, options...)
}

func PutObjectWithURL(signUrl string, fileSuffix string, reader io.Reader) error {
	oc := config.GetOSSConfig()
	if oc == nil {
		panic(consts.OssConfigMissingSentence)
	}

	client, err := oss.New(oc.EndPoint, oc.AccessKeyId, oc.AccessKeySecret)
	if err != nil {
		return err
	}

	bucket, err := client.Bucket(oc.BucketName)
	if err != nil {
		return err
	}
	var options = []oss.Option{
		oss.ACL(oss.ACLPublicRead),
		oss.ContentType("image/" + fileSuffix),
	}

	err = bucket.PutObjectWithURL(signUrl, reader, options...)
	if err != nil {
		return err
	}
	return nil
}

func GetObjectUrl(key string, expiredInSec int64) (string, error) {
	oc := config.GetOSSConfig()
	if oc == nil {
		panic("oss configuration is missing!")
	}

	client, err := oss.New(oc.EndPoint, oc.AccessKeyId, oc.AccessKeySecret)
	if err != nil {
		return "", err
	}

	bucket, err := client.Bucket(oc.BucketName)
	if err != nil {
		return "", err
	}

	return bucket.SignURL(key, oss.HTTPGet, expiredInSec)
}
