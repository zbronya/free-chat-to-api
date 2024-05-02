package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zbronya/free-chat-to-api/chat"
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

func InitRouter(g *gin.Engine) {
	g.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "This is Free Chat To API \n"+
			"https://github.com/zbronya/free-chat-to-api")
	})

	g.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	g.Use(Cors)
	g.OPTIONS("/v1/chat/completions")
	g.POST("/v1/chat/completions", chat.Completions)
}
