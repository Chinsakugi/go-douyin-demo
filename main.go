package main

import (
	"github.com/RaymondCode/simple-demo/config"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	initRouter(r)

	r.Run(config.Cfg.ServerConfig.Addr) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
