package service

import (
	"fmt"
	"github.com/shopspring/decimal"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
	"time"
)

//以下为入参

type ProgressGet struct {
	ID int
}

type ProgressCreate struct {
	Creator      int
	LastModifier int

	DisassemblyID int      `json:"disassembly_id" binding:"required"`
	Date          string   `json:"date" binding:"required"`
	Type          int      `json:"type" binding:"required"`
	Value         *float64 `json:"value" binding:"required"`
	Remark        string   `json:"remark,omitempty"`
}

//以下为出参

type ProgressOutput struct {
	Creator      *int `json:"creator"`
	LastModifier *int `json:"last_modifier"`
	ID           int  `json:"id"`

	DisassemblyID *int     `json:"disassembly_id"`
	Date          *string  `json:"date"`
	Type          *int     `json:"type"`
	Value         *float64 `json:"value"`
	Remark        *string  `json:"remark"`
	DataSource    *string  `json:"data_source"`
}

func (p *ProgressGet) Get() response.Common {
	var result ProgressOutput
	err := global.DB.Model(model.Progress{}).
		Where("id = ?", p.ID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorRecordNotFound)
	}

	//默认格式为这样的string：2019-11-01T00:00:00Z，需要取年月日(前9位)
	if result.Date != nil {
		temp := *result.Date
		*result.Date = temp[:10]
	}

	return response.SuccessWithData(result)
}

func (p *ProgressCreate) Create() response.Common {
	var paramOut model.Progress

	if p.Creator > 0 {
		paramOut.Creator = &p.Creator
	}

	if p.LastModifier > 0 {
		paramOut.LastModifier = &p.LastModifier
	}

	paramOut.DisassemblyID = &p.DisassemblyID

	date, err := time.Parse("2006-01-02", p.Date)
	if err != nil {
		return response.Failure(util.ErrorInvalidDateFormat)
	}
	paramOut.Date = &date

	paramOut.Type = &p.Type

	paramOut.Value = p.Value

	if p.Remark != "" {
		paramOut.Remark = &p.Remark
	}

	//计算有修改值的字段数，分别进行不同处理
	tempParamOut, err := util.StructToMap(paramOut)
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	paramOutForCounting := util.MapCopy(tempParamOut,
		"Creator", "LastModifier", "Deleter", "CreateAt", "UpdatedAt", "DeletedAt")

	if len(paramOutForCounting) == 0 {
		return response.Failure(util.ErrorFieldsToBeCreatedNotFound)
	}

	res := global.DB.FirstOrCreate(&paramOut, model.Progress{
		DisassemblyID: &p.DisassemblyID,
		Date:          &date,
		Type:          &p.Type,
		Value:         p.Value,
	})

	if res.Error != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToCreateRecord)
	}

	if res.RowsAffected == 0 {
		return response.Failure(util.ErrorDuplicateRecord)
	}

	date1, err := time.Parse("2006-01-02", "2025-03-22")

	err = test(5, date1, 38)
	if err != nil {
		fmt.Println(err)
	}

	return response.Success()
}

// 给定拆解id、日期、进度类型，计算自身的进度
func test(disassemblyID int, date time.Time, progressType int) (err error) {
	var progressValue float64

	//获取下级拆解情况
	var subDisassembly []model.Disassembly
	err = global.DB.Where("superior_id = ?", disassemblyID).
		Find(&subDisassembly).Error
	if err != nil {
		return err
	}

	for i := range subDisassembly {
		//下级拆解id是否包含有效记录
		var count int64
		global.DB.Model(&model.Progress{}).
			Where("disassembly_id = ?", subDisassembly[i].ID).
			Where("type = ?", progressType).
			Where("date <= ?", date).
			Count(&count)

		var tempSubProgressValue float64
		if count > 0 {
			err = global.DB.Model(&model.Progress{}).
				Where("disassembly_id = ?", subDisassembly[i].ID).
				Where("type = ?", progressType).
				Where("date <= ?", date).
				Order("date desc").Select("value").
				First(&tempSubProgressValue).Error
			if err != nil {
				return err
			}
		} else {
			tempSubProgressValue = 0
		}

		subProgressValue := decimal.NewFromFloat(tempSubProgressValue)
		weight := decimal.NewFromFloat(0)
		if subDisassembly[i].Weight != nil {
			weight = decimal.NewFromFloat(*subDisassembly[i].Weight)
		}
		res, _ := subProgressValue.Mul(weight).Float64()

		progressValue += res
	}

	//找到"系统计算"的字典值
	var dataSource int
	err = global.DB.Model(&model.DictionaryItem{}).
		Where("name = '系统计算'").Select("id").First(&dataSource).Error
	if err != nil {
		return err
	}

	var progress = model.Progress{
		DisassemblyID: &disassemblyID,
		Date:          &date,
		Type:          &progressType,
		Value:         &progressValue,
		DataSource:    &dataSource,
	}

	err = global.DB.Create(&progress).Error
	if err != nil {
		return err
	}

	return nil
}
