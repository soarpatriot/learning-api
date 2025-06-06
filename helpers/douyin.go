package helpers

import (
	"fmt"
	"learning-api/config"
	"learning-api/models"

	credential "github.com/bytedance/douyin-openapi-credential-go/client"
	openApiSdkClient "github.com/bytedance/douyin-openapi-sdk-go/client"
)

type ThirdPartyClient interface {
	// You can add fields like baseURL, http.Client, etc.
	Jscode2session(code string) (*models.Token, error)
}

type DouyinClient struct {
}

func fetchConfig() config.Config {
	// Load the configuration from the config package
	return config.LoadConfig()
}

func NewDouyinClient() ThirdPartyClient {
	return &DouyinClient{}
}

func GenerateSdkClient() (sdkClient *openApiSdkClient.Client, error error) {
	config := fetchConfig()
	apiKey := config.ApiSecret
	apiSecret := config.ApiKey
	opt := new(credential.Config).
		SetClientKey(apiKey).
		SetClientSecret(apiSecret)
	return openApiSdkClient.NewClient(opt)
}

func (d *DouyinClient) Jscode2session(code string) (*models.Token, error) {
	sdkClient, err := GenerateSdkClient()
	token := &models.Token{}
	if err != nil {
		return nil, err
	}
	config := fetchConfig()
	sdkRequest := constructSessionRequest(code, config.AppID, config.ApiSecret)

	// sdk调用
	sdkResponse, err := sdkClient.V2Jscode2session(sdkRequest)
	if err != nil || sdkResponse == nil || sdkResponse.ErrNo == nil || *sdkResponse.ErrNo != 0 {
		fmt.Println("sdk call err:", err, " response:", sdkResponse)
		return nil, err
	}
	token, err = token.FindOrCreateUserToken(sdkResponse.Data)
	if err != nil {
		fmt.Println("Error finding or creating user token:", err)
		return nil, err
	}
	return token, nil
}

func constructSessionRequest(code string, appid string, secret string) *openApiSdkClient.V2Jscode2sessionRequest {
	sdkRequest := &openApiSdkClient.V2Jscode2sessionRequest{}

	sdkRequest.SetAppid(appid)
	sdkRequest.SetCode(code)
	sdkRequest.SetSecret(secret)

	return sdkRequest
}
