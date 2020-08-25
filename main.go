package main

import (
	//"io/ioutil"
	// "fmt"
	"log"

	// "context"
	"net/http"
	"os"
	"time"

	"bitbucket.org/suthisakch/runex/config"
	"bitbucket.org/suthisakch/runex/config/db"
	"bitbucket.org/suthisakch/runex/docs"
	routes "bitbucket.org/suthisakch/runex/route"

	//routes2 "bitbucket.org/suthisakch/runex/routers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	//cors "github.com/rs/cors/wrapper/gin"

	//"github.com/swaggo/swag" // gin-swagger middleware
	// swagger embed files
	_ "bitbucket.org/suthisakch/runex/docs"
	"github.com/gin-contrib/cors"
	//"./docs"
	//"bitbucket.org/suthisakch/runex/middleware"
)

// type key string

// const (
// 	hostKey     = key("hostKey")
// 	usernameKey = key("usernameKey")
// 	passwordKey = key("passwordKey")
// 	databaseKey = key("databaseKey")
// )

var sslkey string = "runex.co.key"
var sslcert string = "runex.co.crt"

func main() {
	os.Setenv("TZ", "Asia/Bangkok")
	// time.FixedZone("UTC+7", +7*60*60)
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization", "Content-Disposition", "Content-Description", "sessionId"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "*"
		},
		MaxAge: 12 * time.Hour,
	}))
	database, err := db.GetDBCollection()
	if err != nil {
		log.Fatalf("todo: database configuration failed: %v", err)
	}

	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not Found")
	})

	//routes.ProjectRoute(router, database)
	routes.Route(router, database)

	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "petstore.swagger.io"
	docs.SwaggerInfo.BasePath = "/v2"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	//url := ginSwagger.URL("http://localhost:3006/swg/api_v2.json")
	router.GET("/swg/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// use ginSwagger middleware to serve the API docs
	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//router.Static("/upload", "./upload")
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Static("/upload", "./upload")
	router.LoadHTMLGlob("templates/*")
	//init the loc
	//loc, _ := time.LoadLocation("Asia/Bangkok")

	router.Run(config.PORT_WEB_SERVICE)

	// err2 := router.RunTLS(config.PORT_WEB_SERVICE, sslcert, sslkey)
	// if err2 != nil {
	// 	log.Println(err2.Error())
	// }
}
