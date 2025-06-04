package models

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
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

// gen_jwt_access_token generates a JWT access token for the user with a given secret and expiration duration.
func (t *Token) GenJWTAccessToken(secret string, expiresIn time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": t.UserID,
		"exp":     time.Now().Add(expiresIn).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// gen_refresh_token generates a secure random refresh token string.
func (t *Token) GenRefreshToken(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Define V2Jscode2sessionResponse according to the expected SDK response structure
type V2Jscode2sessionResponse struct {
	// Add fields as per the actual response structure
	OpenID     string `json:"open_id"`
	SessionKey string `json:"session_key"`
	ErrCode    int    `json:"err_code"`
	ErrMsg     string `json:"err_msg"`
}

func generate_sdk_client() (*openApiSdkClient.Client, error) {
	// 初始化SDK client
	opt := new(credential.Config).
		SetClientKey("tt******"). // 改成自己的app_id
		SetClientSecret("cbs***") // 改成自己的secret
	return openApiSdkClient.NewClient(opt)
}

func Jscode2session(code string) (string, error) {
	// 初始化SDK client
	sdkClient, err := generate_sdk_client()
	if err != nil {
		return "", err
	}

	user := temp_user()
	token := generate_a_token("83Soi1UKQ6", user.ID)
	sdkRequest := construct_session_request(code, "tt4233**", "83Soi1UKQ6")

	// sdk调用
	sdkResponse, err := sdkClient.V2Jscode2session(sdkRequest)
	if err != nil && sdkResponse.ErrNo != nil && *sdkResponse.ErrNo != 0 {
		fmt.Println("sdk call err:", err, " response:", sdkResponse)
		db.Delete(&user) // Clean up the temporary user if there's an error
		return "", err
	}
	user.OpenID = *sdkResponse.Data.Openid
	user.UnionID = *sdkResponse.Data.Unionid
	user.SessionKey = *sdkResponse.Data.SessionKey
	db.Save(user)     // Update the user in the database
	db.Create(&token) // Save the token to the database
	return token.AccessToken, nil
}

func construct_session_request(code string, appid string, secret string) *openApiSdkClient.V2Jscode2sessionRequest {
	sdkRequest := &openApiSdkClient.V2Jscode2sessionRequest{}

	sdkRequest.SetAppid(appid)
	sdkRequest.SetCode(code)
	sdkRequest.SetSecret(secret)

	// sdkRequest.AccessToken = accessToken // Removed: no such field in SDK struct
	return sdkRequest
}

func temp_user() *User {
	// Create a temporary user for testing purposes
	// generate a user and save it to the database
	user := &User{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if db != nil {
		db.Create(&user)
	}
	return user
}

func generate_a_token(secret string, userID uint) *Token {

	// Save user to the database (omitted for brevity)

	token := &Token{
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),

		AccessTokenExpiresIn: 3600,

		RefreshTokenExpiresIn: 7200,
	}
	token.AccessToken, _ = token.GenJWTAccessToken(secret, time.Hour)
	token.RefreshToken, _ = token.GenRefreshToken(32)

	return token
}
