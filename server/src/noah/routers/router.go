package routers

import (
	"github.com/astaxie/beego"
	"noah/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/api/apps/list", &controllers.AppListController{})
	beego.Router("/api/categories", &controllers.CategoryController{})
	beego.Router("/api/categories/summary", &controllers.CategorySummaryController{})
}
