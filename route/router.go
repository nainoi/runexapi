package route

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	handle_workouts "thinkdev.app/think/runex/runexapi/api/v1/workouts"
	handle_activity_v2 "thinkdev.app/think/runex/runexapi/api/v2/activity"
	handle_event_v2 "thinkdev.app/think/runex/runexapi/api/v2/event"
	"thinkdev.app/think/runex/runexapi/api/v2/migration"
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
		migrationGroup(*api, connectionDB)
		workoutGroup(*api, connectionDB)
		activityGroup(*api, connectionDB)
		eventGroup(*api, connectionDB)
	}
}

func userGroup(g gin.RouterGroup, connectionDB *mongo.Database) {
	userRepoDB := repository.RepoUserDB{
		ConnectionDB: connectionDB,
	}
	userAPI := user.API{
		UserRepo: userRepoDB,
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
		g.POST("/registerFirebase", userAPI.RegFirebase)
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
	repoStrava := repo.RepoStravaDB{
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

func workoutGroup(g gin.RouterGroup, connectionDB *mongo.Database) {
	workoutsRepository := repo.WorkoutsRepositoryMongo{
		ConnectionDB: connectionDB,
	}
	workoutsAPI := handle_workouts.WorkoutsAPI{
		WorkoutsRepository: &workoutsRepository,
	}
	g.Use(oauth.AuthMiddleware())
	{
		g.POST("/workout", workoutsAPI.AddWorkout)
		g.GET("/workouts", workoutsAPI.GetWorkouts)
	}

}

func activityGroup(g gin.RouterGroup, connectionDB *mongo.Database) {
	activityV2Repository := repo.ActivityV2RepositoryMongo{
		ConnectionDB: connectionDB,
	}
	activityV2API := handle_activity_v2.ActivityV2API{
		ActivityV2Repository: &activityV2Repository,
	}
	group := g.Group("/activity")
	{
		group.Use(oauth.AuthMiddleware())
		{
			group.POST("/add", activityV2API.AddActivity)
			group.GET("/getByEvent/:event", activityV2API.GetActivityByEvent)
			group.GET("/getByEvent2/:event", activityV2API.GetActivityByEvent2)
			group.POST("/getHistoryDay", activityV2API.GetHistoryDayByEvent)
			group.POST("/getHistoryMonth", activityV2API.GetHistoryMonthByEvent)
			group.DELETE("/deleteActivity/:id/:activityID", activityV2API.DeleteActivityEvent)
			group.POST("/activityWorkout", activityV2API.AddFromWorkout)
			group.POST("/activitiesWorkout", activityV2API.AddMultipleFromWorkout)
		}
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

func migrationGroup(g gin.RouterGroup, connectionDB *mongo.Database) {
	migrationRepo := repo.MigrationRepositoryMongo{
		ConnectionDB: connectionDB,
	}

	migrationAPI := migration.MigrationAPI{
		MigrationRepository: migrationRepo,
	}
	g.POST("/migrateWorkout/:newCollection", migrationAPI.MigrateWorkout)
}

func eventGroup(g gin.RouterGroup, connectionDB *mongo.Database) {
	eventRepository := repo.EventRepositoryMongo{
		ConnectionDB: connectionDB,
	}
	eventAPI := handle_event_v2.EventAPI{
		EventRepository: &eventRepository,
	}
	group := g.Group("/event")
	{
		group.Use(oauth.AuthMiddleware())
		{
			group.POST("", eventAPI.AddEvent)
			group.GET("/myEvent", eventAPI.MyEvent)
			group.POST("", eventAPI.AddEvent)
			group.PUT("/edit/:id", eventAPI.EditEvent)
			group.DELETE("/delete/:id", eventAPI.DeleteEvent)
			group.POST("/:id/uploadImage", eventAPI.UploadImage)
			group.POST("/:id/addProduct", eventAPI.AddProduct)
			group.POST("/:id/editProduct", eventAPI.EditProduct)
			group.GET("/getProduct/:id", eventAPI.GetProductEvent)

			group.DELETE("/deleteProduct/:id/:productID", eventAPI.DeleteProductEvent)
			group.POST("/:id/addTicket", eventAPI.AddTicket)
			group.POST("/:id/editTicket", eventAPI.EditTicket)
			group.DELETE("/deleteTicket/:id/:ticketID", eventAPI.DeleteTicketEvent)
			group.PUT("/validateSlug", eventAPI.ValidateSlug)
		}
	}

}
