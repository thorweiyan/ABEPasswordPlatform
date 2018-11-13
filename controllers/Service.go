package controllers

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/thorweiyan/ABEPasswordPlatform/models"
	"math/big"
	"strconv"
)

type ApplyCertificatesController struct {
	beego.Controller
}

type LoginController struct {
	beego.Controller
}

func (c *LoginController)Post() {
	u := user{}
	json.NewDecoder(c.Ctx.Request.Body).Decode(&u)

	fmt.Println("解析： ", u.UserName, u.UserPasswordHash)

	//调用special AA合约
	index := Rand2(big.NewInt(0)).Int64()
	sccId := "AA_" + strconv.Itoa(int(index)) + "cc"
	fmt.Println("special ccId: ", sccId)
	sOwnerPriKey, _ := hex.DecodeString(AAkey[index-1].prikey)

	passwordhash := []byte(u.UserPasswordHash)

	models.SdkThirdParty(sccId, sOwnerPriKey, u.UserName, passwordhash)

	//正确执行，返回200
	c.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	c.Ctx.WriteString("200")
}