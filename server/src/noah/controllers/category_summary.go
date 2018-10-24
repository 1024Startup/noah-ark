package controllers

import (
	"github.com/astaxie/beego"
)

type CategorySummaryController struct {
	beego.Controller
}

func (this *CategorySummaryController) Get() {
	this.Ctx.WriteString("You can not see me!")
}
