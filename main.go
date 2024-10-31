package main

import (
	"mywaclient/app/chore/event"
	"mywaclient/app/config"
	"mywaclient/app/database"
	"mywaclient/app/router"
)

func main() {
	// Initialize the database, config and router
	database.Initialize()
	config.Initialize()
	event.InitializeChatbot()
	router.InitializeRouter()
	router.InitializeRoutes()

	// Get the router instance and run it
	routerInstance := router.GetRouterInstance()
	routerInstance.Run(":5000")
}
