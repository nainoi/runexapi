package user

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	guuid "github.com/google/uuid"

	config "thinkdev.app/think/runex/runexapi/config"
	//"github.com/appleboy/gin-jwt/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin/binding"

	"github.com/gin-gonic/gin"
	"thinkdev.app/think/runex/runexapi/api/mail"
	"thinkdev.app/think/runex/runexapi/model"
	"thinkdev.app/think/runex/runexapi/pkg/app"
	"thinkdev.app/think/runex/runexapi/pkg/e"
	"thinkdev.app/think/runex/runexapi/repository"
	"thinkdev.app/think/runex/runexapi/utils"
	//jwt "github.com/dgrijalva/jwt-go"
	//"golang.org/x/crypto/bcrypt"
)

// UserAPI is a representation of a UserAPI
type UserAPI struct {
	UserRepository repository.UserRepository
}

// Login is a api of a signin from jwt
func (api UserAPI) Login(context *gin.Context) (model.UserAuth, error) {
	log.Println("Login EP")
	var user model.UserAuth
	var err error

	var userEmail model.LoginEmail
	err3 := context.ShouldBindBodyWith(&userEmail, binding.JSON)
	log.Println(userEmail)
	if err3 != nil {
		log.Println("error param email", err3.Error())
		//context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return user, err3
	}
	user, err = api.UserRepository.GetUserByEmail(userEmail.Email, userEmail.Password)
	user.PF = userEmail.PF
	return user, err
}

//LoginPD is api of signin with provider
func (api UserAPI) LoginPD(context *gin.Context) (model.UserAuth, error) {
	log.Println("Login PD")
	var user model.UserAuth
	var err error

	var userPD model.LoginProvider
	err3 := context.ShouldBindBodyWith(&userPD, binding.JSON)
	if err3 != nil {
		log.Println("error param provider and email", err3.Error())
		//context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return user, err3
	}
	user, err = api.UserRepository.GetUserByProvider(userPD.Provider, userPD.ProviderID)
	user.PF = userPD.PF
	return user, err
}

// AddEP api for add user from email and password
func (api UserAPI) AddEP(context *gin.Context) /*(model.UserAuth, error)*/ {
	//log.Println(context.Params)
	var user model.UserMail
	err := context.ShouldBindBodyWith(&user, binding.JSON)
	if err != nil {
		log.Println("error bind user email", err.Error())
		context.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		//var acc model.UserAuth
		return
	}

	if api.UserRepository.CheckEmail(user.Email) {
		log.Println("email in use")
		context.JSON(http.StatusIMUsed, gin.H{"msg": "Email in use"})
		return
	}
	if user.FullName == "" {
		user.FullName = user.FirstName + " " + user.LastName
	}
	user.Address = []model.Address{}
	user.Password = utils.HashAndSalt([]byte(user.Password))
	err2 := api.UserRepository.AddUserEP(user)
	if err2 != nil {
		log.Println("error AddUserHandeler", err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	mail.SendConfirmRegMail(user.FirstName+user.LastName, user.Email, user.UserID, user.PF, user.Role)
	context.Redirect(http.StatusTemporaryRedirect, "/api/v1/user/login")

}

// AddAdmin api for add user admin
func (api UserAPI) AddAdmin(context *gin.Context) /*(model.UserAuth, error)*/ {
	//log.Println(context.Params)
	var user model.UserMail
	err := context.ShouldBindBodyWith(&user, binding.JSON)
	if err != nil {
		log.Println("error AddUserHandeler", err.Error())
		context.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		//var acc model.UserAuth
		return
	}

	if api.UserRepository.CheckEmail(user.Email) {
		log.Println("email in use")
		context.JSON(http.StatusIMUsed, gin.H{"msg": "Email in use"})
		return
	}
	user.Password = utils.HashAndSalt([]byte(user.Password))
	err2 := api.UserRepository.AddUserAdmin(user)
	if err2 != nil {
		log.Println("error AddUserHandeler", err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	context.Redirect(http.StatusTemporaryRedirect, "/api/v1/user/login")
}

// AddPD api for add user from social provider
func (api UserAPI) AddPD(context *gin.Context) /*(model.UserAuth, error) */ {
	var user model.UserProvider
	err := context.ShouldBindJSON(&user)
	if err != nil {
		log.Println("error AddUserHandeler", err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		//var acc model.UserAuth
		return
	}
	if api.UserRepository.CheckEmail(user.Email) {
		if api.UserRepository.CheckProvider(user) {
			context.Redirect(http.StatusTemporaryRedirect, "/api/v1/user/loginPD")
			return
		}
		err = api.UserRepository.UpdateProvider(user)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}
		context.Redirect(http.StatusTemporaryRedirect, "/api/v1/user/loginPD")
		return
	}
	user.Address = []model.Address{}
	err2 := api.UserRepository.AddUserPD(user)
	if err2 != nil {
		log.Println("error AddUserHandeler", err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	context.Redirect(http.StatusTemporaryRedirect, "/api/v1/user/loginPD")
}

// UserList api for list all user
func (api UserAPI) UserList(context *gin.Context) {
	var usersInfo model.Users
	users, err := api.UserRepository.GetAllUser()
	if err != nil {
		log.Println("error userListHandler", err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	usersInfo.User = users
	context.JSON(http.StatusOK, usersInfo)
}

// Get api for token
func (api UserAPI) Get(context *gin.Context) {
	//claims := jwt.ExtractClaims(context)
	userID, _, _ := utils.GetTokenValue(context)
	if userID != "" {
		user, err := api.UserRepository.GetUser(userID)
		if err != nil {
			log.Println("error GetUserByIDHandler", err.Error())
			context.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}
		res := gin.H{"msg": "success", "data": user}
		context.JSON(http.StatusOK, res)
	}
	// if userID, ok := claims[config.ID_KEY].(string); ok {

	// }
}

func (api UserAPI) Confirm(context *gin.Context) {
	userID, _, _ := utils.GetTokenValue(context)
	if userID != "" {
		err := api.UserRepository.Confirm(userID)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}

		context.JSON(http.StatusOK, gin.H{"msg": "success"})
		return
	}
}

func (api UserAPI) UpdateAvatar(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	id, _, _ := utils.GetTokenValue(c)
	log.Printf("[info] id %s", id)

	file, header, err := c.Request.FormFile("upload")
	filename := header.Filename
	fmt.Println(filename)

	uniqidFilename := guuid.New()
	fmt.Printf("github.com/google/uuid:         %s\n", uniqidFilename.String())

	pathDir := "." + config.UPLOAD_AVATAR
	if _, err := os.Stat(pathDir); os.IsNotExist(err) {
		os.MkdirAll(pathDir, os.ModePerm)
	}

	out, err := os.Create(pathDir + uniqidFilename.String() + ".png")

	path := config.UPLOAD_AVATAR + uniqidFilename.String() + ".png"
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
		//log.Fatal(err)
	}

	err2 := api.UserRepository.UploadAvatar(id, path)
	if err2 != nil {
		log.Println("error UploadImage", err2.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err2.Error()})
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, gin.H{"url": path})
}

// Edit api for update user
func (api UserAPI) Edit(context *gin.Context) {
	var user model.User
	var (
		appG = app.Gin{C: context}
	)
	id, _, _ := utils.GetTokenValue(context)
	err := context.ShouldBindJSON(&user)
	if err != nil {
		log.Println("error bind user", err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	users, err2 := api.UserRepository.EditUser(id, user)
	if err2 != nil {
		log.Println("error update user", err2.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, users)
}

// Delete api for delete user
func (api UserAPI) Delete(context *gin.Context) {
	var (
		appG = app.Gin{C: context}
	)
	userID := context.Param("user_id")
	err := api.UserRepository.DeleteUserByID(userID)
	if err != nil {
		log.Println("error DeleteUserHandler", err.Error())
		appG.Response(http.StatusInternalServerError, e.ERROR, gin.H{"message": err.Error()})
	}
	appG.Response(http.StatusOK, e.SUCCESS, gin.H{"message": e.SUCCESS})
}

//AddAdress to new address
func (api UserAPI) AddAdress(context *gin.Context) {
	var (
		appG = app.Gin{C: context}
	)
	userID, _, _ := utils.GetTokenValue(context)
	var address model.Address
	err := context.ShouldBindBodyWith(&address, binding.JSON)
	log.Println(address)
	if err != nil {
		log.Println("error binding address", err.Error())
		appG.Response(http.StatusInternalServerError, e.ERROR, gin.H{"message": err.Error()})
	}
	user, err2 := api.UserRepository.AddAddress(userID, address)
	if err2 != nil {
		log.Println("error get user", err.Error())
		appG.Response(http.StatusInternalServerError, e.ERROR, gin.H{"message": err.Error()})
	}
	appG.Response(http.StatusOK, e.SUCCESS, user)
}

//UpdateAdress user
func (api UserAPI) UpdateAdress(context *gin.Context) {
	var (
		appG = app.Gin{C: context}
	)
	userID, _, _ := utils.GetTokenValue(context)
	var address model.Address
	err := context.ShouldBindBodyWith(&address, binding.JSON)
	if err != nil {
		log.Println("error binding address", err.Error())
		appG.Response(http.StatusInternalServerError, e.ERROR, gin.H{"message": err.Error()})
	}
	user, err2 := api.UserRepository.UpdateAddress(userID, address)
	if err2 != nil {
		log.Println("error get user", err.Error())
		appG.Response(http.StatusInternalServerError, e.ERROR, gin.H{"message": err.Error()})
	}
	appG.Response(http.StatusOK, e.SUCCESS, user)
}

//ChangePassword user
func (api UserAPI) ChangePassword(context *gin.Context) {
	var (
		appG = app.Gin{C: context}
	)
	userID, _, _ := utils.GetTokenValue(context)
	var user model.UserAuth
	err := context.ShouldBindBodyWith(&user, binding.JSON)
	if err != nil {
		log.Println("error binding user %s", err.Error())
		appG.Response(http.StatusInternalServerError, e.ERROR, gin.H{"message": err.Error()})
		return
	}
	err2 := api.UserRepository.ChangePassword(userID, user.Password, user.NewPassword)
	if err2 != nil {
		log.Println("error get user %s", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, gin.H{"message": "old password not match"})
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, gin.H{"message": e.SUCCESS})
}

//UpdatePassword user
func (api UserAPI) UpdatePassword(context *gin.Context) {
	var (
		appG = app.Gin{C: context}
	)
	tokenString := context.GetHeader("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(config.SECRET_KEY), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims[config.ID_KEY])
		userID := claims[config.ID_KEY].(string)
		var user model.UserAuth
		err = context.ShouldBindBodyWith(&user, binding.JSON)
		if err != nil {
			log.Println("error binding user %s", err.Error())
			appG.Response(http.StatusInternalServerError, e.ERROR, gin.H{"message": err.Error()})
			return
		}
		err = api.UserRepository.UpdatePassword(userID, user.NewPassword)
		if err != nil {
			log.Println("error get user %s", err)
			appG.Response(http.StatusInternalServerError, e.ERROR, gin.H{"message": "update not success"})
			return
		}
		appG.Response(http.StatusOK, e.SUCCESS, gin.H{"message": e.SUCCESS})
	} else {
		fmt.Println(err)
		appG.Response(http.StatusForbidden, e.ERROR, gin.H{"message": "update not success"})
		return
	}

	// claims := jwt.ExtractClaims(context)
	// log.Println(claims)
	// userID, _, _ := utils.GetTokenValue(context)
	// var user model.UserAuth
	// err = context.ShouldBindBodyWith(&user, binding.JSON)
	// if err != nil {
	// 	log.Println("error binding user %s", err.Error())
	// 	appG.Response(http.StatusInternalServerError, e.ERROR, gin.H{"message": err.Error()})
	// 	return
	// }
	// err2 := api.UserRepository.UpdatePassword(userID, user.NewPassword)
	// if err2 != nil {
	// 	log.Println("error get user %s", err)
	// 	appG.Response(http.StatusInternalServerError, e.ERROR, gin.H{"message": "update not success"})
	// 	return
	// }
	// appG.Response(http.StatusOK, e.SUCCESS, gin.H{"message": e.SUCCESS})
}

//ForgotPassword user
func (api UserAPI) ForgotPassword(context *gin.Context) {
	var (
		appG = app.Gin{C: context}
	)
	var user model.UserForgot
	err := context.ShouldBindBodyWith(&user, binding.JSON)
	if err != nil {
		log.Println("error binding user %s", err.Error())
		appG.Response(http.StatusBadRequest, e.ERROR, gin.H{"message": err.Error()})
		return
	}
	acc, err2 := api.UserRepository.ForgotPassword(user.Email)
	if err2 != nil {
		log.Println("error get user %s", err)
		appG.Response(http.StatusBadRequest, e.ERROR, gin.H{"message": "This email does not exist in the system."})
		return
	}
	err = mail.SendForgotMail(acc)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, gin.H{"message": "Can not send email."})
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, gin.H{"message": e.SUCCESS})
}

//ForgotPassword user
func (api UserAPI) ForgotPasswordMobile(context *gin.Context) {
	var (
		appG = app.Gin{C: context}
	)
	var user model.UserForgot
	err := context.ShouldBindBodyWith(&user, binding.JSON)
	if err != nil {
		log.Println("error binding user %s", err.Error())
		appG.Response(http.StatusBadRequest, e.ERROR, gin.H{"message": err.Error()})
		return
	}
	acc, err2 := api.UserRepository.ForgotPassword(user.Email)
	if err2 != nil {
		log.Println("error get user %s", err)
		appG.Response(http.StatusBadRequest, e.ERROR, gin.H{"message": "This email does not exist in the system."})
		return
	}
	token, err := mail.GenarateToken(acc.UserID, acc.Role, acc.PF)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, gin.H{"message": "Can not send email."})
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, gin.H{"message": e.SUCCESS, "data": token})
}
