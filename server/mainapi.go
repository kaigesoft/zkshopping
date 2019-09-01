package main

import (
	"github.com/gin-gonic/gin"
	"server/config"
)

func main() {
	router := gin.Default()
	config.RouterInit(router)
	router.Run(":8082")
}
