package model

import (
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"learn-go/global"
	"time"
)

var (
	err error
)

func Init() {
	//通过gorm连接sqlserver数据库
	global.DB, err = gorm.Open(sqlserver.Open(global.Config.DBConfig.DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	//使用gorm标准格式，创建连接池
	sqlDB, _ := global.DB.DB()
	// Set Max Idle Connections 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// Set Max Open Connections 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(100)
	// Set Connection Max Lifetime 设置了连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Hour)
	err = global.DB.AutoMigrate(
		&RelatedParty{},
		&Project{},
		&Department{},
		&User{},
		&DepartmentAndUser{},
		&Role{},
		&RoleAndUser{},
		&Contract{},
		&ProjectBreakdown{},
		&WorkProgress{},
		&ActualReceiptAndPayment{},
		&PlannedReceiptAndPayment{},
		&PredictedReceiptAndPayment{},
		&Dictionary{},
		&ProjectAndUser{},
		&OperationHistory{},
		&Test{},
	)
	if err != nil {
		panic(err)
	}

	//生成初始数据
	generateData()

}

func generateData() {
	if err = generateRoles(); err != nil {
		panic(err)
	}
	if err = generateDepartments(); err != nil {
		panic(err)
	}

}
