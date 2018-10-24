// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	"owdile-notification/controllers"
)

func init() {
	//beego.Router("/test", &controllers.UserController{}, "get:GetAll;post:Post")

	/*ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
	)
	beego.AddNamespace(ns)*/

	//beego.Router("/test", &controllers.UserController{}, )

	//初始化 namespace
	ns :=
		beego.NewNamespace("/v1",
			beego.NSRouter("/user", &controllers.UserController{}, "get:GetAll"),
			beego.NSRouter("/changepassword", &controllers.UserController{}),
		)
	//注册 namespace
	beego.AddNamespace(ns)

}
