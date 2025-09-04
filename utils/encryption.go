package utils

import (
	"EZ-Encrypt-Middleware/config"

	"github.com/deatil/go-cryptobin/cryptobin/crypto"
)

// Decryption 传入待解密的字符串，返回解密后的字符串
func Decryption(str string, iv string) string {
	cyptde := crypto.
		FromBase64String(str).
		SetKey(config.AppConfig.AESKey).
		SetIv(iv).
		Aes().
		CBC().
		PKCS7Padding().
		Decrypt().
		ToString()

	return cyptde
}
