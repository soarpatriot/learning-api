package models

import (
	"errors"
	"fmt"
	"learning-api/config"
	"time"

	credential "github.com/bytedance/douyin-openapi-credential-go/client"
	openApiSdkClient "github.com/bytedance/douyin-openapi-sdk-go/client"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var db *gorm.DB // Initialize this variable appropriately in your application

// Token represents a token entity
type Token struct {
	ID                    uint      `gorm:"primaryKey" json:"id"`
	UserID                uint      `json:"user_id"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresIn  int       `json:"access_token_expires_in"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresIn int       `json:"refresh_token_expires_in"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// DouyinClient interface for testability
type DouyinClient interface {
	V2Jscode2session(req *openApiSdkClient.V2Jscode2sessionRequest) (*openApiSdkClient.V2Jscode2sessionResponse, error)
}

// generateSdkClientFunc returns a DouyinClient (for production or test)
var generateSdkClientFunc = func() (DouyinClient, error) {
	return generateSdkClient()
}

func generateSdkClient() (DouyinClient, error) {
	config := fetchConfig()
	apiKey := config.ApiSecret
	apiSecret := config.ApiKey
	opt := new(credential.Config).
		SetClientKey(apiKey).
		SetClientSecret(apiSecret)
	return openApiSdkClient.NewClient(opt)
}

func fetchConfig() config.Config {
	// Load the configuration from the config package
	return config.LoadConfig()
}

func Jscode2session(code string) (string, error) {
	sdkClient, err := generateSdkClientFunc()
	if err != nil {
		return "", err
	}
	config := fetchConfig()
	sdkRequest := constructSessionRequest(code, config.AppID, config.ApiSecret)

	// sdk调用
	sdkResponse, err := sdkClient.V2Jscode2session(sdkRequest)
	if err != nil && sdkResponse.ErrNo != nil && *sdkResponse.ErrNo != 0 {
		fmt.Println("sdk call err:", err, " response:", sdkResponse)
		return "", err
	}
	token, err := findOrCreateUserToken(sdkResponse.Data)
	if err != nil {
		fmt.Println("Error finding or creating user token:", err)
		return "", err
	}
	return token.AccessToken, nil
}

func constructSessionRequest(code string, appid string, secret string) *openApiSdkClient.V2Jscode2sessionRequest {
	sdkRequest := &openApiSdkClient.V2Jscode2sessionRequest{}

	sdkRequest.SetAppid(appid)
	sdkRequest.SetCode(code)
	sdkRequest.SetSecret(secret)

	return sdkRequest
}

func findOrCreateUserToken(data *openApiSdkClient.V2Jscode2sessionResponseData) (token *Token, err error) {

	var user *User
	result := db.Where("open_id = ?", data.Openid).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// User not found, use your custom method to create a new one
			user, err = createNewUserWithToken(data)
			if err != nil {
				return nil, err
			}
			return &user.Tokens[0], nil // Return the newly created token
		} else {
			return nil, result.Error
		}
	} else {
		// user found, update the user updated_at and session key

		token = generateToken(*data.Openid)
		user.Tokens = []Token{*token}
		user.SessionKey = *data.SessionKey
		user.UpdatedAt = time.Now()
		if err := db.Save(&user).Error; err != nil {
			return nil, err
		}
	}
	return token, nil
}

func createNewUserWithToken(data *openApiSdkClient.V2Jscode2sessionResponseData) (*User, error) {
	user := &User{
		OpenID:     *data.Openid,
		UnionID:    *data.Unionid,
		SessionKey: *data.SessionKey,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	token := generateToken(*data.Openid)
	user.Tokens = []Token{*token}

	if err := db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func generateToken(openID string) *Token {
	config := fetchConfig()
	const accessTokenExpiresIn = 3600      // 1 hour (seconds)
	const refreshTokenExpiresIn = 31536000 // 1 year (seconds)
	secret := config.ApiSecret
	token := &Token{
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		AccessTokenExpiresIn:  accessTokenExpiresIn,
		RefreshTokenExpiresIn: refreshTokenExpiresIn,
	}
	token.AccessToken, _ = token.genJWTToken(secret, openID, time.Duration(accessTokenExpiresIn)*time.Second)
	token.RefreshToken, _ = token.genJWTToken(secret, openID, time.Duration(refreshTokenExpiresIn)*time.Second)
	return token
}

// gen_jwt_access_token generates a JWT access token for the user with a given secret and expiration duration.
func (t *Token) genJWTToken(secret string, openID string, expiresIn time.Duration) (string, error) {
	//
	const issuer = "learning" // Replace with your actual issuer
	const audience = "douyin" // Replace with your actual audience
	claims := jwt.MapClaims{
		"open_id": openID,
		"exp":     time.Now().Add(expiresIn).Unix(),
		"iat":     time.Now().Unix(),
		"iss":     issuer,   // Replace with your actual issuer
		"aud":     audience, // Replace with your actual audience
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
