package util

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
)

// 给定拆解id，找到所有上级id
func getSuperiorIDs(disassemblyID int64) (superiorIDs []int64) {
	//superior_id可能为空，所以用指针来接收
	var disassembly model.Disassembly
	err := global.DB.Where("id = ?", disassemblyID).
		First(&disassembly).Error

	//如果发生任何错误、或者上级id为空：
	if err != nil || disassembly.SuperiorID == nil {
		return nil
	}

	superiorIDs = append(superiorIDs, *disassembly.SuperiorID)
	res := getSuperiorIDs(*disassembly.SuperiorID)

	superiorIDs = append(superiorIDs, res...)

	return superiorIDs
}

// GetInferiorIDs 给定拆解id，找到所有下级id
func GetInferiorIDs(disassemblyID int64) (inferiorIDs []int64) {
	var disassemblies []model.Disassembly
	err := global.DB.Where("superior_id = ?", disassemblyID).
		Find(&disassemblies).Error

	if err != nil {
		return nil
	}

	for i := range disassemblies {
		inferiorIDs = append(inferiorIDs, disassemblies[i].ID)
		res := GetInferiorIDs(disassemblies[i].ID)

		inferiorIDs = append(inferiorIDs, res...)
	}

	return inferiorIDs
}
