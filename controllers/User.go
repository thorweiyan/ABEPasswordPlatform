package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	"github.com/thorweiyan/ABEPasswordPlatform/chaincodeImpl/wrapper"
)

type user struct {
	Id    int       `form:"-"`
	Name  string 	`form:"username"`
	Age   int       `form:"age"`
	Email string
}

type SignUpController struct {
	beego.Controller
}

func (c *SignUpController)Get()  {
	//c.Data["Website"] = "beego.me"
	//c.Data["Email"] = "astaxie@gmail.com"
	//c.TplName = "register.html"

	json.NewEncoder(c.Ctx.ResponseWriter).Encode("success!") //给前端返回数据
}

func (c *SignUpController)Post()  {
	//pkgname := c.GetString("pkgname")
	//content := c.GetString("content")

	jsoninfo := c.GetString("jsoninfo")
	if jsoninfo == "" {
		c.Ctx.WriteString("jsoninfo is empty")
		return
	}

	u := user{}
	if err := c.ParseForm(&u); err != nil {
		fmt.Println(err)
	} else {
		content := "Name:" + u.Name + " Age: " + strconv.Itoa(u.Age) + " Email: " + u.Email
		c.Ctx.WriteString("解析成功" + content)
	}

	userdata := wrapper.UserData{
		UserName:u.Name,
		UserPasswordHash:nil,
		ChangePasswordPolicy:"",
		GetTipPolicy:"",
		GetTipMessage:"",
		UserAttributes:[]string{},
	}

	fmt.Println(userdata)

	//TODO 调用合约，注册新账户

	//var at models.Article
	//at.Pkgid = pk.Id
	//at.Content = content
	//models.InsertArticle(at)
	//this.Ctx.Redirect(302, "/admin/index")
}


type ChangePasswordController struct {
	beego.Controller
}


type GetTipController struct {
	beego.Controller
}

func (c *GetTipController)Get()  {
	
}

func (c *GetTipController)Post()  {

}