package util

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
)

// FieldIsInModel 验证提交的单个字段是否存在于表中（即数据表是否有相应的字段）
func FieldIsInModel(model model.IModel, field string) bool {
	//获取自定义的数据库表名
	tableName := model.TableName()
	var existedFields []string
	//自行拼接的sql，找出对应表名的所有字段名
	//sqlStatement server的标准写法为：Select Name FROM SysColumns Where id = Object_Id('[某某表]')
	//给 某某表 加上中括号，是因为当表名中含有特殊字符时，直接使用单引号，会出现表名不被识别的问题
	//这里goland编译器莫名报错，函数可以正常运行，可忽略
	sqlStatement := "Select Name FROM SysColumns Where id = OBJECT_ID('[" + tableName + "]')"
	global.DB.Raw(sqlStatement).Find(&existedFields)
	//如果表中字段数量>0且该字段在表的这些字段中
	if len(existedFields) > 0 && IsInSlice(field, existedFields) {
		return true
	}
	return false
}

// FieldsAreInModel 验证提交的多个字段是否存在于表中（即数据表是否有相应的字段）
//func FieldsAreInModel(model model.IModel, fields ...string) bool {
//	for _, field := range fields {
//		res := FieldIsInModel(model, field)
//		//如果有任何一个字段不符合要求,则直接返回false
//		if res == false {
//			return false
//		}
//	}
//	return true
//}

func GetNumberOfPages(numberOfRecords int, pageSize int) (numberOfPages int) {
	//如果没有记录数，或者单页条数小于零：
	if numberOfRecords <= 0 || pageSize < 0 {
		return 0
	}
	//如果有记录数、但单页条数为零（不分页）：
	if numberOfRecords > 0 && pageSize == 0 {
		return 1
	}

	numberOfPages = numberOfRecords / pageSize
	if numberOfRecords%pageSize != 0 {
		numberOfPages++
	}
	return numberOfPages
}
