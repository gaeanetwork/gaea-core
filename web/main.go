package main

import (
	beegoserver "github.com/gaeanetwork/gaea-core/api/beego"
	"github.com/gaeanetwork/gaea-core/common/config"
	"github.com/gaeanetwork/gaea-core/smartcontract/fabric"
	"github.com/gaeanetwork/gaea-core/smartcontract/factory"
	"github.com/gaeanetwork/gaea-core/tee/server"
	ginserver "github.com/gaeanetwork/gaea-core/web/gin"
)

func main() {
	// TODO - remove beego server
	go beegoserver.Start()
	go factory.InitSmartContractService(&fabric.Chaincode{})
	go server.NewTeeServer(config.GRPCAddr).Start()
	ginserver.Start()
}
