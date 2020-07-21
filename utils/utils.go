package utils

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"math"
	"path"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	// "github.com/gin-gonic/gin"
	config "bitbucket.org/suthisakch/runex/config"
	"golang.org/x/crypto/bcrypt"
)

// HashAndSalt hash password user
func HashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

// GetTokenValue get JWT claims
func GetTokenValue(ctx *gin.Context) (string, string, string) {
	claims := jwt.ExtractClaims(ctx)
	userid := claims[config.ID_KEY].(string)
	role := claims[config.ROLE_KEY].(string)
	pf := claims[config.ROLE_KEY].(string)
	return userid, role, pf
}

func ISAdmin(ctx *gin.Context) bool {

	claims := jwt.ExtractClaims(ctx)
	role := claims[config.ROLE_KEY].(string)
	if role == "ADMIN" {
		return true
	} else {
		return false
	}
}

// EncodeMD5 md5 encryption
func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}

// GetExt get the file ext
func GetExt(fileName string) string {
	return path.Ext(fileName)
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
