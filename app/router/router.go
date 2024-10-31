package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

var routerInstance *gin.Engine

func InitializeRouter() {
	fmt.Println("===== Initialize Router =====")

	router := gin.Default()
	router.Use(corsHeaderConfig())
	router.Use(corsConfig())

	routerInstance = router

	fmt.Println("✓ Gin router initialized")
}

func GetRouterInstance() *gin.Engine {
	return routerInstance
}
