package model

import (
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
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
		&DictionaryType{},                        //字典类型
		&DictionaryDetail{},                      //字典项详情
		&RelatedParty{},                          //相关方
		&Project{},                               //项目
		&Organization{},                          //部门
		&RoleAndOrganization{},                   //部门和数据范围的中间表
		&User{},                                  //用户
		&OrganizationAndUser{},                   //部门和用户的中间表
		&Role{},                                  //角色
		&UserAndRole{},                           //角色和用户的中间表
		&Contract{},                              //合同
		&Disassembly{},                           //项目拆解
		&Progress{},                              //工作进度
		&IncomeAndExpenditure{},                  //收付款
		&RequestLog{},                            //操作记录
		&ErrorLog{},                              //错误日志
		&CasbinRule{},                            //casbin规则
		&File{},                                  //上传的文件
		&Test{},                                  //测试
		&Menu{},                                  //菜单
		&RoleAndMenu{},                           //角色和菜单的中间表
		&Api{},                                   //api接口
		&MenuAndApi{},                            //菜单和api的中间表
		&ProjectDailyAndCumulativeIncome{},       //项目累计收款
		&ProjectDailyAndCumulativeExpenditure{},  //项目累计和当日付款
		&ContractDailyAndCumulativeIncome{},      //合同累计收款
		&ContractDailyAndCumulativeExpenditure{}, //合同累计付款
		&DataAuthority{},                         //数据范围表
		&UserAndDataAuthority{},                  //用户和数据范围的中间表
		&Message{},                               //消息表
		&MessageAndUser{},                        //消息和用户的中间表
		&OperationLog{},                          //用户操作日志
	)
	if err != nil {
		global.SugaredLogger.Panicln(err)
	}

	//生成初始数据
	generateInitialData()
}

func generateInitialData() {
	if err = generateDictionaryType(); err != nil {
		global.SugaredLogger.Panicln(err)
	}
	if err = generateDictionaryDetail(); err != nil {
		global.SugaredLogger.Panicln(err)
	}
	if err = generateOrganizations(); err != nil {
		global.SugaredLogger.Panicln(err)
	}
	if err = generateCasbinRules(); err != nil {
		global.SugaredLogger.Panicln(err)
	}
	if err = generateDataAuthorities(); err != nil {
		global.SugaredLogger.Panicln(err)
	}
}
