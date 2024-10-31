package controller

import (
	"mywaclient/app/chore/entity"
	"mywaclient/app/chore/interfaces"
	"mywaclient/app/chore/service"
	"mywaclient/app/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type whatsappController struct {
	service interfaces.WhatsappService
}

func NewWhatsappController() *whatsappController {
	client := config.GetClient()

	return &whatsappController{
		service: service.NewWhatsappService(client),
	}
}

func (c *whatsappController) Register(router *gin.Engine) {
	v1 := router.Group("/v1")

	v1.GET("/whatsapp/is-login", func(ctx *gin.Context) {
		isLogin, err := c.service.CheckDevice()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"is_login": isLogin,
		})
	})

	v1.POST("/whatsapp/send", func(ctx *gin.Context) {
		var req entity.MessageSend
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body request"})
			return
		}

		// Check if the request body is empty
		if req.To == "" || req.Message == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid body request",
			})
			return
		}

		// Send the message
		err := c.service.SendMessage(&req)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Success sending message",
		})
	})

	v1.GET("/whatsapp/reset", func(ctx *gin.Context) {
		err := c.service.ResetLoggedDevice()
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Success reset device",
		})
	})
}
