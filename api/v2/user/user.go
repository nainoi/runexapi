package user

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"thinkdev.app/think/runex/runexapi/api/v2/response"

	//"thinkdev.app/think/runex/runexapi/api/v2/tracer"
	"thinkdev.app/think/runex/runexapi/logger"
	"thinkdev.app/think/runex/runexapi/middleware/oauth"
	"thinkdev.app/think/runex/runexapi/model"
	"thinkdev.app/think/runex/runexapi/repository/v2"
	//stdopentracing "github.com/opentracing/opentracing-go"
	//tracelog "github.com/opentracing/opentracing-go/log"
)

//API struct for user repository
type API struct {
	UserRepo repository.RepoUserDB
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
// @Summary Refresh a JSON Web Token
// @Description Authenticates a user and provides a JWT to refresh Authorize API calls
// @ID Authentication
// @Tags user
// @Consume application/x-www-form-urlencoded
// @Produce json
// @Param refresh_token formData string true "RefreshToken"
// @Success 200 {object} response.ResponseOAuth
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /refreshAccessToken [post]
func (api API) RefreshAccessToken(c *gin.Context) {

	var tokenRequest oauth.RefreshTokenRequestBody
	var (
		res = response.Gin{C: c}
	)

	log.Println(tokenRequest)

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

		newToken, newRefreshToken, err := oauth.GenerateTokenPair(u)
		if err != nil {
			logger.Logger.Errorf(err.Error())
			res.Response(http.StatusBadRequest, "Cannot Generate New Access Token", nil)
			c.Abort()
			return
		}
		var (
			responseJWT = response.ResponseOAuth{
				AccessToken:  newToken,
				RefreshToken: newRefreshToken,
			}
		)
		res.Response(http.StatusOK, "success", responseJWT)
		c.Abort()
		return
	}

	res.Response(http.StatusBadRequest, "Error Found ", nil)
}

// GetUser godoc
// @Summary Get user info
// @Description get user info API calls
// @Security bearerAuth
// @Tags user
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=model.User}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /user [get]
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

// AddUser godoc
// @Summary Add new user
// @Description add new user API calls
// @Consume application/x-www-form-urlencoded
// @Accept  application/json
// @Produce application/json
// @Tags user
// @Param payload body model.UserProviderRequest true "payload"
// @Success 200 {object} response.ResponseOAuth
// @Success 208 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /signup [post]
func (api API) AddUser(c *gin.Context) {
	var userProvider model.UserProviderRequest
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

// UpdateUser update user info
// @Summary Update user info
// @Description put user info to update user info API calls
// @Security bearerAuth
// @Tags user
// @Accept  application/json
// @Produce application/json
// @Param payload body model.User true "payload"
// @Success 200 {object} response.Response{data=model.User}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /user [put]
func (api API) UpdateUser(c *gin.Context) {
	var user model.User
	err := c.ShouldBindJSON(&user)
	var (
		res = response.Gin{C: c}
	)
	userID, _ := oauth.GetValuesToken(c)
	if err != nil {
		fmt.Println(err.Error())
		//tracer.OnErrorLog(span, err)
		res.Response(http.StatusBadRequest, err.Error(), nil)
		//var acc model.UserAuth
		return
	}
	user, err = api.UserRepo.UpdateUser(user, userID)
	if err != nil {
		fmt.Println(err.Error())
		//tracer.OnErrorLog(span, err)
		res.Response(http.StatusInternalServerError, err.Error(), nil)
		//var acc model.UserAuth
		return
	}
	res.Response(http.StatusOK, "update success", user)
}

// LoginUser godoc
// @Summary user login
// @Description login or add new user from open id API calls
// @Consume application/x-www-form-urlencoded
// @Tags user
// @Accept  application/json
// @Produce application/json
// @Param payload body model.UserProviderRequest true "payload"
// @Success 200 {object} response.ResponseOAuth
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /login [post]
func (api API) LoginUser(c *gin.Context) {
	var userProvider model.UserProviderRequest
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
		//var acc model.UserAuth
		return
	}

	u, err := api.UserRepo.Signin(userProvider)
	if err == nil {
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

	}

	// Return if there is an error in creating the JWT return an internal server error
	//tracer.OnErrorLog(span, err)
	res.Response(http.StatusInternalServerError, "Could not generate token", nil)
	return
}

// UpdateUserStrava Method for update new provider id
// @Summary Update User sync user with strava
// @Description Update User sync user with strava API calls
// @Consume application/x-www-form-urlencoded
// @Tags user
// @Security bearerAuth
// @Accept  application/json
// @Produce application/json
// @Param payload body model.UserStravaSyncRequest true "payload"
// @Success 200 {object} response.Response{data=model.User}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /syncStrava [put]
func (api API) UpdateUserStrava(c *gin.Context) {
	var userStrava model.UserStravaSyncRequest
	// span, err := tracer.CreateTracerAndSpan("login", c)
	// if err != nil {
	// 	logger.Logger.Errorf(err.Error())
	// }
	userID, _ := oauth.GetValuesToken(c)
	var (
		res = response.Gin{C: c}
	)
	if userID != "" {
		err := c.ShouldBindJSON(&userStrava)

		if err != nil {
			fmt.Println(err.Error())
			//tracer.OnErrorLog(span, err)
			res.Response(http.StatusBadRequest, err.Error(), nil)
			c.Abort()
			return
		}

		u, err := api.UserRepo.UpdateUserStrava(userStrava, userID)
		if err != nil {
			res.Response(http.StatusInternalServerError, "Could not update user sync with strava", nil)
			c.Abort()
			return
		}

		res.Response(http.StatusOK, "sync strava success", u)
		c.Abort()
		return
	}

	res.Response(http.StatusUnauthorized, "Unauthorized", nil)
	c.Abort()
	return
}

/*
// UpdateUserProvider Method for update new provider id
// @Summary Update User Provider
// @Description Update User Provider info from login API calls
// @Consume application/x-www-form-urlencoded
// @Tags user
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} response.ResponseOAuth
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /login [put]
func (api API) UpdateUserProvider(c *gin.Context) {
	var userProvider model.UserProviderRequest
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
*/

// LogoutUser Method
// @Summary user logout
// @Description user logout system API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags user
// @Accept  application/json
// @Produce application/json
// @Success 202 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /logout [post]
func (api API) LogoutUser(c *gin.Context) {

	//span, err := tracer.CreateTracerAndSpan("logout", c)

	// if err != nil {
	// 	//logger.Logger.Errorf(err.Error())
	// 	//fmt.Println(err.Error())
	// }
	var (
		res = response.Gin{C: c}
	)

	token := c.GetHeader("Authorization")
	userID, _ := oauth.GetValuesToken(c)
	var firebaseToken model.RegisterTokenRequest
	err := c.BindJSON(&firebaseToken)
	if token == "" {
		// span.LogFields(
		// 	tracelog.String("event", "error"),
		// 	tracelog.String("message", "Authorization token was not provided"),
		// )
		// logger.Logger.Errorf("Authorization token was not provided")
		res.Response(http.StatusUnauthorized, "Authorization Token is required", nil)
		c.Abort()
		return
	}

	extractedToken := strings.Split(token, "Bearer ")

	err = oauth.InvalidateToken(extractedToken[1])
	if err != nil {
		//tracer.OnErrorLog(span, err)
		res.Response(http.StatusInternalServerError, err.Error(), nil)
		c.Abort()
		return
	}

	// span.LogFields(
	// 	tracelog.String("event", "success"),
	// 	tracelog.Int("status", http.StatusAccepted),
	// )
	if firebaseToken.FirebaseToken != "" {
		api.UserRepo.FirebaseRemove(firebaseToken, userID)
	}

	res.Response(http.StatusAccepted, "Done", nil)
	return

}

// RegFirebase godoc
// @Summary user firebase register
// @Description user firebase register API calls
// @Consume application/x-www-form-urlencoded
// @Tags user
// @Accept  application/json
// @Produce application/json
// @Param payload body model.RegisterTokenRequest true "payload"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /registerFirebase [post]
func (api API) RegFirebase(c *gin.Context) {
	var firebaseRequest model.RegisterTokenRequest
	// span, err := tracer.CreateTracerAndSpan("login", c)
	// if err != nil {
	// 	logger.Logger.Errorf(err.Error())
	// }
	err := c.ShouldBindJSON(&firebaseRequest)
	var (
		res = response.Gin{C: c}
	)
	if err != nil {
		fmt.Println(err.Error())
		//tracer.OnErrorLog(span, err)
		res.Response(http.StatusBadRequest, err.Error(), nil)
		//var acc model.UserAuth
		return
	}

	userID, _ := oauth.GetValuesToken(c)

	err = api.UserRepo.FirebaseRegister(firebaseRequest.FirebaseToken, userID)
	if err == nil {
		//span, err := tracer.CreateTracerAndSpan("check_email", c)
		// if err != nil {
		// 	fmt.Println(err.Error())
		// }
		// for _, s := range u.Provider {
		// 	if s.ProviderID == userProvider.ProviderID && s.ProviderName == userProvider.Provider {

		// 	}
		// }
		res.Response(http.StatusOK, "firebase register success", nil)
		c.Abort()
		return

	}

	// Return if there is an error in creating the JWT return an internal server error
	//tracer.OnErrorLog(span, err)
	res.Response(http.StatusInternalServerError, "Could not register firebase", nil)
	return
}
