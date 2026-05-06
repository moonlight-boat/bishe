package database

import (
	"job-backend/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase(databasePath string) error {
	var err error
	
	// 连接SQLite数据库
	DB, err = gorm.Open(sqlite.Open(databasePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	
	if err != nil {
		return err
	}
	
	// 自动迁移数据库表
	err = DB.AutoMigrate(&models.JobPosting{})
	if err != nil {
		return err
	}
	
	log.Println("数据库连接成功，表结构迁移完成")
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// CloseDatabase 关闭数据库连接
func CloseDatabase() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
