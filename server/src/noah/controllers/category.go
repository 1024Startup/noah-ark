package controllers

import (
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
	"noah/models"
)

type CategoryController struct {
	beego.Controller
}

func (category *CategoryController) Get() {

	category.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	category.Ctx.Output.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,DELETE")
	category.Ctx.Output.Header("Access-Control-Allow-Headers", "x-requested-with,content-type")
	category.Ctx.Output.Header("Access-Control-Allow-Credentials", "true")

	//page, _ := app.GetInt("page")
	//category, _ := app.GetInt("category")
	categoryModel := models.GetCategoryInstance()

	result, err := categoryModel.GetCategory()
	//	category.Ctx.WriteString(json)

	if err != nil {
		panic(err.Error())
		category.Data["json"] = bson.M{"ret": -1}
	} else {
		category.Data["json"] = bson.M{"ret": 0, "data": bson.M{"rowset": result}}
	}
	category.ServeJSON()
}

func (category *CategoryController) GetSummary() {
	category.Ctx.WriteString("fuck")
}
