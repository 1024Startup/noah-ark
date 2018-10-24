package routers

import (
	"github.com/astaxie/beego"
	"noah/controllers"
)

func init() {
	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/apps",
			beego.NSRouter("/list", &controllers.AppListController{}),
		),

		beego.NSRouter("/categories", &controllers.CategoryController{}),
		beego.NSNamespace("/categories",
			beego.NSRouter("/summary", &controllers.CategorySummaryController{}),
		),
	)
	beego.AddNamespace(ns)
}
