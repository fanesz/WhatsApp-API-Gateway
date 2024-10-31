package router

import (
	"fmt"
	whatsappController "mywaclient/app/chore/controller"
)

func InitializeRoutes() {
	fmt.Println("===== Initialize Routes =====")
	router := GetRouterInstance()

	// Register the routes
	waController := whatsappController.NewWhatsappController()
	waController.Register(router)
	waController.RegisterStream(router)

	fmt.Println("✓ Initialize", len(router.Routes()), "routes")
}
