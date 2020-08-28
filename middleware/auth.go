package middleware

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	jwt_app "github.com/appleboy/gin-jwt/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"thinkdev.app/think/runex/runexapi/api/v1/user"
	config "thinkdev.app/think/runex/runexapi/config"
	"thinkdev.app/think/runex/runexapi/config/db"
	"thinkdev.app/think/runex/runexapi/logger"
	"thinkdev.app/think/runex/runexapi/model"
	//"thinkdev.app/think/runex/runexapi/repository"
)

var (
	ErrMissingField = "Error missing %v"
	// AtJwtKey is used to create the Access token signature
	AtJwtKey = []byte("my_secret_key")
	// RtJwtKey is used to create the refresh token signature
	RtJwtKey = []byte("my_secret_key_2")
)

type RefreshTokenRequestBody struct {
	RefreshToken string `json:"refresh_token"`
}

type AccessTokenRequestBody struct {
	AccessToken string `json:"access_token"`
}

var api user.UserAPI
var u model.User

// Auth middleware
func Auth(userApi user.UserAPI) *jwt.GinJWTMiddleware {
	api = userApi
	authMiddleware, err := jwt_app.New(&jwt_app.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte(config.SECRET_KEY),
		Timeout:     time.Hour * 72 * 24 * 365,
		MaxRefresh:  time.Hour * 72 * 24 * 365,
		IdentityKey: config.ID_KEY,
		PayloadFunc: func(data interface{}) jwt_app.MapClaims {
			if v, ok := data.(*model.UserAuth); ok {
				return jwt_app.MapClaims{
					config.ID_KEY:   v.UserID,
					config.ROLE_KEY: v.Role,
					config.PF:       v.PF,
				}
			}
			return jwt_app.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt_app.ExtractClaims(c)
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
					return "", jwt_app.ErrMissingLoginValues
				}
				//log.Println(user)
				return &user, nil
			} else {
				user, err := userApi.LoginPD(c)
				if err != nil {
					return "", jwt_app.ErrMissingLoginValues
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

// Validate if the fields are available
func (u *model.User) Validate() error {
	if u.FirstName == "" {
		return fmt.Errorf(ErrMissingField, "FirstName")
	}
	if u.LastName == "" {
		return fmt.Errorf(ErrMissingField, "LastName")
	}
	if u.Email == "" {
		return fmt.Errorf(ErrMissingField, "Email")
	}
	if u.Password == "" {
		return fmt.Errorf(ErrMissingField, "Password")
	}
	return nil
}

func (u *model.User) NewSalt() {
	// h := sha1.New()
	// io.WriteString(h, strconv.Itoa(int(time.Now().UnixNano())))
	// u.Salt = fmt.Sprintf("%x", h.Sum(nil))
}

func CalculatePassHash(pass, salt string) string {
	h := sha1.New()
	io.WriteString(h, salt)
	io.WriteString(h, pass)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func IsBlacklisted(tokenString string) bool {

	status := db.RedisClient.Get(context.Background(), tokenString)

	val, _ := status.Result()

	if val == "" {
		return false
	}

	return true
}

// AuthMiddleware checks if the JWT sent is valid or not. This function is involked for every API route that needs authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.GetHeader("Authorization")
		if clientToken == "" {
			logger.Logger.Errorf("Authorization token was not provided")
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Authorization Token is required"})
			c.Abort()
			return
		}

		claims := jwt.MapClaims{}

		extractedToken := strings.Split(clientToken, "Bearer ")

		// Verify if the format of the token is correct
		if len(extractedToken) == 2 {
			clientToken = strings.TrimSpace(extractedToken[1])
		} else {
			logger.Logger.Errorf("Incorrect Format of Authn Token")
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Incorrect Format of Authorization Token "})
			c.Abort()
			return
		}

		foundInBlacklist := IsBlacklisted(extractedToken[1])

		if foundInBlacklist == true {
			logger.Logger.Infof("Found in Blacklist")
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Invalid Token"})
			c.Abort()
			return
		}

		// Parse the claims
		parsedToken, err := jwt.ParseWithClaims(clientToken, claims, func(token *jwt.Token) (interface{}, error) {
			return AtJwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				logger.Logger.Errorf("Invalid Token Signature")
				c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Invalid Token Signature"})
				c.Abort()
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Bad Request"})
			c.Abort()
			return
		}

		if !parsedToken.Valid {
			logger.Logger.Errorf("Invalid Token")
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Invalid Token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GenerateTokenPair creates and returns a new set of access_token and refresh_token.
func GenerateTokenPair(username string, uuid string) (string, string, error) {

	tokenString, err := GenerateAccessToken(username, uuid)
	if err != nil {
		return "", "", err
	}

	// Create Refresh token, this will be used to get new access token.
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshToken.Header["kid"] = "signin_2"

	// Expiration time is 180 minutes
	expirationTimeRefreshToken := time.Now().Add(180 * time.Minute).Unix()

	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = uuid
	rtClaims["exp"] = expirationTimeRefreshToken

	refreshTokenString, err := refreshToken.SignedString(RtJwtKey)
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, nil
}

// ValidateToken is used to validate both access_token and refresh_token. It is done based on the "Key ID" provided by the JWT
func ValidateToken(tokenString string) (bool, string, string, error) {

	var key []byte

	var keyID string

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {

		keyID = token.Header["kid"].(string)
		// If the "kid" (Key ID) is equal to signin_1, then it is compared against access_token secret key, else if it
		// is equal to signin_2 , it is compared against refresh_token secret key.
		if keyID == "signin_1" {
			key = AtJwtKey
		} else if keyID == "signin_2" {
			key = RtJwtKey
		}
		return key, nil
	})

	// Check if signatures are valid.
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			logger.Logger.Errorf("Invalid Token Signature")
			return false, "", keyID, err
		}
		return false, "", keyID, err
	}

	if !token.Valid {
		logger.Logger.Errorf("Invalid Token")
		return false, "", keyID, err
	}

	return true, claims["sub"].(string), keyID, nil
}

// InvalidateToken method marks the access_token as invalid. This token cannot be used for future authentication
func InvalidateToken(tokenString string) error {

	var key []byte

	var keyID string

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {

		keyID = token.Header["kid"].(string)
		// If the "kid" (Key ID) is equal to signin_1, then it is compared against access_token secret key, else if it
		// is equal to signin_2 , it is compared against refresh_token secret key.
		if keyID == "signin_1" {
			key = AtJwtKey
		} else if keyID == "signin_2" {
			key = RtJwtKey
		}
		return key, nil
	})

	// Check if signatures are valid.
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			logger.Logger.Errorf("Invalid Token Signature")
			return err
		}
		return err
	}

	if !token.Valid {
		logger.Logger.Errorf("Invalid Token")
		return err
	}

	// @TODO - Fix the expiration time
	status := db.RedisClient.Set(context.Background(), tokenString, tokenString, 0)

	// if status.Err() != nil {
	// 	logger.Logger.Errorf("Could not set value in Redis")

	// }

	val, err := status.Result()

	if val == "OK" {
		logger.Logger.Infof("User Logged Out")
		return nil
	}

	return status.Err()
}

// GenerateAccessToken method creats a new access token when the user logs in by providing username and password
func GenerateAccessToken(username string, uuid string) (string, error) {
	// Declare the expiration time of the access token
	// Here the expiration is 60 minutes
	expirationTimeAccessToken := time.Now().Add(60 * time.Minute).Unix()

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.New(jwt.SigningMethodHS256)
	token.Header["kid"] = "signin_1"
	claims := token.Claims.(jwt.MapClaims)
	claims["Username"] = username
	claims["exp"] = expirationTimeAccessToken
	claims["sub"] = uuid

	// Create the JWT string
	tokenString, err := token.SignedString(AtJwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
