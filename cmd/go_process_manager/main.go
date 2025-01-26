package main

import (
	_ "github.com/lzh-1625/go_process_manager/boot"
	"github.com/lzh-1625/go_process_manager/internal/app/route"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	route.Route()
}
