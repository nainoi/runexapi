package route

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	handle_workouts "thinkdev.app/think/runex/runexapi/api/v1/workouts"
	handle_activity_v2 "thinkdev.app/think/runex/runexapi/api/v2/activity"
	"thinkdev.app/think/runex/runexapi/api/v2/board"
	"thinkdev.app/think/runex/runexapi/api/v2/config"
	eventV2 "thinkdev.app/think/runex/runexapi/api/v2/event"
	"thinkdev.app/think/runex/runexapi/api/v2/kao"
	"thinkdev.app/think/runex/runexapi/api/v2/migration"
	"thinkdev.app/think/runex/runexapi/api/v2/notification"
	"thinkdev.app/think/runex/runexapi/api/v2/paymethod"
	"thinkdev.app/think/runex/runexapi/api/v2/preorder"
	handle_register_v2 "thinkdev.app/think/runex/runexapi/api/v2/register"
	"thinkdev.app/think/runex/runexapi/api/v2/report"
	"thinkdev.app/think/runex/runexapi/api/v2/strava"
	"thinkdev.app/think/runex/runexapi/api/v2/tambon"
	"thinkdev.app/think/runex/runexapi/api/v2/team"
	"thinkdev.app/think/runex/runexapi/api/v2/upload"
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
		kaoGroup(*api, connectionDB)
		eventGroup(*api, connectionDB)
		registerGroup(*api, connectionDB)
		reportGroup(*api, connectionDB)
		tambonGroup(*api)
		payMethodGroup(*api)
		boardGroup(*api)
	}
}

func tambonGroup(g gin.RouterGroup) {
	g.GET("/tambon/:zipcode", tambon.SearchZipcode)
	g.GET("/province/:province", tambon.SearchProvince)
	g.GET("/amphoe/:amphoe", tambon.SearchAmphoe)
	g.GET("/district/:district", tambon.SearchDistrict)
	g.GET("tambons", tambon.All)
}

func payMethodGroup(g gin.RouterGroup) {
	g.GET("paymethods", paymethod.Get)
}

func userGroup(g gin.RouterGroup, connectionDB *mongo.Database) {
	userRepoDB := repository.RepoUserDB{
		ConnectionDB: connectionDB,
	}
	cfRepoDB := repository.ConfigRepositoryMongo{
		ConnectionDB: connectionDB,
	}
	userAPI := user.API{
		UserRepo: userRepoDB,
	}
	configAPI := config.ConfigAPI{
		ConfigRepo: cfRepoDB,
	}
	g.GET("/test", configAPI.GetTest)
	g.GET("/config", configAPI.GetConfig)
	g.POST("/login", userAPI.LoginUser)
	g.POST("/signup", userAPI.AddUser)
	g.POST("/refreshAccessToken", userAPI.RefreshAccessToken)
	g.POST("/verifyToken", userAPI.VerifyAuthToken)
	g.Use(oauth.AuthMiddleware())
	{
		g.GET("/user", userAPI.GetUser)
		g.PUT("/user", userAPI.UpdateUser)
		g.POST("/user", userAPI.Get)
		g.POST("/logout", userAPI.LogoutUser)
		g.GET("/logout", userAPI.Signout)
		g.PUT("/syncStrava", userAPI.UpdateUserStrava)
		g.POST("/registerFirebase", userAPI.RegFirebase)
		g.POST("/uploads", upload.Uploads)
		g.POST("/uploadCover", upload.UploadCover)
		g.POST("/teamIcon", team.UploadTeamIcon)
		g.POST("/getTeamIcon", team.GetTeamIcon)
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
	g.POST("/workout_hook", workoutsAPI.AddWorkoutHook)
	g.GET("/workoutUpdateLocation/:user_id", workoutsAPI.ManageWorkoutLocation)
	g.Use(oauth.AuthMiddleware())
	{
		g.POST("/workout", workoutsAPI.AddWorkout)
		g.POST("/workouts", workoutsAPI.AddMultiWorkout)
		g.DELETE("/workout", workoutsAPI.RemoveWorkoutActivity)
		g.GET("/workouts", workoutsAPI.GetWorkouts)
		g.POST("/workouts/history", workoutsAPI.GetWorkoutsHistoryMonth)
		g.GET("/workouts/historyAll", workoutsAPI.GetWorkoutsHistoryAll)
		g.GET("/workoutDetail/:id", workoutsAPI.GetWorkoutDetail)
	}

}

func kaoGroup(g gin.RouterGroup, connectionDB *mongo.Database) {
	kaoRepository := repo.KaoRepositoryMongo{
		ConnectionDB: connectionDB,
	}
	kaoAPI := kao.KoaAPI{
		KaoRepository: &kaoRepository,
	}
	g.Use(oauth.AuthMiddleware())
	{
		//g.POST("/kaoActivity", kaoAPI.)
		g.POST("/kao", kaoAPI.GetKaoActivity)
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
		group.POST("/waiting", activityV2API.GetActivityWaiting)
		group.POST("/activityWithState", activityV2API.GetActivityWithStatus)
		group.POST("/activityByUser", activityV2API.AdminGetActivityByUser)
		group.PUT("/update", activityV2API.UpdateActivity)
		group.POST("/adminDelete", activityV2API.AdminDeleteActivity)
		group.Use(oauth.AuthMiddleware())
		{
			group.POST("/add", activityV2API.AddActivity)
			group.DELETE("/remove", activityV2API.RemoveActivityEvent)
			group.GET("/getByEvent/:event", activityV2API.GetActivityByEvent)
			group.GET("/getByEvent2/:event", activityV2API.GetActivityByEvent2)
			group.POST("/dashboard", activityV2API.GetDashboard)
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
	// eventRepository := repository.EventRepositoryMongo{
	// 	ConnectionDB: connectionDB,
	// }
	// eventAPI := handle_event_v2.EventAPI{
	// 	EventRepository: &eventRepository,
	// }
	group := g.Group("/event")
	{
		// group.GET("/findByStatus/:status", eventAPI.GetByStatus)
		//group.GET("/eventInfo/:id", eventAPI.GetByID)
		// group.GET("/eventDetail/:slug", eventAPI.GetBySlug)
		group.GET("/all", eventV2.GetAll)
		// group.GET("/active", eventV2.GetAllActive)
		// group.GET("/getBySlug/:slug", eventAPI.GetBySlug)
		group.GET("/detail/:code", eventV2.GetDetail)
		group.Use(oauth.AuthMiddleware())
		{
			// group.POST("", eventAPI.AddEvent)
			// group.GET("/myEvent", eventAPI.MyEvent)
			// group.PUT("/edit/:id", eventAPI.EditEvent)
			// group.DELETE("/delete/:id", eventAPI.DeleteEvent)
			// group.POST("/:id/uploadImage", eventAPI.UploadImage)
			// group.POST("/:id/addProduct", eventAPI.AddProduct)
			// group.POST("/:id/editProduct", eventAPI.EditProduct)
			// group.GET("/getProduct/:id", eventAPI.GetProductEvent)
			// group.DELETE("/deleteProduct/:id/:productID", eventAPI.DeleteProductEvent)
			// group.POST("/:id/addTicket", eventAPI.AddTicket)
			// group.POST("/:id/editTicket", eventAPI.EditTicket)
			// group.DELETE("/deleteTicket/:id/:ticketID", eventAPI.DeleteTicketEvent)
			// group.PUT("/validateSlug", eventAPI.ValidateSlug)
		}
	}

}

func registerGroup(g gin.RouterGroup, connectionDB *mongo.Database) {
	registerRepository := repository.RegisterRepositoryMongo{
		ConnectionDB: connectionDB,
	}
	registerAPI := handle_register_v2.RegisterAPI{
		RegisterRepository: &registerRepository,
	}
	group := g.Group("/register")
	{
		group.POST("/payment_hook", registerAPI.PaymentHook)
		group.POST("/regsEvent", registerAPI.GetRegEventFromOwner)
		group.POST("/adminUpdatePayment", registerAPI.AdminSendSlip)
		group.Use(oauth.AuthMiddleware())
		{
			group.GET("/all", registerAPI.GetAll)
			group.POST("/add", registerAPI.AddRegister)
			group.POST("/team", registerAPI.AddTeamRegister)
			group.POST("/info", registerAPI.GetRegEventDashboard)
			group.POST("/addRace", registerAPI.AddRaceRegister)
			group.PUT("/edit/:id", registerAPI.EditRegister)
			group.GET("/myRegEvent", registerAPI.GetByUserID)
			group.GET("/checkUserRegisterEvent/:eventID", registerAPI.CheckUserRegisterEvent)
			group.GET("/checkRegEventCode/:code", registerAPI.CheckUserRegisterEventCode)
			group.POST("/checkRegER", registerAPI.CheckUserRegisterER)
			group.GET("/myRegEventActivate", registerAPI.GetMyRegEventActivate)
			group.GET("/regsEvent/:eventID", registerAPI.GetRegEventFromEventer)
			group.POST("/payment", registerAPI.ChargeRegEvent)
			group.POST("/regUpdateUserInfo", registerAPI.RegEventUpdateUserInfo)
			group.POST("/editUserInfo", registerAPI.EditUserInfo)
			group.POST("/sendSlip", registerAPI.SendSlip)
		}
	}

}

func reportGroup(g gin.RouterGroup, connectionDB *mongo.Database) {
	reportRepo := repo.ReportRepositoryMongo{
		ConnectionDB: connectionDB,
	}

	reportAPI := report.ReportAPI{
		ReportRepository: reportRepo,
	}

	group := g.Group("/report")
	{
		group.GET("/dashboard/:eventID", reportAPI.GetDashboardByEvent)
	}

}

//reportGroup ready board
func boardGroup(g gin.RouterGroup) {
	api := g.Group("/board")
	{
		api.POST("/rankings", board.GetAllBoardByEvent)
		api.POST("/finished", board.GetAllRankFinished)

		api.Use(oauth.AuthMiddleware())
		{
			api.POST("/ranking", board.GetBoardByEvent)
			api.POST("/update", board.GetBoardUpdateActivity)

		}
	}
}
