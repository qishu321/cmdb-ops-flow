package models

import (
	"cmdb-ops-flow/conf"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)
var db *gorm.DB
var err error

func InitDb() {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.DbUser,
		conf.DbPassWord,
		conf.DbHost,
		conf.DbPort,
		conf.DbName,
	)
	db, err = gorm.Open("mysql", dns)
	if err != nil {
		fmt.Println("数据库连接失败", err)
	}
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(20)  //设置连接池，空闲
	db.DB().SetMaxOpenConns(100) //打开
	db.DB().SetConnMaxLifetime(time.Second * 30)
	db.AutoMigrate(Cmdb{},)
	db.LogMode(true)

}