package main

import (
	_ "beeblog/routers"
	"beeblog/models"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"beeblog/controllers"
	"os"
)

func init()  {
	models.RegisterDB()
}

func main() {
	orm.Debug = true

	// 创建 table
	orm.RunSyncdb("default", false, true) //自动建表

	// 注册beego路由
	beego.Router("/", &controllers.HomeController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/category", &controllers.CategoryController{})
	beego.Router("/reply", &controllers.ReplyController{})
	beego.Router("/reply/add", &controllers.ReplyController{}, "post:Add")
	beego.Router("/reply/delete", &controllers.ReplyController{}, "get:Delete")

	beego.Router("/topic", &controllers.TopicController{})
	beego.AutoRouter(&controllers.TopicController{}) //自动路由

	//创建附件目录
	os.Mkdir("attachment", os.ModePerm)
	////作为静态文件
	//beego.SetStaticPath("/attachment", "attachment")

	beego.Router("/attachment/:all", &controllers.AttachController{})


	beego.Run()
}

