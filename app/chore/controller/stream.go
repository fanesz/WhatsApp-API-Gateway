package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *whatsappController) RegisterStream(router *gin.Engine) {
	v1 := router.Group("/v1")

	v1.GET("/whatsapp/login-qr", func(ctx *gin.Context) {
		flusher, ok := ctx.Writer.(http.Flusher)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Streaming unsupported",
			})
			return
		}

		ctx.Writer.Header().Set("Content-Type", "image/png")
		ctx.Writer.Header().Set("Cache-Control", "no-cache")
		ctx.Writer.Header().Set("Connection", "keep-alive")

		qrChan, err := c.service.GetLoginQR()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		if qrChan == nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to get QR code",
			})
			return
		}

		doneChan := ctx.Request.Context().Done()

		for {
			select {
			case qrImage, ok := <-qrChan:
				if !ok {
					return
				}
				if qrImage == nil {
					continue
				}

				_, err := ctx.Writer.Write(*qrImage)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"message": err.Error(),
					})
					return
				}

				flusher.Flush()

			case <-doneChan:
				fmt.Println("Client cancelled the QR login request")
				return
			}
		}
	})
}
