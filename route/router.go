package route

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"thinkdev.app/think/runex/runexapi/api/v2/user"
	"thinkdev.app/think/runex/runexapi/middleware/oauth"
	"thinkdev.app/think/runex/runexapi/repository/v2"
)

//Router initialization
func Router(route *gin.Engine, connectionDB *mongo.Database) {
	userRepo := repository.RepoDB{
		ConnectionDB: connectionDB,
	}
	userAPI := user.API{
		UserRepo: userRepo,
	}
	api := route.Group("/api/v2")
	{
		api.POST("/login", userAPI.LoginUser)
		api.POST("/signup", userAPI.AddUser)
		api.POST("/refreshAccessToken", userAPI.RefreshAccessToken)
		api.POST("/verifyToken", userAPI.VerifyAuthToken)
		api.Use(oauth.AuthMiddleware())
		{
			api.GET("/user", userAPI.GetUser)
			api.PUT("/user", userAPI.UpdateUser)
			api.POST("/logout", userAPI.LogoutUser)
		}
	}
}
