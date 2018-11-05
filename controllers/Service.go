package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/thorweiyan/ABEPasswordPlatform/models"
	"math/big"
	"strconv"
)

type userthird struct {
	UserName string
	UserPasswordHash string
}

type ApplyCertificatesController struct {
	beego.Controller
}

type LoginController struct {
	beego.Controller
}

func (c *LoginController)Get()  {
	c.TplName = "thirdparty.html"
}

func (c *LoginController)Post()  {
	//TODO 处理前端信息得到name、hash和属性集
	u := userthird{}
	if err := c.ParseForm(&u); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("解析： ", u.UserName, u.UserPasswordHash)

		//调用special AA合约
		//TODO 随机出special AA
		index := Rand2(big.NewInt(0)).Int64()
		sccId := "AA_" + strconv.Itoa(int(index))
		fmt.Println("special ccId: ", sccId)
		sOwnerPriKey, _ := hex.DecodeString(AAkey[index-1].prikey)

		passwordhash := []byte(fmt.Sprint(sha256.Sum256([]byte(u.UserPasswordHash))))

		models.SdkThirdParty(sccId, sOwnerPriKey, u.UserName, passwordhash)

		//正确执行，返回200
		c.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
		c.Ctx.WriteString("200")
	}
}