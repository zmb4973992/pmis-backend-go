package model

import "pmis-backend-go/global"

type Department struct {
	BaseModel
	Name       string `json:"name" binding:"required"`  //名称
	Level      string `json:"level" binding:"required"` //层级，如公司、事业部、部门等
	SuperiorID *int   `json:"superior_id"`              //上级机构ID，引用自身
	//这里是声名外键关系，并不是实际字段。结构体的字段名随意，首字母大写、否则不会导出，外键名会引用这个字段。
	//不建议用gorm的多对多的设定，不好修改

	//设置外键规则，SuperiorID作为外键，引用自身ID
	//数据库规则限制，自引用不能设置级联更新和级联删除
	//暂时不添加自引用的外键了，因为删除、更新都是麻烦事
	//SuperiorID1 []Department        `gorm:"foreignkey:SuperiorID"`
	//多对多的中间表需要外键，因为需要级联更新、级联删除
	Users    []DepartmentAndUser `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Projects []Project
}

// TableName 修改表名
func (Department) TableName() string {
	return "department"
}

func generateDepartments() error {
	departments := []Department{
		{Name: "北京公司", Level: "公司"},
		{Name: "水泥工程事业部", Level: "事业部"},
		{Name: "水泥工程市场一部", Level: "部门"},
		{Name: "水泥工程市场二部", Level: "部门"},
		{Name: "技术中心", Level: "部门"},
		{Name: "项目管理部", Level: "部门"},
		{Name: "工程项目执行部", Level: "部门"},
		{Name: "水泥延伸业务部", Level: "部门"},
		{Name: "综合管理和法律部", Level: "部门"},
		{Name: "人力资源和海外机构事务部", Level: "部门"},
		{Name: "财务管理部", Level: "部门"},
		{Name: "党建和纪检审计部", Level: "部门"},
		{Name: "储运管理部", Level: "部门"},
		{Name: "进口部/航空技术部", Level: "部门"},
		{Name: "成套业务一部", Level: "部门"},
		{Name: "成套业务二部", Level: "部门"},
		{Name: "成套业务三部", Level: "部门"},
		{Name: "成套业务四部", Level: "部门"},
		{Name: "成套业务五部", Level: "部门"},
		{Name: "成套业务六部", Level: "部门"},
		{Name: "投资企业", Level: "部门"},
		{Name: "海外机构", Level: "部门"},
		{Name: "国内企业管理部", Level: "部门"},
	}
	for _, department := range departments {
		err := global.DB.FirstOrCreate(&Department{}, department).Error
		if err != nil {
			return err
		}
		//添加上级机构id
		var superiorDepartment Department
		if department.Name == "北京公司" { //如果是北京公司，就不做处理
		} else if department.Name == "水泥工程市场一部" || // 如果是水泥工程事业部
			department.Name == "水泥工程市场二部" ||
			department.Name == "技术中心" ||
			department.Name == "项目管理部" ||
			department.Name == "工程项目执行部" ||
			department.Name == "水泥延伸业务部" {
			//查找上级部门的信息
			err = global.DB.Where("name = ?", "水泥工程事业部").First(&superiorDepartment).Error
			if err != nil {
				return err
			}
			//把上级部门的id赋值给本部门
			err = global.DB.Model(&Department{}).Where("name = ?", department.Name).Update("superior_id", superiorDepartment.ID).Error
			if err != nil {
				return err
			}
		} else { //如果非水泥工程事业部的其他部门
			//查找上级部门的信息
			err = global.DB.Where("name = ?", "北京公司").First(&superiorDepartment).Error
			if err != nil {
				return err
			}
			//把上级部门的id赋值给本部门
			err = global.DB.Model(&Department{}).Where("name = ?", department.Name).Update("superior_id", superiorDepartment.ID).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}
