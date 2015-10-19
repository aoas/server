package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
)

func EncodeMd5(value string) string {
	h := md5.New()
	h.Write([]byte(value))
	return hex.EncodeToString(h.Sum(nil))
}

func EncodeSha1(value string) string {
	h := sha1.New()
	h.Write([]byte(value))
	return hex.EncodeToString(h.Sum(nil))
}

func GetRandomString(n int, alphabets ...byte) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		if len(alphabets) == 0 {
			bytes[i] = alphanum[b%byte(len(alphanum))]
		} else {
			bytes[i] = alphabets[b%byte(len(alphabets))]
		}
	}
	return string(bytes)
}
func Base64EncodeByte(data []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(data))
}

func Base64DecodeByte(data []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(data))
}

func Base64Encode(data string) string {
	return string(Base64EncodeByte([]byte(data)))
}

func Base64Decode(data string) (string, error) {
	d, e := Base64DecodeByte([]byte(data))
	return string(d), e
}

func AESEncodeByte(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return data, err
	}
	var iv = key[:aes.BlockSize]
	blockMode := cipher.NewCFBEncrypter(block, iv)
	dest := make([]byte, len(string(data)))
	blockMode.XORKeyStream(dest, data)
	return dest, nil
}
func AESDecodeByte(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return data, err
	}
	var iv = key[:aes.BlockSize]
	blockMode := cipher.NewCFBDecrypter(block, iv)
	dest := make([]byte, len(string(data)))
	blockMode.XORKeyStream(dest, data)
	return dest, nil
}
func AESEncode(data string, key string) (string, error) {
	out, err := AESEncodeByte([]byte(data), []byte(key))
	return string(Base64EncodeByte(out)), err
}
func AESDecode(data string, key string) (string, error) {
	d, e := Base64DecodeByte([]byte(data))
	if e != nil {
		return data, e
	}
	out, err := AESDecodeByte(d, []byte(key))
	return string(out), err
}
