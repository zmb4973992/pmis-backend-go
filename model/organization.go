package model

import (
	"gorm.io/gorm"
	"pmis-backend-go/global"
)

type Organization struct {
	BasicModel
	//连接其他表的id，暂无

	//连接dictionary_item表的id，暂无

	//日期，暂无

	//数字(允许为0、nil)
	SuperiorID *int64 //上级机构ID，引用自身
	//数字(不允许为0、nil，必须有值)
	Sort int //部门在当前层级下的顺序
	//字符串(允许为nil)
	Remarks *string //备注
	//字符串(不允许为nil，必须有值)
	Name string //名称
	//LevelName string `json:"level_name"` //层级名称，如公司、事业部、部门等
	IsValid bool //是否有效
}

// TableName 修改表名
func (o *Organization) TableName() string {
	return "organization"
}

func (o *Organization) BeforeDelete(tx *gorm.DB) error {
	if o.ID == 0 {
		return nil
	}

	//删除相关的子表记录
	//先find，再delete，可以激活相关的钩子函数
	var records []OrganizationAndUser
	err = tx.Where(&OrganizationAndUser{OrganizationID: o.ID}).
		Find(&records).Delete(&records).Error
	if err != nil {
		return err
	}

	return nil
}

func generateOrganizations() error {
	organizations := []Organization{
		{Name: "北京公司", Sort: 1, IsValid: false},
		{Name: "公司领导", Sort: 2, IsValid: false},
		{Name: "水泥工程事业部", Sort: 3, IsValid: true},
		{Name: "水泥工程市场一部", Sort: 4, IsValid: true},
		{Name: "水泥工程市场二部", Sort: 5, IsValid: true},
		{Name: "技术中心", Sort: 6, IsValid: true},
		{Name: "项目管理部", Sort: 7, IsValid: true},
		{Name: "工程项目执行部", Sort: 8, IsValid: true},
		{Name: "水泥延伸业务部", Sort: 9, IsValid: true},
		{Name: "综合管理和法律部", Sort: 10, IsValid: true},
		{Name: "人力资源和海外机构事务部", Sort: 11, IsValid: true},
		{Name: "财务管理部", Sort: 12, IsValid: true},
		{Name: "党建文宣部", Sort: 13, IsValid: true},
		{Name: "纪检审计部", Sort: 14, IsValid: true},
		{Name: "储运管理部", Sort: 15, IsValid: true},
		{Name: "进口部/航空技术部", Sort: 16, IsValid: true},
		{Name: "成套业务一部", Sort: 17, IsValid: true},
		{Name: "成套业务二部", Sort: 18, IsValid: true},
		{Name: "成套业务四部", Sort: 19, IsValid: true},
		{Name: "成套业务五部", Sort: 20, IsValid: true},
		{Name: "成套业务六部", Sort: 21, IsValid: true},
		{Name: "投资企业", Sort: 22, IsValid: true},
		{Name: "海外机构", Sort: 23, IsValid: true},
		{Name: "国内企业管理部", Sort: 24, IsValid: true},
	}
	for _, organization := range organizations {
		err = global.DB.Where(&Organization{Name: organization.Name}).
			Where(&Organization{Sort: organization.Sort}).
			FirstOrCreate(&organization).Error
		if err != nil {
			global.SugaredLogger.Errorln(err)
			return err
		}
		//添加上级机构id
		var superiorOrganization Organization
		if organization.Name == "北京公司" { //如果是北京公司，就不做处理
		} else if organization.Name == "水泥工程市场一部" || // 如果是水泥工程事业部下的部门
			organization.Name == "水泥工程市场二部" ||
			organization.Name == "技术中心" ||
			organization.Name == "项目管理部" ||
			organization.Name == "工程项目执行部" ||
			organization.Name == "水泥延伸业务部" {
			//查找上级部门的信息
			err = global.DB.Where("name = ?", "水泥工程事业部").First(&superiorOrganization).Error
			if err != nil {
				global.SugaredLogger.Errorln(err)
				return err
			}
			//把上级部门的id赋值给本部门
			err = global.DB.Model(&Organization{}).Where("name = ?", organization.Name).Update("superior_id", superiorOrganization.ID).Error
			if err != nil {
				global.SugaredLogger.Errorln(err)
				return err
			}
		} else { //如果非水泥工程事业部的其他部门
			//查找上级部门的信息
			err = global.DB.Where("name = ?", "北京公司").First(&superiorOrganization).Error
			if err != nil {
				global.SugaredLogger.Errorln(err)
				return err
			}
			//把上级部门的id赋值给本部门
			err = global.DB.Model(&Organization{}).Where("name = ?", organization.Name).Update("superior_id", superiorOrganization.ID).Error
			if err != nil {
				global.SugaredLogger.Errorln(err)
				return err
			}
		}
	}
	return nil
}
