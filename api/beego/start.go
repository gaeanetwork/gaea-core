package beegoserver

import (
	"github.com/astaxie/beego"

	// import beego routers
	_ "github.com/gaeanetwork/gaea-core/api/beego/routers"
)

// Start the beego http server
func Start() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "../api/beego/swagger"
		beego.SetStaticPath("/", "dist")
	}

	beego.Run()
}
