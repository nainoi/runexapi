package routers

import (
	handle_activity "bitbucket.org/suthisakch/runex/api/v1/activity"
	handle_event "bitbucket.org/suthisakch/runex/api/v1/event"
	"bitbucket.org/suthisakch/runex/repository"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func EventRoute(route *gin.Engine, connectionDB *mongo.Database) {
	eventRepository := repository.EventRepositoryMongo{
		ConnectionDB: connectionDB,
	}
	eventAPI := handle_event.EventAPI{
		EventRepository: &eventRepository,
	}
	api := route.Group("/api/v1/event")
	{
		api.POST("", eventAPI.AddEvent)
		api.GET("/findByStatus/:status", eventAPI.GetByStatus)
		api.GET("/eventInfo/:id", eventAPI.GetByID)
		api.PUT("/edit/:id", eventAPI.EditEvent)
		api.DELETE("/delete/:id", eventAPI.DeleteEvent)
		api.POST("/:id/uploadImage", eventAPI.UploadImage)
		api.GET("/all/", eventAPI.GetAll)
		api.POST("/:id/addProduct", eventAPI.AddProduct)
		api.DELETE("/deleteProduct/:id/:productID", eventAPI.DeleteProductEvent)

	}
}

func ActivityRoute(route *gin.Engine, connectionDB *mongo.Database) {
	activityRepository := repository.ActivityRepositoryMongo{
		ConnectionDB: connectionDB,
	}
	activityAPI := handle_activity.ActivityAPI{
		ActivityRepository: &activityRepository,
	}
	api := route.Group("/api/v1/activity")
	{
		api.POST("/add", activityAPI.AddActivity)
		//api.POST("/byEvent", activityAPI.GetActivityByEvent)
		api.GET("/getByEvent/:event", activityAPI.GetActivityByEvent)
		api.GET("/getByEvent2/:event", activityAPI.GetActivityByEvent2)

	}
}
