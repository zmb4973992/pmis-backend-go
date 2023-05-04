package util

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
)

// 给定拆解id，找到所有上级id
func getSuperiorSnowIDs(disassemblySnowID int64) (superiorSnowIDs []int64) {
	//superior_id可能为空，所以用指针来接收
	var disassembly model.Disassembly
	err := global.DB.Where("id = ?", disassemblySnowID).
		First(&disassembly).Error

	//如果发生任何错误、或者上级id为空：
	if err != nil || disassembly.SuperiorSnowID == nil {
		return nil
	}

	superiorSnowIDs = append(superiorSnowIDs, *disassembly.SuperiorSnowID)
	res := getSuperiorSnowIDs(*disassembly.SuperiorSnowID)

	superiorSnowIDs = append(superiorSnowIDs, res...)

	return superiorSnowIDs
}

// GetInferiorSnowIDs 给定拆解id，找到所有下级id
func GetInferiorSnowIDs(disassemblySnowID int64) (inferiorSnowIDs []int64) {
	var disassemblies []model.Disassembly
	err := global.DB.Where("superior_id = ?", disassemblySnowID).
		Find(&disassemblies).Error

	if err != nil {
		return nil
	}

	for i := range disassemblies {
		inferiorSnowIDs = append(inferiorSnowIDs, disassemblies[i].SnowID)
		res := GetInferiorSnowIDs(disassemblies[i].SnowID)

		inferiorSnowIDs = append(inferiorSnowIDs, res...)
	}

	return inferiorSnowIDs
}
