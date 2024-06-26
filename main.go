package main

import (
	"msm/boot"
	"msm/route"

	"github.com/gin-gonic/gin"
)

func main() {
	boot.Boot()
	// go termui.TermuiInit()
	gin.SetMode(gin.ReleaseMode)
	route.Route()
}
