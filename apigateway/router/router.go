package router

import (
	"apigw/handler"
	"apigw/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	route := gin.Default()
	route.Use(middleware.SetTracer(), middleware.SetHystrixBreaker())
	route.POST("/account/register", handler.RegisterHandler)
	return route
}
