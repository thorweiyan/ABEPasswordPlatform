package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Lab"] = "复旦区块链实验室"
	c.Data["Email"] = "zry_nuaa@163.com"
	c.TplName = "index.tpl"
}
