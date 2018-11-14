package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
)

type usertest struct {
	Id    int       `form:"-"`
	Name  string 	`form:"username"`
	Age   int       `form:"age"`
	Email string
}

type MainController struct {
	beego.Controller
}

//func (c *MainController) Get() {
//	c.Data["Lab"] = "复旦区块链实验室"
//	c.Data["Email"] = "zry_nuaa@163.com"
//	c.TplName = "index.tpl"
//}

func (c *MainController)Post()  {
	u := usertest{}
	if err := c.ParseForm(&u); err != nil {
		fmt.Println("解析成功！")

	}
	c.Ctx.WriteString("解析成功！" + u.Name + strconv.Itoa(u.Age) + u.Email)
}
