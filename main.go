package main

import (
	_ "beeblog/routers"
	"beeblog/models"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"beeblog/controllers"
)

func init()  {
	models.RegisterDB()
}

func main() {
	orm.Debug = true

	// 创建 table
	orm.RunSyncdb("default", false, true) //自动建表

	// 注册beego路由
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/category", &controllers.CategoryController{})
	beego.Router("/topic", &controllers.TopicController{})
	beego.AutoRouter(&controllers.TopicController{}) //自动路由

	beego.Run()
}

