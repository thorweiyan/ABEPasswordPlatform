package routers

import (
	"github.com/thorweiyan/ABEPasswordPlatform/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})

	beego.Router("/", &controllers.SignUpController{})
	beego.Router("ChangePassword", &controllers.SignUpController{})
	beego.Router("GetTip", &controllers.GetTipController{})
}
