package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zbronya/free-chat-to-api/config"
	"github.com/zbronya/free-chat-to-api/logger"
	"github.com/zbronya/free-chat-to-api/router"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	g := gin.New()
	g.Use(gin.Recovery())
	router.InitRouter(g)
	host := config.Bind
	logger.GetLogger().Info(fmt.Sprint("Server started on http://", host, ":", config.Port))
	_ = g.Run(fmt.Sprint(config.Bind, ":", config.Port))
}
