package dao

import (
	"learn-go/dto"
	"learn-go/global"
	"learn-go/model"
)

type departmentDAO struct{}

func (departmentDAO) Get(departmentID int) *dto.DepartmentGetDTO {
	var department model.Department

	err := global.DB.Where("id = ?", departmentID).First(&department).Error
	if err != nil {
		return nil
	}

	var paramOut dto.DepartmentGetDTO
	if department.Name != "" {
		paramOut.Name = department.Name
	}

	if department.Level != "" {
		paramOut.Level = department.Level
	}

	if department.SuperiorID != nil {
		paramOut.SuperiorID = department.SuperiorID
	}

	return &paramOut

	//默认嵌套递归次数上限为4次，太多了降低效率，而且没必要
	//return getDepartmentWithRecursionLimit(departmentID, 4, 0)
}

//由于get方法有递归调用，所以需要在这里多加2个参数进行限制。标准的get方法调用这个内部函数，达到封装的效果
//func getDepartmentWithRecursionLimit(departmentID int, recursionTimesLimit int, recursionTimes int) *dto.DepartmentGetDTO {
//	var departmentGetDTO = dto.DepartmentGetDTO{}
//	//把基础的部门信息查出来
//	var department model.Department
//	err := global.DB.Where("id = ?", departmentID).First(&department).Error
//	if err != nil {
//		return nil
//	}
//	//把所有查出的结果赋值给输出变量
//	departmentGetDTO.Name = department.Name
//	departmentGetDTO.Level = department.Level
//
//	//递归查询上级部门信息
//	if department.SuperiorID != nil {
//		recursionTimes += 1
//		if recursionTimes <= recursionTimesLimit {
//			departmentGetDTO.SuperiorID = getDepartmentWithRecursionLimit(*department.SuperiorID, recursionTimesLimit, recursionTimes)
//		} else {
//			departmentGetDTO.SuperiorID = "递归深度超过" + strconv.Itoa(recursionTimesLimit) + "次，可能存在循环递归，请检查数据是否正确"
//		}
//	}
//	return &departmentGetDTO
//}

// Create 这里是只负责新增，不写任何业务逻辑。只要收到参数就创建数据库记录，然后返回错误
func (departmentDAO) Create(param *model.Department) error {
	err := global.DB.Create(param).Error
	return err
}

// Update 这里是只负责更新，不写任何业务逻辑。只要收到id和更新参数，然后返回错误
func (departmentDAO) Update(param *model.Department) error {
	//注意，这里就算没有找到记录，也不会报错，只有更新字段出现问题才会报错。详见gorm的update用法
	err := global.DB.Where("id = ?", param.ID).Omit("created_at", "creator").Save(param).Error
	return err
}

func (departmentDAO) Delete(departmentID int) error {
	//注意，这里就算没有找到记录，也不会报错。详见gorm的delete用法
	err := global.DB.Delete(&model.Department{}, departmentID).Error
	return err
}
