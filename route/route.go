package route

import (
	_ "fmt"
	"log"
	"net/http"
	"time"

	// handle_activity "bitbucket.org/suthisakch/runex/api/v1/activity"
	// handle_admin "bitbucket.org/suthisakch/runex/api/v1/admin"
	handle_banner "thinkdev.app/think/runex/runexapi/api/v1/banner"
	"thinkdev.app/think/runex/runexapi/api/v1/board"

	//handle_importData "bitbucket.org/suthisakch/runex/api/v1/importdata"
	// handle_register "thinkdev.app/think/runex/runexapi/api/v1/register"
	// "thinkdev.app/think/runex/runexapi/api/v1/uploads"
	// handle_user "thinkdev.app/think/runex/runexapi/api/v1/user"
	// auth "thinkdev.app/think/runex/runexapi/middleware"
	// "thinkdev.app/think/runex/runexapi/repository"

	handle_activity "thinkdev.app/think/runex/runexapi/api/v1/activity"
	handle_activity_v2 "thinkdev.app/think/runex/runexapi/api/v2/activity"

	handle_admin "thinkdev.app/think/runex/runexapi/api/v1/admin"

	// handle_banner "thinkdev.app/think/runex/runexapi/api/v1/banner"
	handle_category "thinkdev.app/think/runex/runexapi/api/v1/category"
	handle_coupon "thinkdev.app/think/runex/runexapi/api/v1/coupon"
	handle_event "thinkdev.app/think/runex/runexapi/api/v1/event"
	handle_runHistory "thinkdev.app/think/runex/runexapi/api/v1/runHistory"

	//handle_importData "thinkdev.app/think/runex/runexapi/api/v1/importdata"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	handle_register "thinkdev.app/think/runex/runexapi/api/v1/register"
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

	//api.Use(static.Serve("/img", static.LocalFile("./img", true)))

	EventRoute(route, connectionDB, middleware)
	ActivityRoute(route, connectionDB, middleware)
	CategoryRoute(route, connectionDB, middleware)
	RegisterRoute(route, connectionDB, middleware)
	BannerRoute(route, connectionDB, middleware)
	CouponRoute(route, connectionDB, middleware)
	RunHistoryRoute(route, connectionDB, middleware)
	AdminRoute(route, connectionDB, middleware)
	BoardRoute(route, connectionDB, middleware)

	//ActivityV2Route(route, connectionDB, middleware)
	//WorkoutsRoute(route, connectionDB, middleware)
	//ImportDataRoute(route, connectionDB, middleware)
}

// EventRoute for manage Event
func EventRoute(route *gin.Engine, connectionDB *mongo.Database, middleware *jwt.GinJWTMiddleware) {
	eventRepository := repository.EventRepositoryMongo{
		ConnectionDB: connectionDB,
	}
	eventAPI := handle_event.EventAPI{
		EventRepository: &eventRepository,
	}

	api := route.Group("/api/v1/event")
	{
		api.GET("/findByStatus/:status", eventAPI.GetByStatus)
		api.GET("/eventInfo/:id", eventAPI.GetByID)
		api.GET("/all", eventAPI.GetAll)
		api.GET("/active", eventAPI.GetAllActive)

		api.Use(middleware.MiddlewareFunc())
		{

			api.GET("/myEvent", eventAPI.MyEvent)
			api.POST("", eventAPI.AddEvent)
			api.PUT("/edit/:id", eventAPI.EditEvent)
			api.DELETE("/delete/:id", eventAPI.DeleteEvent)
			api.POST("/:id/uploadImage", eventAPI.UploadImage)
			api.POST("/:id/addProduct", eventAPI.AddProduct)
			api.POST("/:id/editProduct", eventAPI.EditProduct)
			api.GET("/getProduct/:id", eventAPI.GetProductEvent)

			api.DELETE("/deleteProduct/:id/:productID", eventAPI.DeleteProductEvent)
			api.POST("/:id/addTicket", eventAPI.AddTicket)
			api.POST("/:id/editTicket", eventAPI.EditTicket)
			api.DELETE("/deleteTicket/:id/:ticketID", eventAPI.DeleteTicketEvent)
		}
	}
	api2 := route.Group("/api/v1")
	{
		api2.POST("/search/event", eventAPI.SearchEvent)
	}

}

// ActivityRoute for manage activity
func ActivityRoute(route *gin.Engine, connectionDB *mongo.Database, middleware *jwt.GinJWTMiddleware) {
	activityRepository := repository.ActivityRepositoryMongo{
		ConnectionDB: connectionDB,
	}
	activityAPI := handle_activity.ActivityAPI{
		ActivityRepository: &activityRepository,
	}
	api := route.Group("/api/v1/activity")
	{
		api.Use(middleware.MiddlewareFunc())
		{
			api.POST("/add", activityAPI.AddActivity)
			api.POST("/multiadd", activityAPI.AddMultiActivity)
			//api.POST("/byEvent", activityAPI.GetActivityByEvent)
			api.GET("/getByEvent/:event", activityAPI.GetActivityByEvent)
			api.GET("/getByEvent2/:event", activityAPI.GetActivityByEvent2)
			api.POST("/getHistoryDay", activityAPI.GetHistoryDayByEvent)
			api.POST("/getHistoryMonth", activityAPI.GetHistoryMonthByEvent)
			api.DELETE("/deleteActivity/:id/:activityID", activityAPI.DeleteActivityEvent)
			api.GET("/getAllEventActivity/:event_id", activityAPI.GetActivityAllInfo)
		}
	}
}

// CategoryRoute for manage category
func CategoryRoute(route *gin.Engine, connectionDB *mongo.Database, middleware *jwt.GinJWTMiddleware) {
	categoryRepository := repository.CategoryRepositoryMongo{
		ConnectionDB: connectionDB,
	}
	categoryAPI := handle_category.CategoryAPI{
		CategoryRepository: &categoryRepository,
	}
	api := route.Group("/api/v1/category")
	{

		api.GET("/all", categoryAPI.GetAll)
		api.Use(middleware.MiddlewareFunc())
		{
			api.POST("/add", categoryAPI.AddCategory)
			api.PUT("/edit/:id", categoryAPI.EditCategory)
			api.DELETE("/delete/:id", categoryAPI.DeleteCategory)

		}
	}
}

// RegisterRoute for manage register
func RegisterRoute(route *gin.Engine, connectionDB *mongo.Database, middleware *jwt.GinJWTMiddleware) {
	registerRepository := repository.RegisterRepositoryMongo{
		ConnectionDB: connectionDB,
	}
	registerAPI := handle_register.RegisterAPI{
		RegisterRepository: &registerRepository,
	}
	api := route.Group("/api/v1/register")
	{

		api.Use(middleware.MiddlewareFunc())
		{
			api.GET("/all", registerAPI.GetAll)
			api.POST("/add", registerAPI.AddRegister)
			api.POST("/addRace", registerAPI.AddRaceRegister)
			api.PUT("/edit/:id", registerAPI.EditRegister)
			api.GET("/findByEvent/:eventID", registerAPI.GetByEvent)
			api.POST("/sendSlip/:id", registerAPI.SendSlipTransfer)
			api.POST("/adminUpSlip/:id", registerAPI.AdminUpSlip)
			api.GET("/countRegisterEvent/:eventID", registerAPI.CountRegisterEvent)
			api.POST("/sendMail", registerAPI.SendMailRegister)
			api.GET("/checkUserRegisterEvent/:eventID", registerAPI.CheckUserRegisterEvent)
			api.GET("/myRegEvent", registerAPI.GetByUserID)
			api.GET("/getRegEvent/:regID", registerAPI.GetRegEvent)
			api.GET("/myRegEventActivate", registerAPI.GetMyRegEventActivate)
			api.POST("/payment", registerAPI.ChargeRegEvent)
			api.POST("/report", registerAPI.GetReport)
			api.POST("/reportAll", registerAPI.GetReportAll)
			api.POST("/findPersonRegEvent", registerAPI.FindPersonRegEvent)
			api.PUT("/updateStatus", registerAPI.UpdateStatus)
		}
	}
}

// BannerRoute for manage banner
func BannerRoute(route *gin.Engine, connectionDB *mongo.Database, middleware *jwt.GinJWTMiddleware) {
	bannerRepository := repository.BannerRepositoryMongo{
		ConnectionDB: connectionDB,
	}
	bannerAPI := handle_banner.BannerAPI{
		BannerRepository: &bannerRepository,
	}
	api := route.Group("/api/v1/banner")
	{

		api.GET("/all", bannerAPI.GetAll)
		//api.GET("/testmail", bannerAPI.Testmail)
		api.Use(middleware.MiddlewareFunc())
		{
			api.POST("/add", bannerAPI.AddBanner)
			api.DELETE("/delete/:id", bannerAPI.DeleteBanner)

		}
	}
}

// CouponRoute for manage coupon
func CouponRoute(route *gin.Engine, connectionDB *mongo.Database, middleware *jwt.GinJWTMiddleware) {
	couponRepository := repository.CouponRepositoryMongo{
		ConnectionDB: connectionDB,
	}
	couponAPI := handle_coupon.CouponAPI{
		CouponRepository: &couponRepository,
	}
	api := route.Group("/api/v1/coupon")
	{

		api.Use(middleware.MiddlewareFunc())
		{
			api.GET("/couponInfo/:code", couponAPI.GetByCode)
			api.GET("/all", couponAPI.GetAll)
			api.POST("/create", couponAPI.CreateCoupon)
			api.PUT("/edit/:id", couponAPI.EditCoupon)
			api.DELETE("/delete/:id", couponAPI.DeleteCoupon)
			api.POST("/validate", couponAPI.ValidateCode)
		}
	}
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

//BoardRoute ready board
func BoardRoute(route *gin.Engine, connectionDB *mongo.Database, middleware *jwt.GinJWTMiddleware) {
	boardRepository := repository.BoardRepositoryMongo{
		ConnectionDB: connectionDB,
	}
	boardAPI := board.BoardAPI{
		BoardRepository: &boardRepository,
	}
	api := route.Group("/api/v1/board")
	{

		api.Use(middleware.MiddlewareFunc())
		{
			api.GET("/ranking/:eventID", boardAPI.GetBoardByEvent)

		}
	}
}

//AdminRoute route
func AdminRoute(route *gin.Engine, connectionDB *mongo.Database, middleware *jwt.GinJWTMiddleware) {
	adminRepository := repository.AdminRepositoryMongo{
		ConnectionDB: connectionDB,
	}
	adminAPI := handle_admin.AdminAPI{
		AdminRepository: &adminRepository,
	}
	api := route.Group("/api/v1/admin")
	{

		api.Use(middleware.MiddlewareFunc())
		{
			api.PUT("/editShipping/:id", adminAPI.ChangeShppingAddress)
			api.POST("/updateSlip", adminAPI.UpdateSlip)
			api.GET("/getSlip/:regID", adminAPI.GetSlip)
			api.GET("/getRegEvent/:regID", adminAPI.GetRegEvent)
		}
	}
}

// ImportDataRoute for manage banner
// func ImportDataRoute(route *gin.Engine, connectionDB *mongo.Database, middleware *jwt.GinJWTMiddleware) {
// 	importDataRepository := repository.ImportDataRepositoryMongo{
// 		ConnectionDB: connectionDB,
// 	}
// 	importDataAPI := handle_importData.ImportDataAPI{
// 		ImportDataRepository: &importDataRepository,
// 	}
// 	api := route.Group("/api/v1/importdata")
// 	{

// 		//api.GET("/excel/:event", importDataAPI.ImportExcel)
// 		// //api.GET("/testmail", bannerAPI.Testmail)
// 		// api.Use(middleware.MiddlewareFunc())
// 		// {
// 		// 	api.POST("/add", bannerAPI.AddBanner)
// 		// 	api.DELETE("/delete/:id", bannerAPI.DeleteBanner)

// 		// }
// 	}
// }

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

// func WorkoutsRoute(route *gin.Engine, connectionDB *mongo.Database, middleware *jwt.GinJWTMiddleware) {
// 	workoutsRepository := repository.WorkoutsRepositoryMongo{
// 		ConnectionDB: connectionDB,
// 	}
// 	workoutsAPI := handle_workouts.WorkoutsAPI{
// 		WorkoutsRepository: &workoutsRepository,
// 	}
// 	api := route.Group("/api/v1/workout")
// 	{
// 		api.Use(middleware.MiddlewareFunc())
// 		{
// 			api.POST("/add", workoutsAPI.AddWorkout)
// 		}
// 	}
// }
