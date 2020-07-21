package middleware

import (
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Content-Type", "application/json")
		//c.AbortWithStatus(http.StatusOK)

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(200)
            return
        }

        c.Next()
    }
}