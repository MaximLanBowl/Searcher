package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/driverSearch", searchDriver)
	r.Run(":4444")
}