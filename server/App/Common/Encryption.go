package Common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"server/Base"
)

type Encryption struct{}

func (e Encryption) AesEncryptCBC(origData []byte) string {
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher([]byte(Base.AppConfig.AesKey))
	blockSize := block.BlockSize()                                                        // 获取秘钥块的长度
	origData = e.pkcs5Padding(origData, blockSize)                                        // 补全码
	blockMode := cipher.NewCBCEncrypter(block, []byte(Base.AppConfig.AesKey)[:blockSize]) // 加密模式
	encrypted := make([]byte, len(origData))                                              // 创建数组
	blockMode.CryptBlocks(encrypted, origData)                                            // 加密
	return base64.URLEncoding.EncodeToString(encrypted)
}
func (e Encryption) AesDecryptCBC(str string) (decrypted []byte, err error) {
	encrypted, err := base64.URLEncoding.DecodeString(str)
	if err != nil {
		return decrypted, err
	}

	block, err := aes.NewCipher([]byte(Base.AppConfig.AesKey))
	if err != nil {
		return decrypted, err
	} // 分组秘钥
	blockSize := block.BlockSize()                                                        // 获取秘钥块的长度
	blockMode := cipher.NewCBCDecrypter(block, []byte(Base.AppConfig.AesKey)[:blockSize]) // 加密模式
	decrypted = make([]byte, len(encrypted))                                              // 创建数组
	blockMode.CryptBlocks(decrypted, encrypted)                                           // 解密
	decrypted = e.pkcs5UnPadding(decrypted)                                               // 去除补全码
	return decrypted, nil
}

func (Encryption) pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func (Encryption) pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
