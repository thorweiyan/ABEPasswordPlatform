package controllers

import "github.com/astaxie/beego"

type SignUpController struct {
	beego.Controller
}

func (c *SignUpController)Get()  {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "register.html"
}


type ChangePasswordController struct {
	beego.Controller
}


type GetTipController struct {
	beego.Controller
}