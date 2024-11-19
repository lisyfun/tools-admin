package db

import (
	"fmt"
	"time"
	"tools-admin/backend/common/config"
	"tools-admin/backend/model"
	"tools-admin/backend/pkg/log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

func init() {
	var err error
	db := config.Config.Db
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		db.Username, db.Password, db.Host, db.Port, db.Database)
	Db, err = gorm.Open(mysql.Open(url), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}
	if Db.Error != nil {
		panic(Db.Error)
	}

	sqlDB, err := Db.DB()
	if err != nil {
		panic(err)
	}

	// 设置连接池配置
	sqlDB.SetMaxIdleConns(db.MaxIdleConns)
	sqlDB.SetMaxOpenConns(db.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(db.ConnMaxLifetime) * time.Second)

	// 自动迁移数据库
	err = Db.AutoMigrate(
		&model.User{},
		&model.Menu{},
		&model.Task{},
		&model.TaskLog{},
	)
	if err != nil {
		log.Error("数据库迁移失败: " + err.Error())
		panic(err)
	}

	log.Info("数据库初始化成功!")
}
