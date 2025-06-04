package handlers

import (
	"net/http"
	"testing"

	credential "github.com/bytedance/douyin-openapi-credential-go/client"
	openApiSdkClient "github.com/bytedance/douyin-openapi-sdk-go/client"
	"github.com/gin-gonic/gin"
)

// TokenRequest represents the expected request body for /token
// Only a 'code' param is required

type TokenRequest struct {
	Code string `json:"code" binding:"required"`
}

// PostToken handles POST /token
func PostToken(c *gin.Context) {
	var req TokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code param is required"})
		return
	}
	// TODO: Implement your code-to-token logic here
	c.JSON(http.StatusOK, gin.H{"message": "Token endpoint received code", "code": req.Code})
}

func TestV2Jscode2session(t *testing.T) {
	// 初始化SDK client
	opt := new(credential.Config).
		SetClientKey("tt******"). // 改成自己的app_id
		SetClientSecret("cbs***") // 改成自己的secret
	sdkClient, err := openApiSdkClient.NewClient(opt)
	if err != nil {
		t.Log("sdk init err:", err)
		return
	}

	/* 构建请求参数，该代码示例中只给出部分参数，请用户根据需要自行构建参数值
	   	token:
	   	   1.若用户自行维护token,将用户维护的token赋值给该参数即可
	          2.SDK包中有获取token的函数，请根据接口path在《OpenAPI SDK 总览》文档中查找获取token函数的名字
	            在使用过程中，请注意token互刷问题
	       header:
	          sdk中默认填充content-type请求头，若不需要填充除content-type之外的请求头，删除该参数即可
	*/
	sdkRequest := &openApiSdkClient.V2Jscode2sessionRequest{}
	sdkRequest.SetAnonymousCode("EpFjYuaoOl")
	sdkRequest.SetAppid("tt4233**")
	sdkRequest.SetCode("YAH9JFcEc4")
	sdkRequest.SetSecret("83Soi1UKQ6")

	// sdk调用
	sdkResponse, err := sdkClient.V2Jscode2session(sdkRequest)
	if err != nil {
		t.Log("sdk call err:", err)
		return
	}
	t.Log(sdkResponse)
}
