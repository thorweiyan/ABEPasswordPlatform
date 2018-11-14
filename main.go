package main

import (
	"github.com/astaxie/beego/context"
	_ "github.com/thorweiyan/ABEPasswordPlatform/routers"
	"github.com/astaxie/beego"
	"net/http"
	"strings"
)

func main() {
	//这两句是为了可以直接访问静态文件
	//beego.InsertFilter("/", beego.BeforeRouter, TransparentStatic)
	//beego.InsertFilter("/main.html", beego.BeforeRouter, TransparentStatic)
	beego.SetStaticPath("/","static")

	beego.Run()
}

func TransparentStatic(ctx * context.Context) {
	if strings.Index(ctx.Request.URL.Path, "v1/") >= 0 {
		return
	}
	http.ServeFile(ctx.ResponseWriter, ctx.Request, "static/"+ctx.Request.URL.Path)
}