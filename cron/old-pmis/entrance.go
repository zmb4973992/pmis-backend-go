package old_pmis

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/service"
)

func ImportDataForCron() {
	var user model.User
	err := global.DB.
		Where("username = 'z0030975'").
		First(&user).Error
	if err != nil {
		global.SugaredLogger.Panicln(err)
	}

	err = ImportData(user.Id)
	if err != nil {
		param := service.ErrorLogCreate{Detail: err.Error()}
		param.Create()
	}
}

func ImportData(userId int64) error {
	err := connectToDatabase()
	if err != nil {
		return err
	}

	err = importDisassembly(userId)
	if err != nil {
		return err
	}

	err = importProgress(userId)
	if err != nil {
		return err
	}

	return nil
}
