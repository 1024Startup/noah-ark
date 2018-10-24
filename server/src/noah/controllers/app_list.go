package controllers

import (
	"github.com/astaxie/beego"
)

type AppListController struct {
	beego.Controller
}

func (this *AppListController) Get() {
	this.Ctx.WriteString("App list")
}
