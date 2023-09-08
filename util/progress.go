package util

import (
	"github.com/gookit/goutil/arrutil"
	"github.com/shopspring/decimal"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"sort"
	"time"
)

// UpdateProgressOfSuperiors 给定拆解id、进度类型，计算所有上级的进度
func UpdateProgressOfSuperiors(disassemblyID int64, progressType int64, userID int64) (err error) {
	superiorIDs := GetSuperiorIDs(disassemblyID)

	for i := range superiorIDs {
		err = UpdateOwnProgress(superiorIDs[i], progressType, userID)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdateOwnProgress 给定拆解id、进度类型，计算自身的进度
func UpdateOwnProgress(disassemblyID int64, progressType int64, userID int64) (err error) {
	//找到"进度的数据来源"的字典类型值
	var dataSourceOfProgress model.DictionaryType
	err = global.DB.Where("name = '进度的数据来源'").
		First(&dataSourceOfProgress).Error
	if err != nil {
		return err
	}

	//找到"系统计算"的字典详情值
	var systemCalculation model.DictionaryDetail
	err = global.DB.
		Where("dictionary_type_id = ?", dataSourceOfProgress.ID).
		Where("name = '系统计算'").
		First(&systemCalculation).Error
	if err != nil {
		return err
	}

	//删除相关进度,防止产生重复数据
	global.DB.Where("disassembly_id = ?", disassemblyID).
		Where("data_source = ?", systemCalculation.ID).
		Where("type = ?", progressType).
		Delete(&model.Progress{})

	//获取下级拆解情况
	var subDisassembly []model.Disassembly
	err = global.DB.Where("superior_id = ?", disassemblyID).
		Find(&subDisassembly).Error
	if err != nil {
		return err
	}

	//获取相应进度类型的日期数组
	var tempDates []string
	for i := range subDisassembly {
		var subDates []string
		global.DB.Model(&model.Progress{}).
			Where("disassembly_id = ?", subDisassembly[i].ID).
			Where("type = ?", progressType).
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

		err = updateOwnProgress1(disassemblyID, date, progressType, userID)
		if err != nil {
			return err
		}
	}
	return nil
}

// 给定拆解id、日期、进度类型，计算自身的进度
func updateOwnProgress1(disassemblyID int64, date time.Time, progressType int64, userID int64) (err error) {
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
	var dataSourceOfProgress model.DictionaryType
	err = global.DB.Where("name = '进度的数据来源'").
		First(&dataSourceOfProgress).Error
	if err != nil {
		return err
	}

	var systemCalculation model.DictionaryDetail
	err = global.DB.
		Where("dictionary_type_id = ?", dataSourceOfProgress.ID).
		Where("name = '系统计算'").
		First(&systemCalculation).Error
	if err != nil {
		return err
	}

	var progress = model.Progress{
		DisassemblyID: &disassemblyID,
		Date:          &date,
		Type:          &progressType,
		Value:         &sumOfProgress,
		DataSource:    &systemCalculation.ID,
		BasicModel:    model.BasicModel{Creator: &userID},
	}

	err = global.DB.Create(&progress).Error
	if err != nil {
		return err
	}

	return nil
}
