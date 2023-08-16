package lvmin

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/service"
)

func ImportDataForCron() {
	var user model.User
	err := global.DB.Where("username = 'z0030975'").First(&user).Error
	if err != nil {
		global.SugaredLogger.Panicln(err)
	}

	err = ImportData(user.ID)
	if err != nil {
		param := service.ErrorLogCreate{Detail: err.Error()}
		param.Create()
	}
}

func ImportData(userID int64) error {
	err := ConnectToDatabase()
	if err != nil {
		return err
	}

	err = ImportRelatedParty(userID)
	if err != nil {
		return err
	}

	err = ImportProject(userID)
	if err != nil {
		return err
	}

	err = UpdateExchangeRageOfProject(userID)
	if err != nil {
		return err
	}

	err = ImportContract(userID)
	if err != nil {
		return err
	}

	err = UpdateExchangeRageOfContract(userID)
	if err != nil {
		return err
	}

	err = ImportActualExpenditure(userID)
	if err != nil {
		return err
	}

	err = ImportForecastedExpenditure(userID)
	if err != nil {
		return err
	}

	err = ImportPlannedExpenditure(userID)
	if err != nil {
		return err
	}

	err = ImportActualIncome(userID)
	if err != nil {
		return err
	}

	return nil
}
