package oss

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/polaris-team/dingtalk-sdk-golang/encrypt"
	"hash"
	"io"
	"strconv"
	"strings"
	"time"
)

type PostPolicyInfo struct {
	Policy    string
	Signature string
	Dir       string
	Expire    string
	AccessId  string
	Host      string
	Region    string
	Bucket    string
}

func PostPolicyWithCallback(dir string, expire int64, maxFileSize int64, callback string) *PostPolicyInfo {
	if maxFileSize == 0{
		maxFileSize = 167772160
	}

	oc := config.GetOSSConfig()
	if oc == nil {
		panic("oss configuration is missing!")
	}

	ex := time.Now().Add(time.Duration(expire) * time.Millisecond)

	secretAccessKey := oc.AccessKeySecret
	postPolicy := GeneratePostPolicy(oc.BucketName, ex, maxFileSize, dir, callback)

	encodedPolicy := encrypt.BASE64([]byte(postPolicy))
	sign := CalculatePostSignature(encodedPolicy, secretAccessKey)

	host := strings.Join(strings.Split(oc.EndPoint, "//"), "//"+oc.BucketName+".")

	return &PostPolicyInfo{
		Policy:    encodedPolicy,
		Signature: sign,
		Dir:       dir,
		Expire:    strconv.FormatInt(ex.UnixNano()/1e6/1000, 10),
		AccessId:  oc.AccessKeyId,
		Host:      host,
		Region:    oc.EndPoint,
		Bucket:    oc.BucketName,
	}
}

func PostPolicy(dir string, expire int64, maxFileSize int64) *PostPolicyInfo {
	return PostPolicyWithCallback(dir, expire, maxFileSize, "")
}

func GeneratePostPolicy(bucket string, expire time.Time, maxFileSize int64, startsWith string, callback string) string {
	formatedExpiration := expire.UTC().Format("2006-01-02T15:04:05.999Z07:00")

	conditions := []interface{}{
		map[string]interface{}{"bucket":bucket},
		[]interface{}{"content-length-range", 0, maxFileSize},
	}

	//前缀
	if startsWith != ""{
		conditions = append(conditions, []interface{}{"starts-with", "$key", startsWith})
	}

	//后缀
	if callback != ""{
		conditions = append(conditions, map[string]interface{}{"callback":callback})
	}

	jsonizedConds := map[string]interface{}{
		"expiration": formatedExpiration,
		"conditions": conditions,
	}
	postPolicy := json.ToJsonIgnoreError(jsonizedConds)
	return postPolicy
}

func CalculatePostSignature(encodedPolicy, secretAccessKey string) string {
	return ComputeSignature(secretAccessKey, encodedPolicy)
}

func ComputeSignature(key, data string) string {
	return sign([]byte(key), data)
}

func sign(key []byte, data string) string {
	h := hmac.New(func() hash.Hash { return sha1.New() }, key)
	io.WriteString(h, data)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signedStr
}
