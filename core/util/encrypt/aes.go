package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/strs"
)

var key = "&polaris*aes#key"

func AesEncrypt(orig string) (string, error) {
	// 转成字节数组
	origData := []byte(orig)
	k := []byte(key)

	// 分组秘钥
	block, err := aes.NewCipher(k)
	if err != nil {
		logger.GetDefaultLogger().Error(strs.ObjectToString(err))
		return "", err
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = PKCS7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)

	return Base58Encode(cryted), nil
}

func AesDecrypt(cryted string) (string, error) {
	// 转成字节数组
	crytedByte, err := Base58Decode(cryted)
	if err != nil {
		logger.GetDefaultLogger().Error(strs.ObjectToString(err))
		return "", err
	}
	k := []byte(key)

	// 分组秘钥
	block, err := aes.NewCipher(k)
	if err != nil {
		logger.GetDefaultLogger().Error(strs.ObjectToString(err))
		return "", err
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	if len(crytedByte) == 0 {
		return "", nil
	}
	// 创建数组
	orig := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(orig, crytedByte)
	// 去补全码
	orig = PKCS7UnPadding(orig)

	return string(orig), nil
}
