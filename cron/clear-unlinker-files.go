package cron

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/service"
	"strconv"
)

func clearUnlinkedFiles() {
	var fileIDs []int64
	global.DB.Model(&model.File{}).
		Select("id").Find(&fileIDs)

	for i := range fileIDs {
		var count int64
		global.DB.Model(&model.RelatedParty{}).
			Where("file_ids like ?", "%"+strconv.FormatInt(fileIDs[i], 10)+"%").
			Count(&count)
		if count > 0 {
			continue
		}

		//如果在所有表中都没有找到这个引用，就可以删掉这个文件
		var param service.FileDelete
		param.ID = fileIDs[i]
		res := param.Delete()
		if res.Code != 0 {
			var param1 service.ErrorLogCreate
			param1.Detail = "删除文件失败，文件id为：" + strconv.FormatInt(fileIDs[i], 10) + "请手动删除"
			param1.Create()
		}
	}
}
