package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/thorweiyan/ABEPasswordPlatform/chaincodeImpl/wrapper"
	"net/http"
	"strings"
)

type user struct {
	UserName             string		`json:"UserName"`
	UserPasswordHash     string		`json:"UserPasswordHash"`//"xxxxxxxxx"
	ChangePasswordPolicy string		`json:"ChangePasswordPolicy"`//"CPP:xxxxx"
	GetTipPolicy         string		`json:"GetTipPolicy"`//"GTP:xxxxx"
	GetTipMessage        string		`json:"GetTipMessage"`//"GTM:xxxxx"
	UserAttributes       string		`json:"UserAttributes"`//"xxxxxxxxx"
}


type SignUpController struct {
	beego.Controller
}

//func (c *SignUpController)Get()  {
//	c.TplName = "user.html"
//}

func (c *SignUpController)Post() {
	u := user{}

	json.NewDecoder(c.Ctx.Request.Body).Decode(&u)

	userdata := wrapper.UserData{
		UserName:             u.UserName,
		UserPasswordHash:     []byte(u.UserPasswordHash),
		ChangePasswordPolicy: u.ChangePasswordPolicy,
		GetTipPolicy:         u.GetTipPolicy,
		GetTipMessage:        u.GetTipMessage,
		UserAttributes:       strings.Split(u.UserAttributes, ","),
	}
	fmt.Println(userdata)

	_, err := DoSdk(userdata, "userSignUp")

	if err != nil{
		fmt.Fprint(c.Ctx.ResponseWriter, http.StatusForbidden) //403
	}else {
		fmt.Fprint(c.Ctx.ResponseWriter, http.StatusOK)
	}
}


type ChangePasswordController struct {
	beego.Controller
}

//func (c * ChangePasswordController)Get()  {
//	c.TplName = "user.html"
//}
func (c *ChangePasswordController)Post() {
	//TODO 处理前端信息得到name、hash和属性集
	u := user{}
	json.NewDecoder(c.Ctx.Request.Body).Decode(&u)

	fmt.Println("解析： ", u.UserName, u.UserPasswordHash, u.UserAttributes)

	userdata := wrapper.UserData{
		UserName:         u.UserName,
		UserPasswordHash: []byte(u.UserPasswordHash),
		UserAttributes:   strings.Split(u.UserAttributes, ","),
	}
	_, err := DoSdk(userdata, "userChangePassword")

	if err != nil{
		fmt.Fprint(c.Ctx.ResponseWriter, http.StatusForbidden) //403
	}else {
		fmt.Fprint(c.Ctx.ResponseWriter, http.StatusOK)
	}
}

type GetTipController struct {
	beego.Controller
}

//func (c *GetTipController)Get()  {
//	c.TplName = "user.html"
//}

func (c *GetTipController)Post() {
	u := user{}
	json.NewDecoder(c.Ctx.Request.Body).Decode(&u)

	fmt.Println("解析： ", u.UserName, u.UserAttributes)

	userdata := wrapper.UserData{
		UserName:       u.UserName,
		UserAttributes: strings.Split(u.UserAttributes, ","),
	}
	result,err := DoSdk(userdata, "userGetTip")

	if err != nil{
		fmt.Fprint(c.Ctx.ResponseWriter, http.StatusForbidden) //403
	}else {
		c.Ctx.WriteString(result) //返回tip
	}
}