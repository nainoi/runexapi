package response

import (
	"github.com/gin-gonic/gin"
)

//Gin struct
type Gin struct {
	C *gin.Context
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
	return
}
