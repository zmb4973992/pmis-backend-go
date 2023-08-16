package model

import (
	"pmis-backend-go/global"
)

type DataAuthority struct {
	BasicModel
	Name string //名称
	Sort int    //排序
}

// TableName 修改表名
func (d *DataAuthority) TableName() string {
	return "data_authority"
}

func generateDataAuthorities() error {
	records := []DataAuthority{
		{Name: "所有部门", Sort: 1},
		{Name: "所属部门和子部门", Sort: 2},
		{Name: "所属部门", Sort: 3},
		{Name: "无权限", Sort: 4},
	}
	for i := range records {
		err = global.DB.Where("name = ?", records[i].Name).
			FirstOrCreate(&records[i]).Error
		if err != nil {
			return err
		}
	}
	return nil
}
