package cron

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/service"
	"pmis-backend-go/util"
	"strconv"
)

func clearUnlinkedFiles() {
	var fileIds []int64
	global.DB.Model(&model.File{}).
		Select("id").Find(&fileIds)

	for i := range fileIds {
		var count int64
		global.DB.Model(&model.RelatedParty{}).
			Where("file_ids like ?", "%"+strconv.FormatInt(fileIds[i], 10)+"%").
			Count(&count)
		if count > 0 {
			continue
		}

		//如果在所有表中都没有找到这个引用，就可以删掉这个文件
		var param service.FileDelete
		param.Id = fileIds[i]
		errCode := param.Delete()
		if errCode != util.Success {
			var param1 service.ErrorLogCreate
			param1.Detail = "删除文件失败，文件id为：" + strconv.FormatInt(fileIds[i], 10) + "请手动删除"
			param1.Create()
		}
	}
}
