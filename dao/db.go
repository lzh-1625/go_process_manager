package dao

import (
	"log"
	zlog "msm/log"
	"msm/model"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

var defaultConfig = gorm.Session{PrepareStmt: true, SkipDefaultTransaction: true}

func init() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  true,
		},
	)
	gdb, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		zlog.Logger.Panicf("sqlite数据库初始化失败！\n错误原因：%v", err)
	}
	zlog.Logger.Info("sqlite初始化成功")
	db = gdb.Session(&defaultConfig)
	// db = gdb.Session(&defaultConfig).Debug()
	db.AutoMigrate(&model.Process{}, &model.User{}, &model.Permission{}, &model.Push{}, &model.Config{})
}
