package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/hyperledger/fabric/common/flogging"
	"gitlab.com/jaderabbit/go-rabbit/i18n"
)

var logger = flogging.MustGetLogger("tee.api.controllers")

func responseJSON(c beego.Controller, success bool, data interface{}) {
	responseData := map[string]interface{}{
		"success": success,
		"result":  "",
		"code":    "000000",
		"message": "",
	}
	if !success {
		if err, ok := data.(i18n.Error); !ok {
			responseData["message"] = data
		} else {
			responseData["code"] = err.Code()
			responseData["message"] = err.Error()
		}
	} else {
		responseData["result"] = data
	}

	output, err := json.Marshal(responseData)
	if err != nil {
		logger.Errorf("Error marshal responseData: %s", err)
	}
	c.Ctx.ResponseWriter.Write(output)
}
