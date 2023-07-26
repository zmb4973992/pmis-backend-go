package cron

import (
	"fmt"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/service"
	"strconv"
	"strings"
)

func importActualIncome() error {
	err := ImportActualIncomeFromTabShouKuan()
	if err != nil {
		return err
	}

	err = ImportActualIncomeFromTabShouHui()
	if err != nil {
		return err
	}

	return nil
}

type tabShouKuan struct {
	Date                     string  `gorm:"column:F10565"`
	Amount                   float64 `gorm:"column:F10617"`
	Currency                 string  `gorm:"column:F10576"`
	ContractCode             string  `gorm:"column:F10624"`
	ProjectCode              string  `gorm:"column:F10589"`
	ImportedRelatedPartyName string  `gorm:"column:F10581"`
	ImportedID               string  `gorm:"column:ID"`
	IOrd                     string  `gorm:"column:iOrd"`
}

func ImportActualIncomeFromTabShouKuan() error {
	fmt.Println("★★★★★开始处理人民币实际收款记录......★★★★★")

	var records []tabShouKuan
	global.DB2.Table("tabShouKuan").Where("F10617 > 0").
		Find(&records)

	for i := range records {
		if i > 0 && i%100 == 0 {
			process, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(i)/float64(len(records))), 64)
			fmt.Println("已处理", i, "条实际收款记录，当前进度：", fmt.Sprintf("%.0f", process*100), "%")
		}

		var tempCount int64
		global.DB.Model(&model.IncomeAndExpenditure{}).
			Where("imported_approval_id = ?", records[i].ImportedID+records[i].IOrd).
			Count(&tempCount)
		if tempCount == 0 {
			if records[i].Date != "" {
				records[i].Date = records[i].Date[:10]
			}

			var currencyID int64
			if records[i].Currency != "CNY" {
				param := service.ErrorLogCreate{
					Detail: "tabShouKuan视图的记录中发现无法匹配的币种：" +
						records[i].Currency + "，ID为：" + records[i].ImportedID,
				}
				param.Create()
			}

			var dictionaryTypeID int64
			err := global.DB.Model(&model.DictionaryType{}).
				Where("name = ?", "币种").Select("id").
				First(&dictionaryTypeID).Error
			if err != nil {
				param := service.ErrorLogCreate{
					Detail: "dictionaryType表中找不到”币种“这个名称",
				}
				param.Create()
			} else {
				err = global.DB.Model(&model.DictionaryDetail{}).
					Where("dictionary_type_id = ?", dictionaryTypeID).
					Where("name = ?", "人民币").Select("id").
					First(&currencyID).Error
				if err != nil {
					param := service.ErrorLogCreate{
						Detail: "dictionaryDetail表中找不到“人民币”这个名称",
					}
					param.Create()
				}
			}

			var contractID int64
			if records[i].ContractCode != "" {
				err = global.DB.Model(&model.Contract{}).
					Where("code = ?", records[i].ContractCode).Select("id").
					First(&contractID).Error
				if err != nil {
					param := service.ErrorLogCreate{
						Detail: "tabShouKuan视图的记录中发现无法匹配的合同编号：" +
							records[i].ContractCode,
					}
					param.Create()
					records[i].ContractCode = ""
				}
			}

			var projectID int64
			if records[i].ProjectCode != "" {
				err = global.DB.Model(&model.Project{}).
					Where("code = ?", records[i].ProjectCode).Select("id").
					First(&projectID).Error
				if err != nil {
					param := service.ErrorLogCreate{
						Detail: "tabShouKuan视图的记录中发现无法匹配的项目编号：" +
							records[i].ProjectCode,
					}
					param.Create()
					records[i].ProjectCode = ""
				}
			}

			var relatedPartyID int64
			if records[i].ImportedRelatedPartyName != "" {
				err = global.DB.Model(&model.RelatedParty{}).
					Where("name = ?", strings.TrimSpace(records[i].ImportedRelatedPartyName)).Select("id").
					First(&relatedPartyID).Error
				if err != nil {
					err = global.DB.Model(&model.RelatedParty{}).
						Where("imported_original_name like ?", "%"+strings.TrimSpace(records[i].ImportedRelatedPartyName)+"%").
						Select("id").
						First(&relatedPartyID).Error
					if err != nil {
						param := service.ErrorLogCreate{
							Detail: "tabShouKuan视图的记录中发现无法匹配的相关方名称：" +
								records[i].ImportedRelatedPartyName,
						}
						param.Create()
						records[i].ImportedRelatedPartyName = ""
					}
				}
			}

			exchangeRate := 1.0

			newRecord := service.IncomeAndExpenditureCreate{
				ProjectID:          projectID,
				ContractID:         contractID,
				FundDirection:      "收款",
				Currency:           currencyID,
				Kind:               "实际",
				Date:               records[i].Date,
				Amount:             &records[i].Amount,
				ExchangeRate:       &exchangeRate,
				Type:               0,
				Term:               0,
				Remarks:            "",
				Attachment:         "",
				ImportedApprovalID: records[i].ImportedID + records[i].IOrd,
			}

			res := newRecord.Create()

			if res.Code != 0 {
				param := service.ErrorLogCreate{
					Detail: "导入tabShouKuan视图的记录时发生错误：" +
						res.Message + "，ID为：" + records[i].ImportedID + "，iOrd为：" + records[i].IOrd,
				}
				param.Create()
			}
		}
	}

	return nil
}

type tabShouHui struct {
	Date                     string  `gorm:"column:F14395"`
	Amount                   float64 `gorm:"column:F14166"`
	Currency                 string  `gorm:"column:F14168"`
	ProjectCode              string  `gorm:"column:F16856"`
	ImportedRelatedPartyName string  `gorm:"column:F14394"`
	BankSerialID             string  `gorm:"column:F14165"`
	IOrd                     string  `gorm:"column:iOrd"`
	ExchangeRate             float64
}

func ImportActualIncomeFromTabShouHui() error {
	fmt.Println("★★★★★开始处理实际收汇记录......★★★★★")

	var records []tabShouHui
	global.DB2.Table("tabShouHui").Where("F14166 > 0").
		Find(&records)

	for i := range records {
		if i > 0 && i%100 == 0 {
			process, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(i)/float64(len(records))), 64)
			fmt.Println("已处理", i, "条实际收汇记录，当前进度：", fmt.Sprintf("%.0f", process*100), "%")
		}

		var tempCount int64
		global.DB.Model(&model.IncomeAndExpenditure{}).
			Where("imported_approval_id = ?", records[i].BankSerialID+records[i].IOrd).
			Count(&tempCount)
		if tempCount == 0 {
			if records[i].Date != "" {
				records[i].Date = records[i].Date[:10]
			}

			var currencyID int64
			if records[i].Currency != "" {
				switch records[i].Currency {
				case "RMB":
					records[i].Currency = "人民币"
				case "CNY":
					records[i].Currency = "人民币"
				case "1":
					records[i].Currency = "人民币"
				case "2":
					records[i].Currency = "美元"
				case "USD":
					records[i].Currency = "美元"
				case "3":
					records[i].Currency = "欧元"
				case "EUR":
					records[i].Currency = "欧元"
				}

				var dictionaryTypeID int64
				err := global.DB.Model(&model.DictionaryType{}).
					Where("name = ?", "币种").Select("id").
					First(&dictionaryTypeID).Error
				if err != nil {
					param := service.ErrorLogCreate{
						Detail: "dictionaryType表中找不到”币种“这个名称",
					}
					param.Create()
					records[i].Currency = ""
				} else {
					err = global.DB.Model(&model.DictionaryDetail{}).
						Where("name = ?", records[i].Currency).Select("id").
						First(&currencyID).Error
					if err != nil {
						param := service.ErrorLogCreate{
							Detail: "tabShouHui视图的记录中发现无法匹配的币种：" +
								records[i].Currency + "，银行流水号为：" + records[i].BankSerialID,
						}
						param.Create()
						records[i].Currency = ""
					}
				}
			}

			var projectID int64
			if records[i].ProjectCode != "" {
				err := global.DB.Model(&model.Project{}).
					Where("code = ?", records[i].ProjectCode).Select("id").
					First(&projectID).Error
				if err != nil {
					param := service.ErrorLogCreate{
						Detail: "tabShouHui视图的记录中发现无法匹配的项目编号：" +
							records[i].ProjectCode,
					}
					param.Create()
					records[i].ProjectCode = ""
				}
			}

			var relatedPartyID int64
			if records[i].ImportedRelatedPartyName != "" {
				err := global.DB.Model(&model.RelatedParty{}).
					Where("name = ?", strings.TrimSpace(records[i].ImportedRelatedPartyName)).Select("id").
					First(&relatedPartyID).Error
				if err != nil {
					err = global.DB.Model(&model.RelatedParty{}).
						Where("imported_original_name like ?", "%"+strings.TrimSpace(records[i].ImportedRelatedPartyName)+"%").
						Select("id").
						First(&relatedPartyID).Error
					if err != nil {
						param := service.ErrorLogCreate{
							Detail: "tabShouHui视图的记录中发现无法匹配的相关方名称：" +
								records[i].ImportedRelatedPartyName,
						}
						param.Create()
						records[i].ImportedRelatedPartyName = ""
					}
				}
			}

			switch records[i].Currency {
			case "人民币":
				records[i].ExchangeRate = 1
			case "美元":
				records[i].ExchangeRate = 7.2
			case "欧元":
				records[i].ExchangeRate = 7.8
			case "港币":
				records[i].ExchangeRate = 0.9
			case "新加坡元":
				records[i].ExchangeRate = 5.2
			case "马来西亚币":
				records[i].ExchangeRate = 1.5
			default:
				records[i].ExchangeRate = 1
			}

			newRecord := service.IncomeAndExpenditureCreate{
				ProjectID:          projectID,
				ContractID:         0,
				FundDirection:      "收款",
				Currency:           currencyID,
				Kind:               "实际",
				Date:               records[i].Date,
				Amount:             &records[i].Amount,
				ExchangeRate:       &records[i].ExchangeRate,
				Type:               0,
				Term:               0,
				Remarks:            "",
				Attachment:         "",
				ImportedApprovalID: records[i].BankSerialID + records[i].IOrd,
			}

			res := newRecord.Create()

			if res.Code != 0 {
				param := service.ErrorLogCreate{
					Detail: "导入tabShouHui视图的记录时发生错误：" +
						res.Message + "，银行流水ID为：" + records[i].BankSerialID + "，iOrd为：" + records[i].IOrd,
				}
				param.Create()
			}
		}
	}

	return nil
}
