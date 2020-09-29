package route

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"thinkdev.app/think/runex/runexapi/api/v2/notification"
	"thinkdev.app/think/runex/runexapi/api/v2/preorder"
	"thinkdev.app/think/runex/runexapi/api/v2/strava"
	"thinkdev.app/think/runex/runexapi/api/v2/user"
	"thinkdev.app/think/runex/runexapi/middleware/oauth"
	repo "thinkdev.app/think/runex/runexapi/repository"
	"thinkdev.app/think/runex/runexapi/repository/v2"
)

//Router initialization
func Router(route *gin.Engine, connectionDB *mongo.Database) {

	api := route.Group("/api/v2")
	{
		userGroup(*api, connectionDB)
		preOrderGroup(*api, connectionDB)
		syncGroup(*api, connectionDB)
		notiGroup(*api, connectionDB)
	}
}

func userGroup(g gin.RouterGroup, connectionDB *mongo.Database) {
	userRepo := repository.RepoDB{
		ConnectionDB: connectionDB,
	}
	userAPI := user.API{
		UserRepo: userRepo,
	}
	g.POST("/login", userAPI.LoginUser)
	g.POST("/signup", userAPI.AddUser)
	g.POST("/refreshAccessToken", userAPI.RefreshAccessToken)
	g.POST("/verifyToken", userAPI.VerifyAuthToken)
	g.Use(oauth.AuthMiddleware())
	{
		g.GET("/user", userAPI.GetUser)
		g.PUT("/user", userAPI.UpdateUser)
		g.POST("/logout", userAPI.LogoutUser)
		g.PUT("/syncStrava", userAPI.UpdateUserStrava)
	}
}

func preOrderGroup(g gin.RouterGroup, connectionDB *mongo.Database) {
	preRepo := repo.PreorderRepositoryMongo{
		ConnectionDB: connectionDB,
	}

	preAPI := preorder.API{
		PreRepo: preRepo,
	}
	g.POST("/searchPreOrder", preAPI.SearchPreOrder)
}

func syncGroup(g gin.RouterGroup, connectionDB *mongo.Database) {
	repoStrava := repo.RepoDB{
		ConnectionDB: connectionDB,
	}

	stravaAPI := strava.API{
		Repo: repoStrava,
	}
	group := g.Group("/strava")
	{
		group.POST("/activity", stravaAPI.AddStravaActivity)
		group.GET("/activities", stravaAPI.GetStravaActivities)
	}

}

func notiGroup(g gin.RouterGroup, connectionDB *mongo.Database) {
	// repoStrava := repo.RepoDB{
	// 	ConnectionDB: connectionDB,
	// }

	// stravaAPI := strava.API{
	// 	Repo: repoStrava,
	// }
	g.POST("/notificationOne", notification.SendOneNotification)

}