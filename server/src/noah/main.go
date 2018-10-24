package main

import (
	"github.com/astaxie/beego"
	_ "noah/routers"
)

func main() {
	beego.BConfig.WebConfig.DirectoryIndex = true
	beego.Run()
}
