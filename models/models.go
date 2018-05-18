package models

import (
	"time"
	"github.com/Unknwon/com"
	"os"
	"path"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3" //导入驱动
)

const (
	_DB_NAME	    = "data/beeblog"
	_SQLITES_DRIVER = "sqlite3"
)

type Category struct {
	Id		 		int64
	Title    		string
	Created 		time.Time   `orm:"index"`
	Views   		int64     	`orm:"index"`
	TopicTime       time.Time   `orm:"index"`
	TopicCount 		int64
	TopicLastUserId int64
}

type Topic struct {
	Id 				int64
	Uid 			int64
	Title 			string
	Content 		string      `orm:"size(5000)"`
	Attachment 		string
	Created		    time.Time   `orm:"index"`
	Updated 	 	time.Time   `orm:"index"`
	views  			int64       `orm:"index"`
	Author 			string
	ReplyTime 		time.Time   `orm:"index"`
	ReplyCount 		int64
	ReplyLastUserId int64
}

func RegisterDB()  {
	if !com.IsExist(_DB_NAME) {
		os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm)
		os.Create(_DB_NAME)
	}

	// 注册定义的 model
	orm.RegisterModel(new(Category), new(Topic))
	//注册驱动
	orm.RegisterDriver(_SQLITES_DRIVER, orm.DRSqlite)
	// 设置默认数据库
	//数据库存放位置：_DB_NAME ， 数据库别名：default
	orm.RegisterDataBase("default", _SQLITES_DRIVER, _DB_NAME, 10)

}