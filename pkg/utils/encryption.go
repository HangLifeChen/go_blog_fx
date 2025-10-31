package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/mr-tron/base58"
)

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Sha256(encode_str string) [32]byte {
	return sha256.Sum256([]byte(encode_str))
}

func Base58(hash [32]byte) string {
	return base58.Encode(hash[:])
}

func AesEncryptGCM(plainText, key string) (string, error) {
	hashKey := sha256.Sum256([]byte(key))
	// 创建AES密码块
	block, err := aes.NewCipher(hashKey[:])
	if err != nil {
		return "", err
	}

	// 创建GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 生成随机Nonce（推荐12字节）
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 加密并附加Nonce到结果中
	decoded := base64.StdEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(plainText), nil))

	return decoded, nil
}

func AesDecryptGCM(ciphertext, key string) (string, error) {
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	ciphertext = string(ciphertextBytes)
	hashKey := sha256.Sum256([]byte(key))
	// 创建AES密码块
	block, err := aes.NewCipher(hashKey[:])
	if err != nil {
		return "", err
	}

	// 创建GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 获取Nonce大小
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	// 分离Nonce和实际密文
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// 解密并验证
	plainText, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
