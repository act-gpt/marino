package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const SHA265 = "0123456789abcdefghijklmnopqrstuvwxyz"

func getId(str string) int {
	id := 0
	for i := 0; i < len(str); i++ {
		id += int(str[i])
	}
	return id
}

func Timestamp(str string, chars string) int64 {
	// 所有字符
	c := strings.Split(chars, "")
	id := getId(str)
	// 取特定字符
	s := c[id%len(c)]
	// 特定字符里的次数
	count := strings.Count(strings.ToLower(str), s)
	// 当前时间
	date := int64(time.Now().UnixMilli())
	// 计算时间戳
	t := int64(count)
	if t > 0 {
		t = date % int64(count)
	}
	return date - t + int64(count)
}

func CheckTimeStamp(timestamp float64, str, chars string) bool {
	// 所有字符
	c := strings.Split(chars, "")
	// 取特定字符
	id := getId(str)
	s := c[int(id)%len(c)]
	// 特定字符里的次数
	count := strings.Count(strings.ToLower(str), s)
	t := count
	if t > 0 {
		t = int(timestamp) % count
	}
	return t == 0
}

func isObject(v interface{}) bool {
	_, ok := v.(map[string]interface{})
	return ok
}

func Marshall(params map[string]interface{}) string {
	if params == nil {
		params = make(map[string]interface{})
	}
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var kvs []string
	for _, k := range keys {
		v := params[k]
		if v == nil {
			delete(params, k)
			continue
		}
		kvs = append(kvs, k+"="+params[k].(string))
	}
	return strings.Join(kvs, "&")
}

func Sign(str string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Equal(item map[string]interface{}, secret string, signStr string) bool {
	sign := signStr
	delete(item, "sign")
	params := Marshall(item)
	return sign == Sign(params, secret)
}

// pkcs7Padding 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	//判断缺少几位长度。最少1，最多 blockSize
	padding := blockSize - len(data)%blockSize
	//补足位数。把切片[]byte{byte(padding)}复制padding个
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// pkcs7UnPadding 填充的反向操作
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("加密字符串错误！")
	}
	//获取填充的个数
	unPadding := int(data[length-1])
	return data[:(length - unPadding)], nil
}

// AesEncrypt 加密
func AesEncrypt(data []byte, key []byte) ([]byte, error) {
	//创建加密实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//判断加密快的大小
	blockSize := block.BlockSize()
	//填充
	encryptBytes := pkcs7Padding(data, blockSize)
	//初始化加密数据接收切片
	crypted := make([]byte, len(encryptBytes))
	//使用cbc加密模式
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	//执行加密
	blockMode.CryptBlocks(crypted, encryptBytes)
	return crypted, nil
}

// AesDecrypt 解密
func AesDecrypt(data []byte, key []byte) ([]byte, error) {
	//创建实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//获取块的大小
	blockSize := block.BlockSize()
	//使用cbc
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	//初始化解密数据接收切片
	crypted := make([]byte, len(data))
	//执行解密
	blockMode.CryptBlocks(crypted, data)
	//去除填充
	crypted, err = pkcs7UnPadding(crypted)
	if err != nil {
		return nil, err
	}
	return crypted, nil
}

// EncryptByAes Aes加密 后 base64 再加
func EncryptByAes(PwdKey []byte, data []byte) (string, error) {
	res, err := AesEncrypt(data, PwdKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(res), nil
}

// DecryptByAes Aes 解密
func DecryptByAes(PwdKey []byte, data string) ([]byte, error) {
	dataByte, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return AesDecrypt(dataByte, PwdKey)
}

// 密码 hash
func Password2Hash(password string) (string, error) {
	passwordBytes := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	return string(hashedPassword), err
}

// 验证密码
func ValidatePasswordAndHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// 加密请求字符串
func Encode(params map[string]interface{}, secret string) string {
	str := Marshall(params)
	return Sign(str, secret)
}

// 检查请求字符串
func Check(data map[string]interface{}, code string) bool {
	_, hasSign := data["sign"]
	_, hasNonce := data["nonce"]
	if !hasSign || !hasNonce || len(data["nonce"].(string)) < 8 {
		return false
	}
	return Equal(data, code, data["sign"].(string))
}

// Nonce
func Nonce(length int) string {
	if length <= 0 {
		length = 8
	}
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	pos := len(chars)
	nonces := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		nonces[i] = chars[r.Intn(pos)]
	}
	return string(nonces)
}
