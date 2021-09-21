package main

import (
	//"io/ioutil"
	// "fmt"
	"fmt"
	"io"
	"log"

	// "context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	//"thinkdev.app/think/runex/runexapi/config"
	"thinkdev.app/think/runex/runexapi/config/db"
	"thinkdev.app/think/runex/runexapi/docs"
	"thinkdev.app/think/runex/runexapi/logger"
	routes "thinkdev.app/think/runex/runexapi/route"

	//cors "github.com/rs/cors/wrapper/gin"

	//"github.com/swaggo/swag" // gin-swagger middleware
	// swagger embed files
	"github.com/gin-contrib/cors"
	stdopentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	_ "thinkdev.app/think/runex/runexapi/docs"
	//"thinkdev.app/think/runex/runexapi/middleware"
)

var sslkey string = "runex.co.key"
var sslcert string = "runex.co.crt"

func initJaeger(service string) (stdopentracing.Tracer, io.Closer) {

	// Uncomment the lines below, if sending traces directly to the collector
	//tracerIP := GetEnv("TRACER_HOST", "localhost")
	//tracerPort := GetEnv("TRACER_PORT", "14268")

	agentIP := db.GetEnv("JAEGER_AGENT_HOST", "localhost")
	agentPort := db.GetEnv("JAEGER_AGENT_PORT", "6831")

	logger.Logger.Infof("Sending traces to %s %s", agentIP, agentPort)

	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: agentIP + ":" + agentPort,
			// Uncomment the line below, if sending traces directly to the collector
			//			CollectorEndpoint: "http://" + tracerIP + ":" + tracerPort + "/api/traces",
		},
	}
	tracer, closer, err := cfg.New(service, config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

func main() {
	os.Setenv("TZ", "Asia/Bangkok")
	// time.FixedZone("UTC+7", +7*60*60)
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err.Error()))
	}

	if !viper.GetBool("app.debug") {
		gin.SetMode(gin.ReleaseMode)
	}

	//create your file with desired read/write permissions
	// f, err := os.OpenFile("log.info", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	// if err != nil {
	// 	fmt.Println("Could not open file ", err)
	// 	logger.Logger.Infof("Could not open file")
	// } else {
	// 	logger.InitLogger(f)
	// }

	// Swagger 2.0 Meta Information
	docs.SwaggerInfo.Title = "RUNEX Aplication - Runex API"
	docs.SwaggerInfo.Description = "RUNEX Aplication - Runex API"
	docs.SwaggerInfo.Version = "2.0"
	docs.SwaggerInfo.Host = viper.GetString("swagger.base_api")
	docs.SwaggerInfo.BasePath = viper.GetString("swagger.base_path")
	docs.SwaggerInfo.Schemes = []string{"https"}

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

	redisClient := db.ConnectRedisDB()
	// tracer, closer := initJaeger("user")

	// stdopentracing.SetGlobalTracer(tracer)

	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not Found")
	})

	routes.Route(router, database)
	routes.Router(router, database)

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Static("/upload", "./upload")
	router.LoadHTMLGlob("templates/*")
	//init the loc
	//loc, _ := time.LoadLocation("Asia/Bangkok")

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//router.Run(viper.GetString("app.port"))

	s := &http.Server{
		Addr:           viper.GetString("app.port"),
		Handler:        router,
		ReadTimeout:    100 * time.Second,
		WriteTimeout:   100 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()

	// err2 := router.RunTLS(config.PORT_WEB_SERVICE, sslcert, sslkey)
	// if err2 != nil {
	// 	log.Println(err2.Error())
	// }

	defer db.CloseDB(database)

	// defer closer.Close()

	// defer f.Close()

	defer redisClient.Close()
}
