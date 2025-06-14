package models

import (
	"errors"
	"learning-api/config"
	"time"

	openApiSdkClient "github.com/bytedance/douyin-openapi-sdk-go/client"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

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
	User                  User      // One-to-one relationship with User
}

func NewToken() *Token {
	return &Token{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func FindToken(accessToken, refreshToken string) (*Token, error) {
	var token Token
	result := db.Where("access_token = ? AND refresh_token = ?", accessToken, refreshToken).First(&token)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Token not found
		}
		return nil, result.Error
	}
	return &token, nil
}

func FindOrCreateUserToken(data *openApiSdkClient.V2Jscode2sessionResponseData) (token *Token, err error) {

	var user *User

	result := db.Where("open_id = ?", data.Openid).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// User not found, use your custom method to create a new one
			user, err = createNewUserWithToken(data)
			if err != nil {
				return nil, err
			}
			// convert the token in &user.Tokens[0]  and return it

			return &user.Tokens[0], nil // Return the newly created token
		} else {
			return nil, result.Error
		}
	} else {
		// user found, update the user updated_at and session key

		token := NewToken()
		err := token.GenTokenWithDate()
		if err != nil {
			return nil, err
		}
		// delete all the tokens where user_id = user.ID
		if err := db.Where("user_id = ?", user.ID).Delete(&Token{}).Error; err != nil {
			return nil, err
		}

		user.Tokens = []Token{*token}
		user.SessionKey = *data.SessionKey
		user.UpdatedAt = time.Now()
		if err := db.Save(&user).Error; err != nil {
			return nil, err
		}
		return token, nil
	}
}

func NewUserAndToken(data *openApiSdkClient.V2Jscode2sessionResponseData) (*User, error) {
	return createNewUserWithToken(data)
}

func createNewUserWithToken(data *openApiSdkClient.V2Jscode2sessionResponseData) (*User, error) {
	user := &User{
		OpenID:     *data.Openid,
		UnionID:    *data.Unionid,
		SessionKey: *data.SessionKey,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	token := NewToken()
	err := token.GenTokenWithDate()
	if err != nil {
		return nil, err
	}

	user.Tokens = []Token{*token}

	if err := db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (t *Token) GenTokenWithDate() error {

	const accessTokenExpiresIn = 3600      // 1 hour (seconds)
	const refreshTokenExpiresIn = 31536000 // 1 year (seconds)

	t.AccessTokenExpiresIn = accessTokenExpiresIn
	t.RefreshTokenExpiresIn = refreshTokenExpiresIn

	err := t.SetAccessTokenAndRefreshToken(accessTokenExpiresIn, refreshTokenExpiresIn)
	if err != nil {
		return err
	}

	return nil
}

func (t *Token) RefreshToNewToken() (*Token, error) {

	const accessTokenExpiresIn = 3600 // 1 hour (seconds)
	newToken := NewToken()
	newToken.UserID = t.UserID
	newToken.AccessTokenExpiresIn = accessTokenExpiresIn

	newRefreshTokenExpiresIn := t.RefreshTokenExpiresIn - accessTokenExpiresIn

	err := newToken.SetAccessTokenAndRefreshToken(accessTokenExpiresIn, newRefreshTokenExpiresIn)
	if err != nil {
		return nil, err
	}
	if err := db.Save(newToken).Error; err != nil {
		return nil, err
	}

	return newToken, nil
}

func (t *Token) SetAccessTokenAndRefreshToken(accessTokenExpiresIn int, refreshTokenExpiresIn int) error {
	config := config.LoadConfig()
	secret := config.ClientSecret
	t.AccessTokenExpiresIn = accessTokenExpiresIn
	accessToken, err := GenJWTToken(secret, time.Duration(accessTokenExpiresIn)*time.Second)
	if err != nil {
		return err
	}
	t.AccessToken = accessToken
	refreshToken, err := GenJWTToken(secret, time.Duration(refreshTokenExpiresIn)*time.Second)
	if err != nil {
		return err
	}
	t.RefreshToken = refreshToken
	t.RefreshTokenExpiresIn = refreshTokenExpiresIn

	return nil
}

// GenJWTToken generates a JWT token with the specified secret and expiration duration.

// gen_jwt_access_token generates a JWT access token for the user with a given secret and expiration duration.
func GenJWTToken(secret string, expiresIn time.Duration) (string, error) {
	//
	const issuer = "learning" // Replace with your actual issuer
	const audience = "douyin" // Replace with your actual audience
	claims := jwt.MapClaims{
		"exp": time.Now().Add(expiresIn).Unix(),
		"iat": time.Now().Unix(),
		"iss": issuer,   // Replace with your actual issuer
		"aud": audience, // Replace with your actual audience
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString([]byte(secret))
}
