package routers

import (
	"github.com/thorweiyan/ABEPasswordPlatform/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})

	beego.Router("/signup", &controllers.SignUpController{})
	beego.Router("/changepassword", &controllers.SignUpController{})
	beego.Router("/gettip", &controllers.GetTipController{})

    beego.Router("login", &controllers.LoginController{})
    beego.Router("applycertificates", &controllers.ApplyCertificatesController{})
}
