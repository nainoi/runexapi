package oauth

import (
	"context"
	"fmt"

	//"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"thinkdev.app/think/runex/runexapi/api/v1/user"
	"thinkdev.app/think/runex/runexapi/api/v2/response"
	config "thinkdev.app/think/runex/runexapi/config"
	"thinkdev.app/think/runex/runexapi/config/db"
	"thinkdev.app/think/runex/runexapi/logger"
	"thinkdev.app/think/runex/runexapi/model"
)

var (
	// ErrMissingField is use to error
	ErrMissingField = "Error missing %v"
	// AtJwtKey is used to create the Access token signature
	AtJwtKey = []byte(config.SECRET_KEY)
	// RtJwtKey is used to create the refresh token signature
	RtJwtKey = []byte(config.RE_SECRET_KEY)
)

//RefreshTokenRequestBody refresh token key
type RefreshTokenRequestBody struct {
	RefreshToken string `json:"refresh_token"`
}

//AccessTokenRequestBody refresh token key
type AccessTokenRequestBody struct {
	AccessToken string `json:"access_token"`
}

var api user.UserAPI

// User Struct
type User struct {
	U *model.User
}

//GetValuesToken in claim
func GetValuesToken(c *gin.Context) (string, string) {
	clientToken := c.GetHeader("Authorization")
	extractedToken := strings.Split(clientToken, "Bearer ")
	// Verify if the format of the token is correct
	if len(extractedToken) == 2 {
		clientToken = strings.TrimSpace(extractedToken[1])
	} else {
		return "", ""
	}
	valid, id, role, err := ValidateToken(clientToken)

	if valid == false || err != nil {
		return "", ""
	}
	return id, role
}

// Validate if the fields are available
func (u *User) Validate() error {
	if u.U.Email == "" {
		return fmt.Errorf(ErrMissingField, "Email")
	}
	// if u.Password == "" {
	// 	return fmt.Errorf(ErrMissingField, "Password")
	// }
	return nil
}

//IsBlacklisted check backlist token
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
		var (
			res = response.Gin{C: c}
		)
		clientToken := c.GetHeader("Authorization")
		if clientToken == "" {
			//logger.Logger.Errorf("Authorization token was not provided")
			res.Response(http.StatusUnauthorized, "Authorization Token is required", nil)
			c.Abort()
			return
		}

		claims := jwt.MapClaims{}

		extractedToken := strings.Split(clientToken, "Bearer ")

		// Verify if the format of the token is correct
		if len(extractedToken) == 2 {
			clientToken = strings.TrimSpace(extractedToken[1])
		} else {
			//logger.Logger.Errorf("Incorrect Format of Authn Token")
			res.Response(http.StatusBadRequest, "Incorrect Format of Authorization Token ", nil)
			c.Abort()
			return
		}

		foundInBlacklist := IsBlacklisted(extractedToken[1])

		if foundInBlacklist == true {
			//logger.Logger.Infof("Found in Blacklist")
			res.Response(http.StatusUnauthorized, "Invalid Token", nil)
			c.Abort()
			return
		}

		// Parse the claims
		parsedToken, err := jwt.ParseWithClaims(clientToken, claims, func(token *jwt.Token) (interface{}, error) {
			return AtJwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				//logger.Logger.Errorf("Invalid Token Signature")
				res.Response(http.StatusUnauthorized, "Invalid Token Signature", nil)
				c.Abort()
				return
			}

			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					// Token is either expired or not active yet
					res.Response(444, "Token expire", nil)
					c.Abort()
					return
				}
			}

			res.Response(http.StatusBadRequest, "Bad Request", nil)
			c.Abort()
			return
		}

		if !parsedToken.Valid {
			//logger.Logger.Errorf("Invalid Token")
			res.Response(http.StatusUnauthorized, "Invalid Token", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}

// GenerateTokenPair creates and returns a new set of access_token and refresh_token.
func GenerateTokenPair(u model.User) (string, string, error) {

	tokenString, err := GenerateAccessToken(u)
	if err != nil {
		return "", "", err
	}

	// Create Refresh token, this will be used to get new access token.
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshToken.Header["rnx"] = "signin_2"

	// Expiration time is 180 minutes
	expirationTimeRefreshToken := time.Now().Add(24 * 30 * 3 * time.Hour).Unix()

	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims[config.ID_KEY] = u.UserID
	rtClaims["exp"] = expirationTimeRefreshToken

	refreshTokenString, err := refreshToken.SignedString(RtJwtKey)
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, nil
}

// ValidateRefreshToken is used to validate both access_token and refresh_token. It is done based on the "Key ID" provided by the JWT
func ValidateRefreshToken(tokenString string) (bool, string, string, error) {

	var key []byte

	var keyID string

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {

		keyID = token.Header["rnx"].(string)
		// If the "rnx" (Key ID) is equal to signin_1, then it is compared against access_token secret key, else if it
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

	return true, claims[config.ID_KEY].(string), "", nil
}

// ValidateToken is used to validate both access_token and refresh_token. It is done based on the "Key ID" provided by the JWT
func ValidateToken(tokenString string) (bool, string, string, error) {

	var key []byte

	var keyID string

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {

		keyID = token.Header["rnx"].(string)
		// If the "rnx" (Key ID) is equal to signin_1, then it is compared against access_token secret key, else if it
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

	return true, claims[config.ID_KEY].(string), claims[config.ROLE_KEY].(string), nil
}

// InvalidateToken method marks the access_token as invalid. This token cannot be used for future authentication
func InvalidateToken(tokenString string) error {

	var key []byte

	var keyID string

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {

		keyID = token.Header["rnx"].(string)
		// If the "rnx" (Key ID) is equal to signin_1, then it is compared against access_token secret key, else if it
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
		//logger.Logger.Infof("User Logged Out")
		return nil
	}

	return status.Err()
}

// GenerateAccessToken method creats a new access token when the user logs in by providing username and password
func GenerateAccessToken(u model.User) (string, error) {
	// Declare the expiration time of the access token
	// Here the expiration is 60 minutes
	expirationTimeAccessToken := time.Now().Add(24 * 60 * time.Minute).Unix()
	//expirationTimeAccessToken := time.Now().Add(1 * time.Minute).Unix()

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.New(jwt.SigningMethodHS256)
	token.Header["rnx"] = "signin_1"
	claims := token.Claims.(jwt.MapClaims)
	claims[config.ID_KEY] = u.UserID
	claims[config.ROLE_KEY] = u.Role
	claims["exp"] = expirationTimeAccessToken

	// Create the JWT string
	tokenString, err := token.SignedString(AtJwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
