package core

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var g_db *gorm.DB

func init() {
	config := Config().Mysql
	if &config == nil {
		panic(fmt.Sprintf("mysql not config"))
	}
	mysqlConfig := config.Username + ":" + config.Password + "@(" + config.Path + ")/" + config.Dbname + "?" + config.Config
	fmt.Println(mysqlConfig)
	if db, err := gorm.Open("mysql", mysqlConfig); err != nil {
		panic(fmt.Sprintf("mysql start is error %v", err))
	} else {
		db.DB().SetMaxOpenConns(config.MaxOpenConns)
		db.DB().SetMaxIdleConns(config.MaxIdleConns)
		db.LogMode(config.LogMode)
		g_db = db
	}
	QLog().Info("mysql start success")
}
func Db() *gorm.DB {
	return g_db
}
