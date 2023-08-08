package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	rand2 "crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"io"
	"math/rand"
	"time"
)

func GenerateRandomUserID() int64 {
	node, err := snowflake.NewNode(1) // 传入机器 ID
	if err != nil {
		// 错误处理
	}

	return node.Generate().Int64()
}
func GetRandomString(l int) string {
	str := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
func GenerateRandomNumber() int64 {
	// 设置随机种子，保证每次生成的随机数都不同
	rand.Seed(time.Now().UnixNano())

	// 生成10位随机整数，范围在0到9999999999之间
	return int64(rand.Intn(10000000000))
}

func FixMd5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func Encrypt(key []byte, plaintext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand2.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(key []byte, ciphertext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	decodedCipherText, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	if len(decodedCipherText) < aes.BlockSize {
		return "", fmt.Errorf("加密数据长度错误")
	}

	iv := decodedCipherText[:aes.BlockSize]
	decodedCipherText = decodedCipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(decodedCipherText, decodedCipherText)

	return string(decodedCipherText), nil
}
