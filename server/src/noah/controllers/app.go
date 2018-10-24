package controllers

import (
	"github.com/astaxie/beego"
)

type AppController struct {
	beego.Controller
}

func (app *AppController) Get() {
	app.Ctx.WriteString("nothing here")
}

type GetListData struct {
}

func (app *AppController) GetList() {
	page, _ := app.GetInt("page")
	category, _ := app.GetInt("category")

	app.Data["json"] = map[string]int{"page": page, "category": category}

	app.ServeJSON()
}
