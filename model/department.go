package model

import (
	"gorm.io/gorm"
	"pmis-backend-go/global"
	"time"
)

type Department struct {
	BaseModel
	Name       string  `json:"name"`        //名称
	LevelName  string  `json:"level_name"`  //层级名称，如公司、事业部、部门等
	Sequence   int     `json:"sequence"`    //部门在当前层级下的顺序
	SuperiorID *int    `json:"superior_id"` //上级机构ID，引用自身
	Remarks    *string `json:"remarks"`     //备注

	//这里是声名外键关系，并不是实际字段。结构体的字段名随意，首字母大写、否则不会导出，外键名会引用这个字段。
	//不建议用gorm的多对多的设定，不好修改
	//设置外键规则，SuperiorID作为外键，引用自身ID
	//数据库规则限制，自引用不能设置级联更新和级联删除
	//暂时不添加自引用的外键了，因为删除、更新都是麻烦事
	//多对多的中间表需要外键，因为需要级联更新、级联删除
}

// TableName 修改表名
func (*Department) TableName() string {
	return "department"
}

func (d *Department) BeforeDelete(tx *gorm.DB) error {
	if d.ID > 0 {
		//如果有删除人的id，则记录下来
		if d.Deleter != nil && *d.Deleter > 0 {
			err := tx.Model(&Department{}).Where("id = ?", d.ID).
				Update("deleter", d.Deleter).Error
			if err != nil {
				return err
			}
		}
		//删除相关的子表记录
		err = tx.Model(&DepartmentAndUser{}).Where("department_id = ?", d.ID).
			Updates(map[string]any{
				"deleted_at": time.Now(),
				"deleter":    d.Deleter,
			}).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func generateDepartments() error {
	departments := []Department{
		{Name: "北京公司", LevelName: "公司", Sequence: 1},
		{Name: "水泥工程事业部", LevelName: "事业部", Sequence: 2},
		{Name: "水泥工程市场一部", LevelName: "部门", Sequence: 3},
		{Name: "水泥工程市场二部", LevelName: "部门", Sequence: 4},
		{Name: "技术中心", LevelName: "部门", Sequence: 5},
		{Name: "项目管理部", LevelName: "部门", Sequence: 6},
		{Name: "工程项目执行部", LevelName: "部门", Sequence: 7},
		{Name: "水泥延伸业务部", LevelName: "部门", Sequence: 8},
		{Name: "综合管理和法律部", LevelName: "部门", Sequence: 9},
		{Name: "人力资源和海外机构事务部", LevelName: "部门", Sequence: 10},
		{Name: "财务管理部", LevelName: "部门", Sequence: 11},
		{Name: "党建和纪检审计部", LevelName: "部门", Sequence: 12},
		{Name: "储运管理部", LevelName: "部门", Sequence: 13},
		{Name: "进口部/航空技术部", LevelName: "部门", Sequence: 14},
		{Name: "成套业务一部", LevelName: "部门", Sequence: 15},
		{Name: "成套业务二部", LevelName: "部门", Sequence: 16},
		{Name: "成套业务三部", LevelName: "部门", Sequence: 17},
		{Name: "成套业务四部", LevelName: "部门", Sequence: 18},
		{Name: "成套业务五部", LevelName: "部门", Sequence: 19},
		{Name: "成套业务六部", LevelName: "部门", Sequence: 20},
		{Name: "投资企业", LevelName: "部门", Sequence: 21},
		{Name: "海外机构", LevelName: "部门", Sequence: 22},
		{Name: "国内企业管理部", LevelName: "部门", Sequence: 23},
	}
	for _, department := range departments {
		err := global.DB.FirstOrCreate(&Department{}, department).Error
		if err != nil {
			global.SugaredLogger.Errorln(err)
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
				global.SugaredLogger.Errorln(err)
				return err
			}
			//把上级部门的id赋值给本部门
			err = global.DB.Model(&Department{}).Where("name = ?", department.Name).Update("superior_id", superiorDepartment.ID).Error
			if err != nil {
				global.SugaredLogger.Errorln(err)
				return err
			}
		} else { //如果非水泥工程事业部的其他部门
			//查找上级部门的信息
			err = global.DB.Where("name = ?", "北京公司").First(&superiorDepartment).Error
			if err != nil {
				global.SugaredLogger.Errorln(err)
				return err
			}
			//把上级部门的id赋值给本部门
			err = global.DB.Model(&Department{}).Where("name = ?", department.Name).Update("superior_id", superiorDepartment.ID).Error
			if err != nil {
				global.SugaredLogger.Errorln(err)
				return err
			}
		}
	}
	return nil
}
