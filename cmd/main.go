package main

import (
	"github.com/MaximLanBowl/Searcher.git/handlers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server for driver search.
// @termsOfService https://your-git-repo-url.com

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:5432
// @BasePath /



func main() {
    r := gin.Default()

    // Добавляем эндпоинты
    r.GET("/driverSearch", handlers.SearchDriver)

    // Добавляем Swagger эндпоинты
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    r.Run(":5432")
}