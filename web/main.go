package main

import (
	beegoserver "github.com/gaeanetwork/gaea-core/web/beego"
)

func main() {
	// TODO - remove beego server
	go beegoserver.Start()
	ginserver.Start()
}
