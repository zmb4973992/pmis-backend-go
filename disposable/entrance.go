package disposable

import (
	"pmis-backend-go/cron/lvmin"
	oldPmis "pmis-backend-go/cron/old-pmis"
	windowsAd "pmis-backend-go/cron/windows-ad"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/service"
)

func Init() {
	err := windowsAd.UpdateUsersByLDAP()
	if err != nil {
		param := service.ErrorLogCreate{Detail: err.Error()}
		param.Create()
		global.SugaredLogger.Panicln(err)
	}

	var user model.User
	err = global.DB.Where("username = 'z0030975'").
		First(&user).Error
	if err != nil {
		param := service.ErrorLogCreate{Detail: "找不到用户名为“z0030975”的用户"}
		param.Create()
		global.SugaredLogger.Panicln(err)
	}

	err = lvmin.ImportData(user.Id)
	if err != nil {
		param := service.ErrorLogCreate{Detail: err.Error()}
		param.Create()
		global.SugaredLogger.Panicln(err)
	}

	err = oldPmis.ImportData(user.Id)
	if err != nil {
		param := service.ErrorLogCreate{Detail: err.Error()}
		param.Create()
		global.SugaredLogger.Panicln(err)
	}
}
