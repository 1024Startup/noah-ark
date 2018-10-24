package controllers

import (
	"github.com/astaxie/beego"
	"noah/models"
	"gopkg.in/mgo.v2/bson"
)

type CategoryController struct {
	beego.Controller
}

func (category *CategoryController) Get() {

	//page, _ := app.GetInt("page")
	//category, _ := app.GetInt("category")
	categoryModel:= models.GetCategoryInstance()

	result, err := categoryModel.GetCategory()
//	category.Ctx.WriteString(json)

	if err != nil {
		panic(err.Error())
		category.Data["json"] = bson.M{"ret": -1}
	} else {
		category.Data["json"] = bson.M{"ret": 0, "data": result}
	}
	category.ServeJSON()
}

func (category *CategoryController) GetSummary() {
	category.Ctx.WriteString("fuck")
}
