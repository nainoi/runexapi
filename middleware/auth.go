package middleware

import (
	"log"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"thinkdev.app/think/runex/runexapi/api/v1/user"
	config "thinkdev.app/think/runex/runexapi/config"
	"thinkdev.app/think/runex/runexapi/model"
	//"thinkdev.app/think/runex/runexapi/repository"
)

var api user.UserAPI

// Auth middleware
func Auth(userApi user.UserAPI) *jwt.GinJWTMiddleware {
	api = userApi
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte(config.SECRET_KEY),
		Timeout:     time.Hour * 72 * 24 * 365,
		MaxRefresh:  time.Hour * 72 * 24 * 365,
		IdentityKey: config.ID_KEY,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.UserAuth); ok {
				return jwt.MapClaims{
					config.ID_KEY:   v.UserID,
					config.ROLE_KEY: v.Role,
					config.PF:       v.PF,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			if value, ok := claims[config.ID_KEY].(string); ok {
				if userID, err := primitive.ObjectIDFromHex(value); err == nil {
					return &model.UserAuth{
						UserID: userID,
						Role:   claims[config.ROLE_KEY].(string),
						PF:     claims[config.ROLE_KEY].(string),
						// Role: claims[identityKey].(string),
					}
				}
			}
			return &model.UserAuth{
				UserID: primitive.NilObjectID,
				Role:   claims[config.ROLE_KEY].(string),
				PF:     claims[config.PF].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			// if path, exist := c.Get("/api/v1/user/login"); exist {
			// 	log.Println("------------------------")
			// 	log.Println(path)
			// }
			if c.Request.URL.Path == "/api/v1/user/login" {
				user, err := userApi.Login(c)
				if err != nil {
					return "", jwt.ErrMissingLoginValues
				}
				//log.Println(user)
				return &user, nil
			} else {
				user, err := userApi.LoginPD(c)
				if err != nil {
					return "", jwt.ErrMissingLoginValues
				}
				//log.Println(user)
				return &user, nil
			}

			//return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*model.UserAuth); ok {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code": code,
				"msg":  message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
	return authMiddleware
}
