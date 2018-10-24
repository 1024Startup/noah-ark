package controllers

import (
	"github.com/astaxie/beego"
)

type CategoryController struct {
	beego.Controller
}

func (category *CategoryController) Get() {
	category.Ctx.WriteString("Category!")
}

func (category *CategoryController) GetSummary() {
	category.Ctx.WriteString("fuck")
}
