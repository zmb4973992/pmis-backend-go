package util

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
)

// GetSuperiorIds 给定拆解id，找到所有上级id
func GetSuperiorIds(disassemblyId int64) (superiorIds []int64) {
	//superior_id可能为空，所以用指针来接收
	var disassembly model.Disassembly
	err := global.DB.Where("id = ?", disassemblyId).
		First(&disassembly).Error

	//如果发生任何错误、或者上级id为空：
	if err != nil || disassembly.SuperiorId == nil {
		return nil
	}

	superiorIds = append(superiorIds, *disassembly.SuperiorId)
	res := GetSuperiorIds(*disassembly.SuperiorId)

	superiorIds = append(superiorIds, res...)

	return superiorIds
}

// GetInferiorIds 给定拆解id，找到所有下级id
func GetInferiorIds(disassemblyId int64) (inferiorIds []int64) {
	var disassemblies []model.Disassembly
	err := global.DB.Where("superior_id = ?", disassemblyId).
		Find(&disassemblies).Error

	if err != nil {
		return nil
	}

	for i := range disassemblies {
		inferiorIds = append(inferiorIds, disassemblies[i].Id)
		res := GetInferiorIds(disassemblies[i].Id)

		inferiorIds = append(inferiorIds, res...)
	}

	return inferiorIds
}
