package routers

import (
	"github.com/zrynuaa/ABEPasswordPlatform/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})

	beego.Router("/", &controllers.SignUpController{})
	beego.Router("ChangePassword", &controllers.SignUpController{})
	beego.Router("GetTip", &controllers.GetTipController{})
}
