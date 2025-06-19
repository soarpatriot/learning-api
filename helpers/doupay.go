package helpers

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"strings"
)

func GetByteAuthorization(privateKeyStr, data, appId, nonceStr, timestamp, keyVersion string) (string, error) {
	var byteAuthorization string
	// 读取私钥
	key, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(privateKeyStr, "\n", ""))
	if err != nil {
		return "", err
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(key)
	if err != nil {
		return "", err
	}
	// 生成签名
	signature, err := getSignature("POST", "/requestOrder", timestamp, nonceStr, data, privateKey)
	if err != nil {
		return "", err
	}
	// 构造byteAuthorization
	byteAuthorization = fmt.Sprintf("SHA256-RSA2048 appid=%s,nonce_str=%s,timestamp=%s,key_version=%s,signature=%s", appId, nonceStr, timestamp, keyVersion, signature)
	return byteAuthorization, nil
}

func getSignature(method, url, timestamp, nonce, data string, privateKey *rsa.PrivateKey) (string, error) {
	fmt.Printf("method:%s\n url:%s\n timestamp:%s\n nonce:%s\n data:%s", method, url, timestamp, nonce, data)
	targetStr := method + "\n" + url + "\n" + timestamp + "\n" + nonce + "\n" + data + "\n"
	h := sha256.New()
	h.Write([]byte(targetStr))
	digestBytes := h.Sum(nil)

	signBytes, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, digestBytes)
	if err != nil {
		return "", err
	}
	sign := base64.StdEncoding.EncodeToString(signBytes)

	return sign, nil
}

func RandStr(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}
