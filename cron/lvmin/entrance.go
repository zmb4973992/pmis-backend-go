package lvmin

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/service"
)

func ImportDataForCron() {
	var user model.User
	err := global.DB.Where("username = 'z0030975'").
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
	err := ConnectToDatabase()
	if err != nil {
		return err
	}

	err = ImportRelatedParty(userId)
	if err != nil {
		return err
	}

	err = ImportProject(userId)
	if err != nil {
		return err
	}

	err = UpdateExchangeRageOfProject(userId)
	if err != nil {
		return err
	}

	err = ImportContract(userId)
	if err != nil {
		return err
	}

	err = UpdateExchangeRageOfContract(userId)
	if err != nil {
		return err
	}

	err = ImportActualExpenditure(userId)
	if err != nil {
		return err
	}

	err = ImportForecastedExpenditure(userId)
	if err != nil {
		return err
	}

	err = ImportPlannedExpenditure(userId)
	if err != nil {
		return err
	}

	err = ImportActualIncome(userId)
	if err != nil {
		return err
	}

	return nil
}
