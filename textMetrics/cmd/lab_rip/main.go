package main

import (
	"fmt"

	"lab1_rip/internal/app/config"
	"lab1_rip/internal/app/dsn"
	"lab1_rip/internal/app/handler"
	"lab1_rip/internal/app/repository"
	"lab1_rip/internal/pkg/app"

	_ "lab1_rip/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {
	router := gin.Default()
	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}
	configCORS := cors.DefaultConfig()
	configCORS.AllowOrigins = []string{"*"}
	configCORS.AllowHeaders = []string{"*"}
	configCORS.AllowMethods = []string{"*"}
	configCORS.AllowCredentials = true
	router.Use(cors.New(configCORS))

	postgresString := dsn.FromEnv()
	fmt.Println(postgresString)

	rep, errRep := repository.NewRepository(postgresString, conf)
	if errRep != nil {
		logrus.Fatalf("error initializing repository: %v", errRep)
	}
	redis, err := repository.NewRedis(conf)
	if errRep != nil {
		logrus.Fatalf("error initializing redis: %v", err)
	}
	hand := handler.NewHandler(rep, conf, redis)

	application := app.NewApp(conf, router, hand)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.PersistAuthorization(true)))
	application.RunApp()
}
