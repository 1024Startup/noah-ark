package controllers

import (
	"github.com/astaxie/beego"
	"noah/models"
)

type AppController struct {
	beego.Controller
}

func (app *AppController) Get() {
	app.Ctx.WriteString("nothing here")
}

func (app *AppController) GetList() {
	page, _ := app.GetInt("page")
	category, _ := app.GetInt("category")

	var app_model models.Apps

	json, err := app_model.GetList(page, category)
	if err != nil {
		app.Data["json"] = map[string]string{"ret": "0", "rowset": json}
	} else {
		print(err)
		app.Data["json"] = map[string]string{"ret": "-1"}
	}
	app.ServeJSON()
}
