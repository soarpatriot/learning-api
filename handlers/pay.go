package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"learning-api/config"
	"learning-api/helpers"

	"crypto/tls"
	"crypto/x509"

	"github.com/gin-gonic/gin"
)

// PayOrderRequest is the input struct for /pay/order
type PayOrderRequest struct {
	TotalAmount int    `json:"total_amount"`
	Subject     string `json:"subject"`
	Body        string `json:"body"`
	CpExtra     string `json:"cp_extra"`
}

// DouyinOrderRequest is the struct for Douyin API
type DouyinOrderRequest struct {
	AppID       string `json:"app_id"`
	OutOrderNo  string `json:"out_order_no"`
	TotalAmount int    `json:"total_amount"`
	Subject     string `json:"subject"`
	Body        string `json:"body"`
	ValidTime   int    `json:"valid_time"`
	Sign        string `json:"sign"`
	NotifyURL   string `json:"notify_url"`
	StoreUid    string `json:"store_uid"`
}

func randomOrderNo() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36) + RandString(4)
}

func RandString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func createSecureHTTPClient() (*http.Client, error) {
	roots, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: roots},
	}
	client := &http.Client{Transport: tr}
	return client, nil
}

func createInsecureHTTPClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}
}

func PayOrder(c *gin.Context) {
	var req PayOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cfg := config.LoadConfig()
	orderNo := randomOrderNo()
	order := DouyinOrderRequest{
		AppID:       cfg.AppID,
		OutOrderNo:  orderNo,
		TotalAmount: req.TotalAmount,
		Subject:     req.Subject,
		Body:        req.Body,
		ValidTime:   180,
		StoreUid:    "75169185453352082020",
		NotifyURL:   "https://1l2v8anoldbg6-env-KfJ4EiJx5I.service.douyincloud.run/pay/callback",
	}
	// Prepare sign params (as map)
	signParams := map[string]interface{}{
		"app_id":       order.AppID,
		"out_order_no": order.OutOrderNo,
		"total_amount": order.TotalAmount,
		"subject":      order.Subject,
		"body":         order.Body,
		"valid_time":   order.ValidTime,
		"notify_url":   order.NotifyURL,
	}
	order.Sign = helpers.RequestSign(signParams)

	jsonBody, _ := json.Marshal(order)
	fmt.Println("Request Body:", string(jsonBody))

	client := createInsecureHTTPClient()

	reqHttp, err := http.NewRequest(
		"POST",
		"https://developer.toutiao.com/api/apps/ecpay/v1/create_order",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	reqHttp.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(reqHttp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	// Parse the response to merge out_order_no
	var douyinResponse map[string]interface{}
	if err := json.Unmarshal(body, &douyinResponse); err != nil {
		// If parsing fails, return original response
		c.Data(resp.StatusCode, "application/json", body)
		return
	}

	// Add out_order_no to the response
	douyinResponse["out_order_no"] = orderNo

	// Return the merged response
	c.JSON(resp.StatusCode, douyinResponse)
}

func PayOrderCallback(c *gin.Context) {
	// This is a placeholder for the callback handler
	// You would typically handle the payment notification here
	fmt.Println("Payment callback received")
	c.JSON(http.StatusOK, gin.H{"message": "Payment callback received"})
}

func PayDouOrder(c *gin.Context) {
	cfg := config.LoadConfig()
	var (
		// 请求时间戳
		timestamp = strconv.FormatInt(time.Now().Unix(), 10)
		// 开发者填入自己的小程序app_id
		appId = cfg.AppID
		// 随机字符串
		nonceStr = helpers.RandStr(10)
		// 应用公钥版本,每次重新上传公钥后需要更新,可通过「开发管理-开发设置-密钥设置」处获取
		keyVersion = "1"
		// 应用私钥,用于加签 重要：1.测试时请修改为开发者自行生成的私钥;2.请勿将示例密钥用于生产环境;3.建议开发者不要将私钥文本写在代码中
		privateKeyStr       = cfg.PrivateKey
		privateKeyBase64Str = base64.StdEncoding.EncodeToString([]byte(privateKeyStr))
		// 生成好的data
		data = "{\"skuList\":[{\"skuId\":\"657\",\"price\":1,\"quantity\":1,\"title\":\"test_title\",\"imageList\":[\"https://xxxx.com/xxxxx.jpg\"],\"type\":301,\"tagGroupId\":\"tag_group_7272625659888058380\"}],\"outOrderNo\":\"test_out_order_no\",\"totalAmount\":1,\"payExpireSeconds\":300,\"orderEntrySchema\":{\"path\":\"\",\"params\":\"\"},\"payNotifyUrl\":\"https://xxxxx/xxx\"}"
		//data = "{\"skuList\":[{\"skuId\":\"1\",\"price\":9999,\"quantity\":1,\"title\":\"标题\",\"imageList\":[\"https://dummyimage.com/234x60\"],\"type\":301,\"tagGroupId\":\"tag_group_7272625659888058380\"}],\"outOrderNo\":\"1213\",\"totalAmount\":9999,\"limitPayWayList\":[],\"payExpireSeconds\":3000,\"orderEntrySchema\":{\"path\":\"page/index/index\",\"params\":\"{\\\"poi\\\":\\\"6601248937917548558\\\",\\\"aweme_useTemplate\\\":1}\"}}"
	)
	byteAuthorization, err := helpers.GetByteAuthorization(privateKeyBase64Str, data, appId, nonceStr, timestamp, keyVersion)
	if err != nil {
		fmt.Println("getByteAuthorization err:", err)
	} else {
		fmt.Println("getByteAuthorization res:", byteAuthorization)
	}
	// return json response and { "auth": byteAuthorization.to_s}
	c.JSON(http.StatusOK, gin.H{"auth": string(byteAuthorization)})
}
