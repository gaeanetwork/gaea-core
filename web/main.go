package main

import (
	beegoserver "github.com/gaeanetwork/gaea-core/api/beego"
	ginserver "github.com/gaeanetwork/gaea-core/web/gin"
)

func main() {
	// TODO - remove beego server
	go beegoserver.Start()
	ginserver.Start()
}
