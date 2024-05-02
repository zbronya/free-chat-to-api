package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zbronya/free-chat-to-api/chat"
	"github.com/zbronya/free-chat-to-api/logger"
	"net/http"
)

func Cors(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}
	c.Next()
}

func RequestLogger(c *gin.Context) {
	c.Next()

	logger.GetLogger().Info(fmt.Sprintf("IP: %s, Method: %s, UserAgent: %s, URL: %s, Status: %d",
		c.Request.RemoteAddr,
		c.Request.Method,
		c.Request.UserAgent(),
		c.Request.URL.Path,
		c.Writer.Status()))
}

func InitRouter(g *gin.Engine) {
	g.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK,
			"This is Free Chat To API \n"+
				"https://github.com/zbronya/free-chat-to-api")
	})

	g.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	g.Use(Cors)
	g.Use(RequestLogger)
	g.OPTIONS("/v1/chat/completions")
	g.POST("/v1/chat/completions", chat.Completions)
	g.POST("/backend-anon/sentinel/chat-requirements", chat.ReverseProxy)
	g.POST("/backend-anon/conversation", chat.ReverseProxy)
}
