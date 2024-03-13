package main

import (
	"crypto/aes"
	"crypto/cipher"
	"strings"
)

func getVt(iv string) []byte {
	iv = strings.Replace(iv, "0x", "", 1)[:16]

	return []byte(iv)
}

func decrypt(contentKey, vt []byte, encryptedBytes []byte) ([]byte, error) {
	// 使用AES算法创建解密器
	block, err := aes.NewCipher(contentKey)
	if err != nil {
		return nil, err
	}

	// CBC模式解密器
	decrypter := cipher.NewCBCDecrypter(block, vt)

	// 解密数据
	decryptedData := make([]byte, len(encryptedBytes))
	decrypter.CryptBlocks(decryptedData, encryptedBytes)

	// 去除填充
	decryptedData = pkcs7UnPad(decryptedData)

	return decryptedData, nil
}

// pkcs7UnPad 去除PKCS7填充
func pkcs7UnPad(data []byte) []byte {
	length := len(data)
	unPadding := int(data[length-1])
	return data[:(length - unPadding)]
}
