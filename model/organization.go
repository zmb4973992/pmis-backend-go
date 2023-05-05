package model

import (
	"github.com/yitter/idgenerator-go/idgen"
	"gorm.io/gorm"
	"pmis-backend-go/global"
)

type Organization struct {
	BasicModel
	//连接其他表的id，暂无

	//连接dictionary_item表的id，暂无

	//日期，暂无

	//数字(允许为0、nil)
	SuperiorSnowID *int64 //上级机构SnowID，引用自身
	//数字(不允许为0、nil，必须有值)
	Sequence int //部门在当前层级下的顺序
	//字符串(允许为nil)
	Remarks *string //备注
	//字符串(不允许为nil，必须有值)
	Name string //名称
	//LevelName string `json:"level_name"` //层级名称，如公司、事业部、部门等
}

// TableName 修改表名
func (o *Organization) TableName() string {
	return "organization"
}

func (o *Organization) BeforeDelete(tx *gorm.DB) error {
	if o.SnowID == 0 {
		return nil
	}

	//删除相关的子表记录
	//先find，再delete，可以激活相关的钩子函数
	var records []OrganizationAndUser
	err = tx.Where(&OrganizationAndUser{OrganizationSnowID: o.SnowID}).
		Find(&records).Delete(&records).Error
	if err != nil {
		return err
	}

	return nil
}

func generateOrganizations() error {
	organizations := []Organization{
		{Name: "北京公司", Sequence: 1},
		{Name: "水泥工程事业部", Sequence: 2},
		{Name: "水泥工程市场一部", Sequence: 3},
		{Name: "水泥工程市场二部", Sequence: 4},
		{Name: "技术中心", Sequence: 5},
		{Name: "项目管理部", Sequence: 6},
		{Name: "工程项目执行部", Sequence: 7},
		{Name: "水泥延伸业务部", Sequence: 8},
		{Name: "综合管理和法律部", Sequence: 9},
		{Name: "人力资源和海外机构事务部", Sequence: 10},
		{Name: "财务管理部", Sequence: 11},
		{Name: "党建和纪检审计部", Sequence: 12},
		{Name: "储运管理部", Sequence: 13},
		{Name: "进口部/航空技术部", Sequence: 14},
		{Name: "成套业务一部", Sequence: 15},
		{Name: "成套业务二部", Sequence: 16},
		{Name: "成套业务三部", Sequence: 17},
		{Name: "成套业务四部", Sequence: 18},
		{Name: "成套业务五部", Sequence: 19},
		{Name: "成套业务六部", Sequence: 20},
		{Name: "投资企业", Sequence: 21},
		{Name: "海外机构", Sequence: 22},
		{Name: "国内企业管理部", Sequence: 23},
	}
	for _, organization := range organizations {
		err = global.DB.Where(&Organization{Name: organization.Name}).
			Where(&Organization{Sequence: organization.Sequence}).
			Attrs(&Organization{BasicModel: BasicModel{SnowID: idgen.NextId()}}).
			FirstOrCreate(&organization).Error
		if err != nil {
			global.SugaredLogger.Errorln(err)
			return err
		}
		//添加上级机构id
		var superiorOrganization Organization
		if organization.Name == "北京公司" { //如果是北京公司，就不做处理
		} else if organization.Name == "水泥工程市场一部" || // 如果是水泥工程事业部
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
			err = global.DB.Model(&Organization{}).Where("name = ?", organization.Name).Update("superior_snow_id", superiorOrganization.SnowID).Error
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
			err = global.DB.Model(&Organization{}).Where("name = ?", organization.Name).Update("superior_snow_id", superiorOrganization.SnowID).Error
			if err != nil {
				global.SugaredLogger.Errorln(err)
				return err
			}
		}
	}
	return nil
}
