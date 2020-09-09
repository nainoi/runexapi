package user

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"thinkdev.app/think/runex/runexapi/api/v2/response"
	"thinkdev.app/think/runex/runexapi/api/v2/tracer"
	"thinkdev.app/think/runex/runexapi/logger"
	"thinkdev.app/think/runex/runexapi/middleware/oauth"
	"thinkdev.app/think/runex/runexapi/model/v2"
	"thinkdev.app/think/runex/runexapi/repository/v2"

	//stdopentracing "github.com/opentracing/opentracing-go"
	tracelog "github.com/opentracing/opentracing-go/log"
)

//API struct for user repository
type API struct {
	UserRepo repository.UserRepository
}

// VerifyAuthToken checks to see if the JWT was present in blacklist table and validates it's authenticity
func (api API) VerifyAuthToken(c *gin.Context) {

	var accessTokenRequest oauth.AccessTokenRequestBody
	var (
		res = response.Gin{C: c}
	)

	err := c.ShouldBindJSON(&accessTokenRequest)
	if err != nil {
		message := err.Error()
		res.Response(http.StatusBadRequest, message, nil)
		return
	}

	foundInBlacklist := oauth.IsBlacklisted(accessTokenRequest.AccessToken)

	if foundInBlacklist == true {
		logger.Logger.Infof("Found in Blacklist")
		res.Response(http.StatusUnauthorized, "Invalid Token", nil)
		c.Abort()
		return
	}

	valid, _, key, err := oauth.ValidateToken(accessTokenRequest.AccessToken)
	if valid == false || err != nil {
		message := err.Error()
		logger.Logger.Errorf(message)
		res.Response(http.StatusUnauthorized, "Invalid Key. User Not Authorized", nil)
		c.Abort()
		return
	}

	// Make sure that key passed was not a refresh token
	if key != "signin_1" {
		logger.Logger.Errorf("Invalid Key Type")
		res.Response(http.StatusUnauthorized, "Provide a valid access token", nil)
		c.Abort()
		return
	}

	// Send StatusOK to indicate the access token was valid
	logger.Logger.Infof("Successfully verified user")
	res.Response(http.StatusOK, "Token Valid. User Authorized", nil)
}

// Paths Information

// RefreshAccessToken godoc
// @Summary Provides a JSON Web Token
// @Description Authenticates a user and provides a JWT to refresh Authorize API calls
// @ID Authentication
// @Consume application/x-www-form-urlencoded
// @Produce json
// @Param refresh_token formData string true "RefreshToken"
// @Success 200 {object} response.ResponseOAuth
// @Failure 401 {object} response.Response
// @Router /refreshAccessToken [post]
func (api API) RefreshAccessToken(c *gin.Context) {

	var tokenRequest oauth.RefreshTokenRequestBody
	var (
		res = response.Gin{C: c}
	)

	err := c.ShouldBindJSON(&tokenRequest)
	if err != nil {
		message := err.Error()
		res.Response(http.StatusBadRequest, message, nil)
		c.Abort()
		return
	}

	valid, id, _, err := oauth.ValidateRefreshToken(tokenRequest.RefreshToken)
	if valid == false || err != nil {
		message := err.Error()
		res.Response(http.StatusUnauthorized, message, nil)
		c.Abort()
		return
	}

	if valid == true && id != "" {

		//var user model.User

		// Retreive the username from users DB. This will verify if the user ID passed with JWT was legit or not.
		p, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			message := err.Error()
			res.Response(http.StatusBadRequest, message, nil)
			c.Abort()
			return
		}

		u, error := api.UserRepo.GetUser(p)

		if error != nil {
			message := "User " + error.Error()
			logger.Logger.Errorf(message)
			res.Response(http.StatusBadRequest, "Invalid refresh token", nil)
			c.Abort()
			return
		}

		newToken, err := oauth.GenerateAccessToken(u)
		if err != nil {
			logger.Logger.Errorf(err.Error())
			res.Response(http.StatusBadRequest, "Cannot Generate New Access Token", nil)
			c.Abort()
			return
		}
		var (
			responseJWT = response.ResponseOAuth{
				AccessToken: newToken,
				RefreshToken: tokenRequest.RefreshToken,
			}
		)
		res.Response(http.StatusOK, "success", responseJWT)
		c.Abort()
		return
	}

	res.Response(http.StatusBadRequest, "Error Found ", nil)
}

//GetUser get user info
func (api API) GetUser(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)
	userID, _ := oauth.GetValuesToken(c)
	if userID != "" {
		p, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			message := err.Error()
			res.Response(http.StatusBadRequest, message, nil)
			c.Abort()
			return
		}
		user, err := api.UserRepo.GetUser(p)
		if err != nil {
			log.Println("error GetUserByIDHandler", err.Error())
			res.Response(http.StatusInternalServerError, "user not found", nil)
			c.Abort()
			return
		}
		res.Response(http.StatusOK, "success", user)
		return
	}
}

//AddUser with separate user
func (api API) AddUser(c *gin.Context) {
	var userProvider model.UserProvider
	// span, err := tracer.CreateTracerAndSpan("login", c)
	// if err != nil {
	// 	logger.Logger.Errorf(err.Error())
	// }
	err := c.ShouldBindJSON(&userProvider)
	if err != nil {
		fmt.Println(err.Error())
		//tracer.OnErrorLog(span, err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		c.Abort()
		return
	}
	var (
		res = response.Gin{C: c}
	)
	_, isMail, isProvider := api.UserRepo.CheckEmail(userProvider)
	if !isMail && !isProvider {
		//span, err := tracer.CreateTracerAndSpan("check_email", c)
		/*if err != nil {
			fmt.Println(err.Error())
		}
		for _, s := range u.Provider {
			if s.ProviderID == userProvider.ProviderID && s.ProviderName == userProvider.Provider {
				accessToken, refreshToken, err := oauth.GenerateTokenPair(u)
				if err != nil || accessToken == "" || refreshToken == "" {
					// Return if there is an error in creating the JWT return an internal server error
					//tracer.OnErrorLog(span, err)
					res.Response(http.StatusInternalServerError, "Could not generate token", nil)
					return
				}
				// span.LogFields(
				// 	tracelog.String("event", "success"),
				// 	tracelog.String("message", "returned token"),
				// 	tracelog.Int("status", http.StatusOK),
				// )
				res.Response(http.StatusOK, "login success", gin.H{"access_token": accessToken, "refresh_token": refreshToken})
				return
			}
		}*/

		//span, err = tracer.CreateTracerAndSpan("add user", c)
		u2, err := api.UserRepo.AddUser(userProvider)
		if err != nil {
			//tracer.OnErrorLog(span, err)
			log.Println("error AddUserHandeler", err.Error())
			res.Response(http.StatusInternalServerError, "add user error", nil)
			c.Abort()
			return
		}
		accessToken, refreshToken, err := oauth.GenerateTokenPair(u2)
		if err != nil || accessToken == "" || refreshToken == "" {
			// Return if there is an error in creating the JWT return an internal server error
			//tracer.OnErrorLog(span, err)
			res.Response(http.StatusInternalServerError, "Could not generate token", nil)
			c.Abort()
			return
		}
		// span.LogFields(
		// 	tracelog.String("event", "success"),
		// 	tracelog.String("message", "returned token"),
		// 	tracelog.Int("status", http.StatusOK),
		// )
		res.Response(http.StatusOK, "login success", gin.H{"access_token": accessToken, "refresh_token": refreshToken})
		return
	}

	res.Response(http.StatusAlreadyReported, "Email in use", nil)
	return

}

//UpdateUser update user info
func (api API) UpdateUser(c *gin.Context) {
}

// GetUsers accepts a context and returns all the users in json format
/*func GetUsers(c *gin.Context) {
	var users []oauth.UserResponse
	span, err := tracer.CreateTracerAndSpan("get_all_users", c)

	if err != nil {
		logger.Logger.Errorf(err.Error())
	}

	logger.Logger.Infof("Retrieving All Users")

	error := db.Collection.Find(nil).All(&users)

	if error != nil {
		tracer.OnErrorLog(span, error)
		message := "Users " + error.Error()
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": message})
		return
	}

	span.LogFields(
		tracelog.String("event", "success"),
		tracelog.Int("status", http.StatusOK),
	)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": users})
}*/

// GetUser accepts context, User ID as param and returns user info
/*func GetUser(c *gin.Context) {
	var user oauth.UserResponse

	span, err := tracer.CreateTracerAndSpan("get_user", c)

	if err != nil {
		logger.Logger.Errorf(err.Error())
	}

	userID := c.Param("id")

	if bson.IsObjectIdHex(userID) {

		error := db.Collection.FindId(bson.ObjectIdHex(userID)).One(&user)

		if error != nil {
			tracer.OnErrorLog(span, error)
			message := "User " + error.Error()
			c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": message})
			return
		}
	} else {
		span.LogFields(
			tracelog.String("event", "error"),
			tracelog.String("message", "Incorrect Format for UserID"),
		)
		message := "Incorrect Format for UserID"
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": message})
		return
	}

	span.LogFields(
		tracelog.String("event", "success"),
		tracelog.Int("status", http.StatusOK),
	)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": user})
}*/

// RegisterUser accepts context and inserts data to the db
/*func RegisterUser(c *gin.Context) {

	var user oauth.User

	error := c.ShouldBindJSON(&user)

	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Incorrect Field Name(s)/ Value(s)"})
		return
	}

	error = user.Validate()

	if error != nil {
		message := "User " + error.Error()
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": message})
		return
	}

	// Inserts ID for the user object
	user.ID = bson.NewObjectId()

	user.Password = auth.CalculatePassHash(user.Password, user.Salt)

	error = db.Collection.Insert(&user)

	if error != nil {
		message := "User " + error.Error()
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": message})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "User created successfully!", "resourceId": user.ID})

}*/

// LoginUser Method
func (api API) LoginUser(c *gin.Context) {
	var userProvider model.UserProvider
	// span, err := tracer.CreateTracerAndSpan("login", c)
	// if err != nil {
	// 	logger.Logger.Errorf(err.Error())
	// }
	err := c.ShouldBindJSON(&userProvider)
	if err != nil {
		fmt.Println(err.Error())
		//tracer.OnErrorLog(span, err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		//var acc model.UserAuth
		return
	}
	var (
		res = response.Gin{C: c}
	)
	u, isMail, isProvider := api.UserRepo.CheckEmail(userProvider)
	if isMail && isProvider {
		//span, err := tracer.CreateTracerAndSpan("check_email", c)
		// if err != nil {
		// 	fmt.Println(err.Error())
		// }
		// for _, s := range u.Provider {
		// 	if s.ProviderID == userProvider.ProviderID && s.ProviderName == userProvider.Provider {

		// 	}
		// }
		accessToken, refreshToken, err := oauth.GenerateTokenPair(u)
		if err != nil || accessToken == "" || refreshToken == "" {
			// Return if there is an error in creating the JWT return an internal server error
			//tracer.OnErrorLog(span, err)
			res.Response(http.StatusInternalServerError, "Could not generate token", nil)
			return
		}
		// span.LogFields(
		// 	tracelog.String("event", "success"),
		// 	tracelog.String("message", "returned token"),
		// 	tracelog.Int("status", http.StatusOK),
		// )
		res.Response(http.StatusOK, "login success", gin.H{"access_token": accessToken, "refresh_token": refreshToken})
		return

	} else if isMail {
		res.Response(http.StatusAlreadyReported, "Email in use, Are you want to update account with new provider?", nil)
		return
	}

	//span, err = tracer.CreateTracerAndSpan("add user", c)
	u2, err := api.UserRepo.AddUser(userProvider)
	if err != nil {
		//tracer.OnErrorLog(span, err)
		log.Println("error AddUserHandeler", err.Error())
		res.Response(http.StatusInternalServerError, "add user error", nil)
		return
	}
	accessToken, refreshToken, err := oauth.GenerateTokenPair(u2)
	if err != nil || accessToken == "" || refreshToken == "" {
		// Return if there is an error in creating the JWT return an internal server error
		//tracer.OnErrorLog(span, err)
		res.Response(http.StatusInternalServerError, "Could not generate token", nil)
		return
	}
	// span.LogFields(
	// 	tracelog.String("event", "success"),
	// 	tracelog.String("message", "returned token"),
	// 	tracelog.Int("status", http.StatusOK),
	// )
	res.Response(http.StatusOK, "login success", gin.H{"access_token": accessToken, "refresh_token": refreshToken})
	return
}

// UpdateUserProvider Method for update new provider id
func (api API) UpdateUserProvider(c *gin.Context) {
	var userProvider model.UserProvider
	// span, err := tracer.CreateTracerAndSpan("login", c)
	// if err != nil {
	// 	logger.Logger.Errorf(err.Error())
	// }
	err := c.ShouldBindJSON(&userProvider)
	var (
		res = response.Gin{C: c}
	)
	if err != nil {
		fmt.Println(err.Error())
		//tracer.OnErrorLog(span, err)
		res.Response(http.StatusBadRequest, err.Error(), nil)
		c.Abort()
		return
	}

	u, isHas, isProvider := api.UserRepo.CheckEmail(userProvider)
	if isHas && !isProvider {
		err = api.UserRepo.UpdateProvider(u, userProvider)
		if err != nil {
			res.Response(http.StatusInternalServerError, "Could not update provider", nil)
			c.Abort()
			return
		}
		accessToken, refreshToken, err := oauth.GenerateTokenPair(u)
		if err != nil || accessToken == "" || refreshToken == "" {
			// Return if there is an error in creating the JWT return an internal server error
			//tracer.OnErrorLog(span, err)
			res.Response(http.StatusInternalServerError, "Could not generate token", nil)
			c.Abort()
			return
		}
		// span.LogFields(
		// 	tracelog.String("event", "success"),
		// 	tracelog.String("message", "returned token"),
		// 	tracelog.Int("status", http.StatusOK),
		// )
		res.Response(http.StatusOK, "login success", gin.H{"access_token": accessToken, "refresh_token": refreshToken})
		return
	}
	res.Response(http.StatusInternalServerError, "Could not update provider", nil)
}

// LogoutUser Method
func (api API) LogoutUser(c *gin.Context) {

	span, err := tracer.CreateTracerAndSpan("logout", c)

	if err != nil {
		logger.Logger.Errorf(err.Error())
		//fmt.Println(err.Error())
	}

	token := c.GetHeader("Authorization")

	if token == "" {
		span.LogFields(
			tracelog.String("event", "error"),
			tracelog.String("message", "Authorization token was not provided"),
		)
		logger.Logger.Errorf("Authorization token was not provided")
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Authorization Token is required"})
		c.Abort()
		return
	}

	extractedToken := strings.Split(token, "Bearer ")

	err = oauth.InvalidateToken(extractedToken[1])
	if err != nil {
		tracer.OnErrorLog(span, err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": err.Error()})
		c.Abort()
		return
	}

	span.LogFields(
		tracelog.String("event", "success"),
		tracelog.Int("status", http.StatusAccepted),
	)
	c.JSON(http.StatusAccepted, gin.H{"status": http.StatusAccepted, "message": "Done"})

}
