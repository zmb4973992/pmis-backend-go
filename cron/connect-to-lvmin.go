package cron

import (
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"pmis-backend-go/global"
	"pmis-backend-go/service"
	"time"
)

// 连接率敏的数据库
func connectToLvmin() {
	var err error
	global.DB2, err = gorm.Open(
		sqlserver.Open(global.Config.DB2Config.DSN), &gorm.Config{},
	)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		param := service.ErrorLogCreate{
			Detail: err.Error(),
			Date:   time.Now().Format("2006-01-02"),
		}
		param.Create()
		return
	}

	//使用gorm标准格式，创建连接池
	sqlDB2, _ := global.DB.DB()
	// Set Max Idle Connections 设置空闲连接池中连接的最大数量
	sqlDB2.SetMaxIdleConns(10)
	// Set Max Open Connections 设置打开数据库连接的最大数量
	sqlDB2.SetMaxOpenConns(100)
	// Set Connection Max Lifetime 设置了连接可复用的最大时间
	sqlDB2.SetConnMaxLifetime(time.Hour)

	err = UpdateProjectCumulativeIncomeAndExpenditure()
	if err != nil {
		global.SugaredLogger.Errorln(err)
		param := service.ErrorLogCreate{
			Detail: err.Error(),
			Date:   time.Now().Format("2006-01-02"),
		}
		param.Create()
		return
	}

	err = UpdateContractCumulativeIncomeAndExpenditure()
	if err != nil {
		global.SugaredLogger.Errorln(err)
		param := service.ErrorLogCreate{
			Detail: err.Error(),
			Date:   time.Now().Format("2006-01-02"),
		}
		param.Create()
		return
	}

	err = importDataFromLvmin()
	if err != nil {
		global.SugaredLogger.Errorln(err)
		param := service.ErrorLogCreate{
			Detail: err.Error(),
			Date:   time.Now().Format("2006-01-02"),
		}
		param.Create()
		return
	}
}
