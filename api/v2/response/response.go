package response

import (
	"github.com/gin-gonic/gin"
)

//Gin struct
type Gin struct {
	C *gin.Context
}

//ResponseOAuth struct
type ResponseOAuth struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

//Response struct
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Response setting gin.JSON
func (g *Gin) Response(httpCode int, errMsg string, data interface{}) {
	g.C.JSON(httpCode, Response{
		Code: httpCode,
		Msg:  errMsg,
		Data: data,
	})
	//g.C.Abort()
	return
}
