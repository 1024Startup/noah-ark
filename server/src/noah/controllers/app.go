package controllers

import (
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
	"noah/models"
	"log"
	"encoding/json"
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
	page := app.GetString("page")
	category := app.GetString("category")

	app.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	app.Ctx.Output.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,DELETE")
	app.Ctx.Output.Header("Access-Control-Allow-Headers", "x-requested-with,content-type")
	app.Ctx.Output.Header("Access-Control-Allow-Credentials", "true")

	var appModel models.Apps

	result, err := appModel.GetList(page, category)
	if err != nil {
		panic(err.Error())
		app.Data["json"] = bson.M{"ret": -1}
	} else {
		app.Data["json"] = bson.M{"ret": 0, "data": bson.M{"rowset": result}}
	}
	app.ServeJSON()
}

func (app *AppController) PostSave() {
	appdata := models.AppsData{}
	app.ParseForm(&appdata)
	btons := app.Input().Get("buttons")



	var buttons []models.AppButton
	json.Unmarshal([]byte(btons), &buttons)
	appdata.Buttons = buttons
	log.Println("dddd::::V%,,,,v%",appdata,buttons)
	appModel := models.Apps{}
	err := appModel.Add(appdata)
	if err {
		app.Data["json"] = bson.M{"ret": -1}
	} else {
		app.Data["json"] = bson.M{"ret": 0}
	}
	app.ServeJSON()
}
