// Package routers for trusted execution environment
// @APIVersion 1.0.0
// @Title Trusted Execution Environment API
// @Description swagger has a very cool tools to autogenerate documents for your API
package routers

import (
	"github.com/astaxie/beego"
	"github.com/gaeanetwork/gaea-core/api/beego/controllers"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/tee",
			beego.NSInclude(
				&controllers.SharedDataController{},
			),
		),
		beego.NSNamespace("/task",
			beego.NSInclude(
				&controllers.TeeTaskController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
