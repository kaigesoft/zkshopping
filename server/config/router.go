package config

import (
	"net/http"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"server/model"
)

func RouterInit(router *gin.Engine) {
	router.Use(cors.Default())
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "It works On 8082")
	})
	router.GET("/uList", model.GetList)
}