package cryptor

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

func AESEncryptCBC(origData []byte, key []byte) (retEncrypted []byte, retErr error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		retErr = err
		return
	}
	blockSize := block.BlockSize()
	origData = pkcs5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	retEncrypted = make([]byte, len(origData))
	blockMode.CryptBlocks(retEncrypted, origData)
	return
}

func AESDecryptCBC(encrypted []byte, key []byte) (retEncrypted []byte, retErr error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		retErr = err
		return
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	retEncrypted = make([]byte, len(encrypted))
	blockMode.CryptBlocks(retEncrypted, encrypted)
	retEncrypted = pkcs5UnPadding(retEncrypted)
	return
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
