package service

import (
	"fmt"
	"github.com/gookit/goutil/arrutil"
	"github.com/shopspring/decimal"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
	"sort"
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

	//date1, err := time.Parse("2006-01-02", "2025-03-22")

	//err = test(195, date1, 40)
	//err = test1(195, 40)

	err = test3(273, 40)

	if err != nil {
		fmt.Println(err)
	}

	return response.Success()
}

// 给定拆解id、日期、进度类型，计算自身的进度
func test(disassemblyID int, date time.Time, progressType int) (err error) {
	//删除相关进度,防止产生重复数据
	global.DB.Where("disassembly_id = ?", disassemblyID).
		Where("date = ?", date).
		Where("type = ?", progressType).
		Delete(&model.Progress{})

	//获取下级拆解情况
	var subDisassembly []model.Disassembly
	err = global.DB.Where("superior_id = ?", disassemblyID).
		Find(&subDisassembly).Error
	if err != nil {
		return err
	}

	var sumOfProgress float64

	for i := range subDisassembly {
		//下级拆解id是否包含有效记录
		var count int64
		global.DB.Model(&model.Progress{}).
			Where("disassembly_id = ?", subDisassembly[i].ID).
			Where("type = ?", progressType).
			Where("date <= ?", date).
			Count(&count)

		var tempSubProgress float64 = 0

		if count > 0 {
			global.DB.Model(&model.Progress{}).
				Where("disassembly_id = ?", subDisassembly[i].ID).
				Where("type = ?", progressType).
				Where("date <= ?", date).
				Order("date desc").Select("value").
				Limit(1).Find(&tempSubProgress)
		}

		subProgress := decimal.NewFromFloat(tempSubProgress)
		subWeight := decimal.NewFromFloat(0)
		if subDisassembly[i].Weight != nil {
			subWeight = decimal.NewFromFloat(*subDisassembly[i].Weight)
		}

		res, _ := subProgress.Mul(subWeight).Float64()

		sumOfProgress += res
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
		Value:         &sumOfProgress,
		DataSource:    &dataSource,
	}

	err = global.DB.Create(&progress).Error
	if err != nil {
		return err
	}

	return nil
}

// 给定拆解id、进度类型，计算自身的进度
func test1(disassemblyID int, progressType int) (err error) {
	//找到"系统计算"的字典值
	var dataSource int
	err = global.DB.Model(&model.DictionaryItem{}).
		Where("name = '系统计算'").Select("id").First(&dataSource).Error
	if err != nil {
		return err
	}

	//删除相关进度,防止产生重复数据
	global.DB.Where("disassembly_id = ?", disassemblyID).
		Where("data_source = ?", dataSource).
		Where("type = ?", progressType).
		Delete(&model.Progress{})

	//获取下级拆解情况
	var subDisassembly []model.Disassembly
	err = global.DB.Where("superior_id = ?", disassemblyID).
		Find(&subDisassembly).Error
	if err != nil {
		return err
	}

	//获取日期数组
	var tempDates []string

	for i := range subDisassembly {
		var subDates []string
		global.DB.Model(&model.Progress{}).
			Where("disassembly_id = ?", subDisassembly[i].ID).
			Select("date").Find(&subDates)
		tempDates = append(tempDates, subDates...)
	}

	//这里的日期格式为2020-01-01T00:00:00Z，需要转成2020-01-01
	var tempDates1 []string
	for i := range tempDates {
		tempDates1 = append(tempDates1, tempDates[i][0:10])
	}

	//确保日期唯一
	dates := arrutil.Unique(tempDates1)

	//给日期排序，从小到大
	sort.Strings(dates)

	for j := range dates {
		date, err1 := time.Parse("2006-01-02", dates[j])
		if err1 != nil {
			return err1
		}

		err = test(disassemblyID, date, progressType)
		if err != nil {
			return err
		}
	}
	return nil
}

// 给定拆解id，找到所有上级id
func test2(disassemblyID int) (superiorIDs []int) {
	//superior_id可能为空，所以用指针来接收
	var disassembly model.Disassembly
	err := global.DB.Where("id = ?", disassemblyID).
		First(&disassembly).Error

	//如果发生任何错误、或者上级id为空：
	if err != nil || disassembly.SuperiorID == nil {
		return nil
	}

	superiorIDs = append(superiorIDs, *disassembly.SuperiorID)
	res := test2(*disassembly.SuperiorID)
	superiorIDs = append(superiorIDs, res...)

	//倒序排列，确保先找出来的值放后面。这样调用时，就按最底层到最高层的顺序进行
	fmt.Println("原始：", superiorIDs)
	fmt.Println("加工后：", reverseSlice(superiorIDs))

	return superiorIDs
}

// 给定拆解id、进度类型，计算所有上级的进度
func test3(disassemblyID int, progressType int) (err error) {
	superiorIDs := test2(disassemblyID)

	for i := range superiorIDs {
		err = test1(superiorIDs[i], progressType)
		if err != nil {
			return err
		}
	}

	return nil
}

func reverseSlice(param []int) []int {
	if len(param) == 0 {
		return param
	}
	return append(reverseSlice(param[1:]), param[0])
}
