package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

// ECB模式加密
func AesEncryptECB(src, key []byte) ([]byte, error) {
	plaintext := PKCS5Padding(src, aes.BlockSize)
	cipher, err := aes.NewCipher(key[:aes.BlockSize])
	if err != nil {
		return nil, err
	}

	if len(plaintext)%aes.BlockSize != 0 {
		panic("EncryptECB ERROR: Need a multiple of the blocksize 16")
	}

	ciphertext := make([]byte, 0)
	text := make([]byte, 16)
	for len(plaintext) > 0 {
		// 每次运算一个block
		cipher.Encrypt(text, plaintext)
		plaintext = plaintext[aes.BlockSize:]
		ciphertext = append(ciphertext, text...)
	}
	return ciphertext, nil
}

// ECB模式解密
func AesDecryptECB(ciphertext, key []byte) ([]byte, error) {
	cipher, err := aes.NewCipher(key[:aes.BlockSize])
	if err != nil {
		return nil, err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		panic("Need a multiple of the blocksize 16")
	}

	plaintext := make([]byte, 0)
	text := make([]byte, 16)
	for len(ciphertext) > 0 {
		cipher.Decrypt(text, ciphertext)
		ciphertext = ciphertext[aes.BlockSize:]
		plaintext = append(plaintext, text...)
	}
	return plaintext, nil
}

//AES CBC encrypt
func AesEncryptCBC(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	//origData = ZeroPadding(origData, blockSize)
	//blockMode := cipher.e
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

//AES CBC decrypt
func AesDecryptCBC(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

// Padding补全
//pkcs7 padding
func PKCS7Pad(data []byte) []byte {
	padding := aes.BlockSize - len(data)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(0)}, padding)
	return append(data, padtext...)
}

func PKCS7UPad(data []byte) []byte {
	padLength := int(data[len(data)-1])
	return data[:len(data)-padLength]
}

//pksc5 padding
func PKCS5Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	//println("000000000000000000000000000000000000000000000", padding)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//zero padding
func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
