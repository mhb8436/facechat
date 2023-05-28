package router

import (
	"errors"
	handler "facechat/api/handler"
	"facechat/config"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Register() *gin.Engine {
	r := gin.Default()
	r.Use(CorsMiddleware())
	// initUserRouter(r)
	initPushRouter(r)
	r.NoRoute(func(c *gin.Context) {
		fmt.Println("please check reqeust url")
	})

	return r
}

func init() {
	if _, err := os.Stat(config.Conf.Api.ApiBase.UploadDir); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(config.Conf.Api.ApiBase.UploadDir, os.ModePerm)
		if err != nil {
			fmt.Println("init error" + err.Error())
		}
	}
}

func initPushRouter(r *gin.Engine) {
	pushGroup := r.Group("/push")
	pushGroup.Use(CheckSessionId())
	{
		pushGroup.POST("/push", handler.Send)
	}
}

func CheckSessionId() gin.HandlerFunc {
	return func(c *gin.Context) {
		return
	}
}
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		var openCorsFlag = true
		if openCorsFlag {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
			c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, nil)
		}
		c.Next()
	}
}
