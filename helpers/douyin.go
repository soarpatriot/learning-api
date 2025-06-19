// package helpers

// import (
// 	"fmt"
// 	"learning-api/config"
// 	"learning-api/models"

// 	credential "github.com/bytedance/douyin-openapi-credential-go/client"
// 	openApiSdkClient "github.com/bytedance/douyin-openapi-sdk-go/client"
// )

// type ThirdPartyClient interface {
// 	// You can add fields like baseURL, http.Client, etc.
// 	Jscode2session(code string, anonymousCode string) (*models.Token, error)
// }

// type DouyinClient struct {
// }

// func fetchConfig() config.Config {
// 	// Load the configuration from the config package
// 	return config.LoadConfig()
// }

// func NewDouyinClient() ThirdPartyClient {
// 	return &DouyinClient{}
// }

// func GenerateSdkClient() (sdkClient *openApiSdkClient.Client, error error) {
// 	config := fetchConfig()
// 	clintKey := config.ClientKey
// 	clientSecret := config.ClientSecret
// 	opt := new(credential.Config).
// 		SetClientKey(clintKey).
// 		SetClientSecret(clientSecret)

// 	return openApiSdkClient.NewClient(opt)
// }

// func (d *DouyinClient) Jscode2session(code string, anonymousCode string) (*models.Token, error) {
// 	fmt.Println("start to call douyin sdk jscode2session with code")
// 	sdkClient, err := GenerateSdkClient()

// 	if err != nil {
// 		fmt.Println("generate sdk client error:", err)
// 		// Handle the error appropriately, maybe return a custom error or nil
// 		return nil, err
// 	}
// 	config := fetchConfig()
// 	fmt.Println("app id:", config.AppID, " app secret: ", config.AppSecret)

// 	fmt.Println("code:", code, " anonymous code: ", anonymousCode)
// 	sdkRequest := constructSessionRequest(code, anonymousCode, config.AppID, config.AppSecret)

// 	// sdk调用
// 	sdkResponse, err := sdkClient.V2Jscode2session(sdkRequest)
// 	if err != nil || sdkResponse == nil || sdkResponse.ErrNo == nil || *sdkResponse.ErrNo != 0 {

// 		fmt.Println("sdk call err:", err, " response:")
// 		return nil, err
// 	}

// 	token, err := models.FindOrCreateUserToken(sdkResponse.Data)
// 	if err != nil {
// 		fmt.Println("Error finding or creating user token:", err)
// 		return nil, err
// 	}
// 	return token, nil
// }

// func constructSessionRequest(code string, anonymousCode string, appid string, secret string) *openApiSdkClient.V2Jscode2sessionRequest {
// 	sdkRequest := &openApiSdkClient.V2Jscode2sessionRequest{}

// 	sdkRequest.SetAppid(appid)
// 	sdkRequest.SetCode(code)
// 	sdkRequest.SetSecret(secret)

// 	return sdkRequest
// }
