package main

import (
	"github.com/gin-gonic/gin"
	"go-douyin-demo/config"
)

func main() {
	r := gin.Default()

	initSwagger()
	initRouter(r)

	r.Run(config.Cfg.ServerConfig.Addr) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
