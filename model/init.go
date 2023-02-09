package model

import (
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"os"
	"pmis-backend-go/global"
	"time"
)

var err error

func InitDatabase() {
	//通过gorm连接sqlserver数据库
	global.DB, err = gorm.Open(sqlserver.Open(global.Config.DBConfig.DSN), &gorm.Config{})
	if err != nil {
		global.SugaredLogger.Panicln(err)
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
		&RelatedParty{},               //相关方
		&Project{},                    //项目
		&Department{},                 //部门
		&User{},                       //用户
		&DepartmentAndUser{},          //部门和用户的中间表
		&Role{},                       //角色
		&RoleAndUser{},                //角色和用户的中间表
		&Contract{},                   //合同
		&Disassembly{},                //项目拆解
		&WorkProgress{},               //工作进度
		&ActualReceiptAndPayment{},    //实际收付款
		&PlannedReceiptAndPayment{},   //计划收付款
		&PredictedReceiptAndPayment{}, //预测收付款
		&DictionaryType{},             //字典类型
		&DictionaryItem{},             //字典项详情
		&ProjectAndUser{},             //项目和用户的中间表
		&OperationLog{},               //操作记录
		&ErrorLog{},                   //错误日志
		&WorkNote{},                   //工作备注
		&WorkReview{},                 //工作点评
		&CasbinRule{},                 //casbin规则
		&File{},                       //上传的文件
		&Test{},                       //测试
	)
	if err != nil {
		global.SugaredLogger.Panicln(err)
	}

	//创建所需的视图
	createView()

	//生成初始数据
	generateInitialData()
}

func createView() {
	sqlStatement, err := os.ReadFile("./sql/create_view.sql")
	if err != nil {
		global.SugaredLogger.Panicln(err)
	}
	err = global.DB.Exec(string(sqlStatement)).Error
	if err != nil {
		global.SugaredLogger.Panicln(err)
	}
}

func generateInitialData() {
	if err = generateRoles(); err != nil {
		global.SugaredLogger.Panicln(err)
	}
	if err = generateDepartments(); err != nil {
		global.SugaredLogger.Panicln(err)
	}
	if err = generateDictionaryTypes(); err != nil {
		global.SugaredLogger.Panicln(err)
	}
	if err = generateDictionaryItems(); err != nil {
		global.SugaredLogger.Panicln(err)
	}
	if err = generateCasbinRules(); err != nil {
		global.SugaredLogger.Panicln(err)
	}
}
