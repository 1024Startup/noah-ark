package main

import (
	_ "owdile-notification/routers"

	"github.com/astaxie/beego"
	"owdile-notification/models"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	models.Init()

	beego.Run()
}
