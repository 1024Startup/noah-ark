package controllers

import (
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
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
	rowset []bson.M
}

func (app *AppController) GetList() {
	page, _ := app.GetInt("page")
	category, _ := app.GetInt("category")

	var appModel models.Apps

	result, err := appModel.GetList(page, category)
	if err != nil {
		panic(err.Error())
		app.Data["json"] = bson.M{"ret": -1}
	} else {
		app.Data["json"] = bson.M{"ret": 0, "rowset": result}
	}
	app.ServeJSON()
}
