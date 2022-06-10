package route

import (
	"log"
	"net/http"
	"time"

	handle_activity_v2 "thinkdev.app/think/runex/runexapi/api/v2/activity"

	// handle_banner "thinkdev.app/think/runex/runexapi/api/v1/banner"
	handle_runHistory "thinkdev.app/think/runex/runexapi/api/v1/runHistory"

	//handle_importData "thinkdev.app/think/runex/runexapi/api/v1/importdata"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"thinkdev.app/think/runex/runexapi/api/v1/uploads"
	handle_user "thinkdev.app/think/runex/runexapi/api/v1/user"
	auth "thinkdev.app/think/runex/runexapi/middleware"
	"thinkdev.app/think/runex/runexapi/repository"
)

// type Routes struct {
// }

// Route for user and authentication
func Route(route *gin.Engine, connectionDB *mongo.Database) {
	userRepository := repository.UserRepositoryMongo{
		ConnectionDB: connectionDB,
	}
	userAPI := handle_user.UserAPI{
		UserRepository: &userRepository,
	}
	middleware := auth.Auth(userAPI)

	//route.StaticFS("/uploads/images/", http.Dir(upload.GetImagePath()))
	//route.Static("/upload", "./upload")
	api := route.Group("/api/v1")
	{

		api.POST("/user/ep", userAPI.AddEP)
		api.POST("/user/pd", userAPI.AddPD)
		api.POST("/user/login", middleware.LoginHandler)
		api.POST("/user/loginPD", middleware.LoginHandler)
		api.POST("/user/forgotpass", userAPI.ForgotPassword)
		api.POST("/user/forgotpassword", userAPI.ForgotPasswordMobile)
		api.POST("/user/updatepass", userAPI.UpdatePassword)
		api.Use(middleware.MiddlewareFunc())
		{
			api.POST("/uploads", uploads.Uploads)
			api.POST("/uploadCover", uploads.UploadCover)
			api.POST("/user/avatar", userAPI.UpdateAvatar)
			api.POST("/user/address", userAPI.AddAdress)
			api.POST("/uploadSlip", uploads.UploadSlip)
			api.POST("/uploadWithFolder", uploads.UploadWithFolder)
			api.GET("/user", userAPI.Get)
			api.PUT("/user", userAPI.Edit)
			api.GET("/user/confirm", userAPI.Confirm)
			api.POST("/user/changepass", userAPI.ChangePassword)
			api.DELETE("/user/:id", userAPI.Delete)
			api.GET("/user/logout", func(c *gin.Context) {
				log.Println("logout")
				if token, err := middleware.CheckIfTokenExpire(c); err == nil {
					if err2 := token.Valid(); err2 == nil {
						log.Println("valid")
						middleware.DisabledAbort = true
						token["exp"] = time.Now().UTC().Unix()
						c.Abort()
						res := gin.H{"msg": "success"}
						c.JSON(http.StatusOK, res)
					}
				} else {
					log.Println(err)
				}

			})
		}
	}

	RunHistoryRoute(route, connectionDB, middleware)
}

//RunHistoryRoute history
func RunHistoryRoute(route *gin.Engine, connectionDB *mongo.Database, middleware *jwt.GinJWTMiddleware) {
	runHistoryRepository := repository.RunHistoryRepositoryMongo{
		ConnectionDB: connectionDB,
	}
	runHistoryAPI := handle_runHistory.RunHistoryAPI{
		RunHistoryRepository: &runHistoryRepository,
	}
	api := route.Group("/api/v1/runhistory")
	{

		api.Use(middleware.MiddlewareFunc())
		{
			api.GET("/myhistory", runHistoryAPI.MyRunHistory)
			api.POST("/add", runHistoryAPI.AddRunHistory)
			api.DELETE("/deleteActivity/:activityID", runHistoryAPI.DeleteActivityHistory)

		}
	}
}

func ActivityV2Route(route *gin.Engine, connectionDB *mongo.Database, middleware *jwt.GinJWTMiddleware) {
	activityV2Repository := repository.ActivityV2RepositoryMongo{
		ConnectionDB: connectionDB,
	}
	activityV2API := handle_activity_v2.ActivityV2API{
		ActivityV2Repository: &activityV2Repository,
	}
	api := route.Group("/api/v2/activity")
	{
		api.Use(middleware.MiddlewareFunc())
		{
			api.POST("/add", activityV2API.AddActivity)
			api.GET("/getByEvent/:event", activityV2API.GetActivityByEvent)
			api.GET("/getByEvent2/:event", activityV2API.GetActivityByEvent2)
			api.POST("/getHistoryDay", activityV2API.GetHistoryDayByEvent)
			api.POST("/getHistoryMonth", activityV2API.GetHistoryMonthByEvent)
			api.DELETE("/deleteActivity/:id/:activityID", activityV2API.DeleteActivityEvent)
		}
	}
}
