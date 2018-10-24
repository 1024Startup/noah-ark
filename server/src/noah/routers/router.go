package routers

import (
	"github.com/astaxie/beego"
	"noah/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
}
