package utils

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"math"
	"path"
	"regexp"
	"strings"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	// "github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	config "thinkdev.app/think/runex/runexapi/config"
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

//ISAdmin check user is admin
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

//StringToSlug convert event name to slug
func StringToSlug(s string) string {
	str := []byte(strings.ToLower(s))

	// convert all spaces to dash
	regE := regexp.MustCompile("[[:space:]]")
	str = regE.ReplaceAll(str, []byte("-"))

	// remove all blanks such as tab
	regE = regexp.MustCompile("[[:blank:]]")
	str = regE.ReplaceAll(str, []byte(""))

	// remove all punctuations with the exception of dash

	regE = regexp.MustCompile("[!/:-@[-`{-~]")
	str = regE.ReplaceAll(str, []byte(""))

	regE = regexp.MustCompile("/[^\x20-\x7F]/")
	str = regE.ReplaceAll(str, []byte(""))

	regE = regexp.MustCompile("`&(amp;)?#?[a-z0-9]+;`i")
	str = regE.ReplaceAll(str, []byte("-"))

	regE = regexp.MustCompile("`&([a-z])(acute|uml|circ|grave|ring|cedil|slash|tilde|caron|lig|quot|rsquo);`i")
	str = regE.ReplaceAll(str, []byte("\\1"))

	regE = regexp.MustCompile("`[^a-z0-9]`i")
	str = regE.ReplaceAll(str, []byte("-"))

	regE = regexp.MustCompile("`[-]+`")
	str = regE.ReplaceAll(str, []byte("-"))

	strReplaced := strings.Replace(string(str), "&", "", -1)
	strReplaced = strings.Replace(strReplaced, `"`, "", -1)
	strReplaced = strings.Replace(strReplaced, "&", "-", -1)
	strReplaced = strings.Replace(strReplaced, "--", "-", -1)

	if strings.HasPrefix(strReplaced, "-") || strings.HasPrefix(strReplaced, "--") {
		strReplaced = strings.TrimPrefix(strReplaced, "-")
		strReplaced = strings.TrimPrefix(strReplaced, "--")
	}

	if strings.HasSuffix(strReplaced, "-") || strings.HasSuffix(strReplaced, "--") {
		strReplaced = strings.TrimSuffix(strReplaced, "-")
		strReplaced = strings.TrimSuffix(strReplaced, "--")
	}

	// normalize unicode strings and remove all diacritical/accents marks
	// see https://www.socketloop.com/tutorials/golang-normalize-unicode-strings-for-comparison-purpose

	//t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.)
	//slug, _, err := transform.String(t, strReplaced)

	//if err != nil {
	//        return "", err
	//}

	return strings.TrimSpace(strReplaced)
}
