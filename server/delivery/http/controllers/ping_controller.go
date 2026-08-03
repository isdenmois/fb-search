package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PingController struct{}

func (ctrl PingController) ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (ctrl PingController) Bind(r *gin.Engine) error {
	r.GET("/api/ping", ctrl.ping)

	r.Static("/assets", "./public/assets") // Adjust path as needed
	r.NoRoute(func(c *gin.Context) {
		c.File("./public/index.html")
	})

	return nil
}
