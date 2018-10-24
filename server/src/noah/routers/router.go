package routers

import (
	"github.com/astaxie/beego"
	"noah/controllers"
)

func init() {
	beego.Router("/api/apps/list", &controllers.AppController{}, "get:GetList")
	beego.Router("/api/categories", &controllers.CategoryController{})
	beego.Router("/api/categories/summary", &controllers.CategoryController{}, "get:GetSummary")
}
