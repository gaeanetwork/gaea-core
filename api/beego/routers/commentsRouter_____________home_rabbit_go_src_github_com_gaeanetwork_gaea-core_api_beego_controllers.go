package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/gaeanetwork/gaea-core/api/beego/controllers:SharedDataController"] = append(beego.GlobalControllerRouter["github.com/gaeanetwork/gaea-core/api/beego/controllers:SharedDataController"],
        beego.ControllerComments{
            Method: "Upload",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/gaeanetwork/gaea-core/api/beego/controllers:SharedDataController"] = append(beego.GlobalControllerRouter["github.com/gaeanetwork/gaea-core/api/beego/controllers:SharedDataController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/gaeanetwork/gaea-core/api/beego/controllers:SharedDataController"] = append(beego.GlobalControllerRouter["github.com/gaeanetwork/gaea-core/api/beego/controllers:SharedDataController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/gaeanetwork/gaea-core/api/beego/controllers:SharedDataController"] = append(beego.GlobalControllerRouter["github.com/gaeanetwork/gaea-core/api/beego/controllers:SharedDataController"],
        beego.ControllerComments{
            Method: "Update",
            Router: `/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/gaeanetwork/gaea-core/api/beego/controllers:SharedDataController"] = append(beego.GlobalControllerRouter["github.com/gaeanetwork/gaea-core/api/beego/controllers:SharedDataController"],
        beego.ControllerComments{
            Method: "Authorize",
            Router: `/authorize/`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/gaeanetwork/gaea-core/api/beego/controllers:SharedDataController"] = append(beego.GlobalControllerRouter["github.com/gaeanetwork/gaea-core/api/beego/controllers:SharedDataController"],
        beego.ControllerComments{
            Method: "GetHistory",
            Router: `/history/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/gaeanetwork/gaea-core/api/beego/controllers:SharedDataController"] = append(beego.GlobalControllerRouter["github.com/gaeanetwork/gaea-core/api/beego/controllers:SharedDataController"],
        beego.ControllerComments{
            Method: "GetAllNotificaitons",
            Router: `/notification/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/gaeanetwork/gaea-core/api/beego/controllers:SharedDataController"] = append(beego.GlobalControllerRouter["github.com/gaeanetwork/gaea-core/api/beego/controllers:SharedDataController"],
        beego.ControllerComments{
            Method: "Request",
            Router: `/request/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/gaeanetwork/gaea-core/api/beego/controllers:SharedDataController"] = append(beego.GlobalControllerRouter["github.com/gaeanetwork/gaea-core/api/beego/controllers:SharedDataController"],
        beego.ControllerComments{
            Method: "GetAllRequests",
            Router: `/request/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/gaeanetwork/gaea-core/api/beego/controllers:TeeTaskController"] = append(beego.GlobalControllerRouter["github.com/gaeanetwork/gaea-core/api/beego/controllers:TeeTaskController"],
        beego.ControllerComments{
            Method: "Create",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/gaeanetwork/gaea-core/api/beego/controllers:TeeTaskController"] = append(beego.GlobalControllerRouter["github.com/gaeanetwork/gaea-core/api/beego/controllers:TeeTaskController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/gaeanetwork/gaea-core/api/beego/controllers:TeeTaskController"] = append(beego.GlobalControllerRouter["github.com/gaeanetwork/gaea-core/api/beego/controllers:TeeTaskController"],
        beego.ControllerComments{
            Method: "Execute",
            Router: `/`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/gaeanetwork/gaea-core/api/beego/controllers:TeeTaskController"] = append(beego.GlobalControllerRouter["github.com/gaeanetwork/gaea-core/api/beego/controllers:TeeTaskController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
