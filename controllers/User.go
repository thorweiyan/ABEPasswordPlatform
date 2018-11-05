package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/cloudflare/cfssl/scan/crypto/sha256"
	"github.com/thorweiyan/ABEPasswordPlatform/chaincodeImpl/wrapper"
	"strings"
)

type user struct {
	UserName             string		`json:"user_name"`
	UserPasswordHash     string		//"xxxxxxxxx"
	ChangePasswordPolicy string		//"CPP:xxxxx"
	GetTipPolicy         string		//"GTP:xxxxx"
	GetTipMessage        string		//"GTM:xxxxx"
	UserAttributes       string		//"xxxxxxxxx"
}


type SignUpController struct {
	beego.Controller
}

func (c *SignUpController)Get()  {
	c.TplName = "user.html"
}

func (c *SignUpController)Post() {
	u := user{}
	if err := c.ParseForm(&u); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("解析成功！" + u.UserName, u.UserPasswordHash,u.ChangePasswordPolicy,u.GetTipPolicy, u.GetTipMessage, u.UserAttributes)

		userdata := wrapper.UserData{
			UserName:             "UN:" + u.UserName,
			UserPasswordHash:     []byte(fmt.Sprint(sha256.Sum256([]byte(u.UserPasswordHash)))),
			ChangePasswordPolicy: "CPP:" + u.ChangePasswordPolicy,
			GetTipPolicy:         "GTP:" + u.GetTipPolicy,
			GetTipMessage:        "GTM:" + u.GetTipMessage,
			UserAttributes:       strings.Split(u.UserAttributes, ","),
		}
		fmt.Println(userdata)

		DoSdk(userdata, "UserSignUp")

		//正确执行，返回200
		c.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
		c.Ctx.WriteString("200")
	}
}


type ChangePasswordController struct {
	beego.Controller
}

func (c * ChangePasswordController)Get()  {
	c.TplName = "user.html"
}
func (c *ChangePasswordController)Post()  {
	//TODO 处理前端信息得到name、hash和属性集
	u := user{}
	if err := c.ParseForm(&u); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("解析： ", u.UserName, u.UserPasswordHash, u.UserAttributes)

		userdata := wrapper.UserData{
			UserName:         "UN:" + u.UserName,
			UserPasswordHash: []byte(fmt.Sprint(sha256.Sum256([]byte(u.UserPasswordHash)))),
			UserAttributes:   strings.Split(u.UserAttributes, ","),
		}
		DoSdk(userdata, "UserChangePassword")

		//正确执行，返回200
		c.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
		c.Ctx.WriteString("200")
	}
}

type GetTipController struct {
	beego.Controller
}

func (c *GetTipController)Get()  {
	c.TplName = "user.html"
}

func (c *GetTipController)Post() {
	u := user{}
	if err := c.ParseForm(&u); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("解析： ", u.UserName, u.UserAttributes)

		userdata := wrapper.UserData{
			UserName:       "UN:" + u.UserName,
			UserAttributes: strings.Split(u.UserAttributes, ","),
		}
		result := DoSdk(userdata, "UserGetTip")

		c.Ctx.WriteString("tips: " + result)
	}
}