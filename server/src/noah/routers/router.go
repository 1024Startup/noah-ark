package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"noah/controllers"
)

func init() {
	beego.Get("/", func(context *context.Context) {
		context.Redirect(302, "/static/")
	})
	beego.Router("/api/apps/list", &controllers.AppController{}, "get:GetList")
	beego.Router("/api/categories", &controllers.CategoryController{})
	beego.Router("/api/categories/summary", &controllers.CategoryController{}, "get:GetSummary")
}
