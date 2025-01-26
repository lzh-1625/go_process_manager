package repository

import (
	"log"
	"os"
	"time"

	"github.com/lzh-1625/go_process_manager/internal/app/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db            *gorm.DB
	defaultConfig = gorm.Session{PrepareStmt: true, SkipDefaultTransaction: true}
)

func InitDb() {
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
		log.Panicf("sqlite数据库初始化失败！\n错误原因：%v", err)
	}
	sqlDB, err := gdb.DB()
	if err != nil {
		log.Panicf("sqlite数据库初始化失败！\n错误原因：%v", err)
	}
	sqlDB.SetConnMaxLifetime(time.Hour)
	db = gdb.Session(&defaultConfig)
	// db = gdb.Session(&defaultConfig).Debug()
	db.AutoMigrate(&model.Process{}, &model.User{}, &model.Permission{}, &model.Push{}, &model.Config{}, &model.ProcessLog{}, &model.Task{})
}
