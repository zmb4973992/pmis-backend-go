package old_pmis

import (
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"pmis-backend-go/global"
	"time"
)

func connectToDatabase() (err error) {
	global.DBForOldPmis, err = gorm.Open(
		sqlserver.Open(global.Config.DbForOldPmis.DSN), &gorm.Config{},
	)
	if err != nil {
		return err
	}

	//使用gorm标准格式，创建连接池
	sqlDB, _ := global.DB.DB()
	// Set Max Idle Connections 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// Set Max Open Connections 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(100)
	// Set Connection Max Lifetime 设置了连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Hour)

	return nil
}
