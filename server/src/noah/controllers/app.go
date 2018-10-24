package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"noah/models"
)

type AppController struct {
	beego.Controller
}

func (app *AppController) Get() {
	app.Ctx.WriteString("nothing here")
}

func (app *AppController) Post() {
	//var appModel models.Apps
}

type appReturn struct {
	ret    int
	rowset []models.Apps
}

func (app *AppController) GetList() {
	page, _ := app.GetInt("page")
	category, _ := app.GetInt("category")

	var appModel models.Apps

	result, err := appModel.GetList(page, category)
	if err != nil {
		print(err)
		app.Data["json"] = appReturn{-1, nil}
	} else {
		app.Data["json"] = appReturn{0, result}
	}
	fmt.Print(result)
	app.ServeJSON()
}
