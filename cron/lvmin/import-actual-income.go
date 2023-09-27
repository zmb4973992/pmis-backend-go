package lvmin

import (
	"fmt"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/service"
	"pmis-backend-go/util"
	"strconv"
	"strings"
)

func ImportActualIncome(userId int64) error {
	err := ImportActualIncomeFromTabShouKuan(userId)
	if err != nil {
		return err
	}

	err = ImportActualIncomeFromTabShouHui(userId)
	if err != nil {
		return err
	}

	err = ImportActualIncomeFromTabShouPiao(userId)
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
	ImportedId               string  `gorm:"column:ContractId"`
	IOrd                     string  `gorm:"column:iOrd"`
}

func ImportActualIncomeFromTabShouKuan(userId int64) error {
	fmt.Println("★★★★★开始处理人民币实际收款记录......★★★★★")

	var records []tabShouKuan
	global.DBForLvmin.Table("tabShouKuan").
		Where("F10617 > 0").
		Find(&records)

	var currency model.DictionaryType
	err := global.DB.
		Where("name = ?", "币种").
		First(&currency).Error
	if err != nil {
		param := service.ErrorLogCreate{
			Detail: "dictionaryType表中找不到”币种“这个名称",
		}
		param.Create()
		return err
	}

	var fundDirection model.DictionaryType
	var income model.DictionaryDetail
	err = global.DB.
		Where("name = '收付款的资金方向'").
		First(&fundDirection).Error
	if err != nil {
		param := service.ErrorLogCreate{
			Detail: "在dictionary_type表中找不到”收付款的资金方向“这个名称",
		}
		param.Create()
		return err
	}

	err = global.DB.
		Where("dictionary_type_id =?", fundDirection.Id).
		Where("name = '收款'").
		First(&income).Error
	if err != nil {
		param := service.ErrorLogCreate{
			Detail: "在dictionary_detail表中找不到”收款“这个名称",
		}
		param.Create()
		return err
	}

	var kind model.DictionaryType
	var actual model.DictionaryDetail
	err = global.DB.
		Where("name = '收付款的种类'").
		First(&kind).Error
	if err != nil {
		param := service.ErrorLogCreate{
			Detail: "在dictionary_type表中找不到”收付款的种类“这个名称",
		}
		param.Create()
		return err
	}

	err = global.DB.
		Where("dictionary_type_id =?", kind.Id).
		Where("name = '实际'").
		First(&actual).Error
	if err != nil {
		param := service.ErrorLogCreate{
			Detail: "在dictionary_detail表中找不到”实际“这个名称(dictionaryType为：收付款的种类)",
		}
		param.Create()
		return err
	}

	var affectedProjectIds []int64
	var affectedContractIds []int64

	for i := range records {
		if i > 0 && i%1000 == 0 {
			process, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(i)/float64(len(records))), 64)
			fmt.Println("已处理", i, "条实际收款记录，当前进度：", fmt.Sprintf("%.0f", process*100), "%")
		} else if i == len(records)-1 {
			fmt.Println("已处理", i, "条实际收款记录，当前进度：100 %")
		}

		var tempCount int64
		global.DB.Model(&model.IncomeAndExpenditure{}).
			Where("imported_approval_id = ?", records[i].ImportedId+records[i].IOrd).
			Count(&tempCount)

		if tempCount > 0 {
			continue
		}

		records[i].Date = records[i].Date[:10]

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

		var detailedCurrency model.DictionaryDetail
		if records[i].Currency != "" {
			err = global.DB.
				Where("dictionary_type_id = ?", currency.Id).
				Where("name = ?", records[i].Currency).
				First(&detailedCurrency).Error
			if err != nil {
				param := service.ErrorLogCreate{
					Detail: "tabShouKuan视图的记录中发现无法匹配的币种：" +
						records[i].Currency + "，审批id为：" + records[i].ImportedId,
				}
				param.Create()
			}
		}

		var contract model.Contract
		if records[i].ContractCode != "" {
			err = global.DB.Model(&model.Contract{}).
				Where("code = ?", records[i].ContractCode).
				First(&contract).Error
			if err != nil {
				param := service.ErrorLogCreate{
					Detail: "tabShouKuan视图的记录中发现无法匹配的合同编号：" +
						records[i].ContractCode + "，合同编码为：" + records[i].ContractCode,
				}
				param.Create()
			}

			affectedContractIds = append(affectedContractIds, contract.Id)
		}

		var project model.Project
		if records[i].ProjectCode != "" {
			err = global.DB.
				Where("code = ?", records[i].ProjectCode).
				First(&project).Error
			if err != nil {
				param := service.ErrorLogCreate{
					Detail: "tabShouKuan视图的记录中发现无法匹配的项目编号：" +
						records[i].ProjectCode,
				}
				param.Create()
			}

			affectedProjectIds = append(affectedProjectIds, project.Id)
		}

		var relatedParty model.RelatedParty
		if records[i].ImportedRelatedPartyName != "" {
			err = global.DB.
				Where("name = ?", strings.TrimSpace(records[i].ImportedRelatedPartyName)).
				First(&relatedParty).Error
			if err != nil {
				err = global.DB.
					Where("imported_original_name like ?", "%"+strings.TrimSpace(records[i].ImportedRelatedPartyName)+"%").
					First(&relatedParty).Error
				if err != nil {
					param := service.ErrorLogCreate{
						Detail: "tabShouKuan视图的记录中发现无法匹配的相关方名称：" +
							records[i].ImportedRelatedPartyName + "审批id为：" + records[i].ImportedId,
					}
					param.Create()
				}
			}
		}

		exchangeRate := 1.0

		newRecord := service.IncomeAndExpenditureCreate{
			IgnoreUpdatingCumulativeIncomeAndExpenditure: true,
			UserId:             userId,
			ProjectId:          project.Id,
			ContractId:         contract.Id,
			Kind:               "实际",
			FundDirection:      "收款",
			Currency:           currency.Id,
			Date:               records[i].Date,
			Amount:             &records[i].Amount,
			ExchangeRate:       &exchangeRate,
			ImportedApprovalId: records[i].ImportedId + records[i].IOrd,
			DataSource:         "收款",
		}

		errCode := newRecord.Create()

		if errCode != util.Success {
			param := service.ErrorLogCreate{
				Detail: "导入tabShouKuan视图的记录时发生错误：" +
					util.GetErrorDescription(errCode) + "，id为：" +
					records[i].ImportedId + "，iOrd为：" + records[i].IOrd,
			}
			param.Create()
		}
	}

	err = updateCumulativeIncome(userId, affectedProjectIds, affectedContractIds)
	if err != nil {
		return err
	}

	fmt.Println("★★★★★人民币实际收款记录处理完成......★★★★★")

	return nil
}

type tabShouHui struct {
	Date                     string  `gorm:"column:F14395"`
	Amount                   float64 `gorm:"column:F16851"`
	Currency                 string  `gorm:"column:F14168"`
	ProjectCode              string  `gorm:"column:F16856"`
	ImportedRelatedPartyName string  `gorm:"column:F14394"`
	BankSerialId             string  `gorm:"column:F14165"`
	IOrd                     string  `gorm:"column:iOrd"`
}

func ImportActualIncomeFromTabShouHui(userId int64) error {
	fmt.Println("★★★★★开始处理实际收汇记录......★★★★★")

	//只处理“完成申报”的记录，别的状态可能会发生修改
	var records []tabShouHui
	global.DBForLvmin.Table("tabShouHui").
		Where("F16851 > 0").
		Where("F14173 = '完成申报'").
		Find(&records)

	var currency model.DictionaryType
	err := global.DB.
		Where("name = ?", "币种").
		First(&currency).Error
	if err != nil {
		param := service.ErrorLogCreate{
			Detail: "dictionaryType表中找不到”币种“这个名称",
		}
		param.Create()
		return err
	}

	var fundDirection model.DictionaryType
	var income model.DictionaryDetail
	err = global.DB.
		Where("name = '收付款的资金方向'").
		First(&fundDirection).Error
	if err != nil {
		param := service.ErrorLogCreate{
			Detail: "在dictionary_type表中找不到”收付款的资金方向“这个名称",
		}
		param.Create()
		return err
	}
	err = global.DB.
		Where("dictionary_type_id =?", fundDirection.Id).
		Where("name = '收款'").
		First(&income).Error
	if err != nil {
		param := service.ErrorLogCreate{
			Detail: "在dictionary_detail表中找不到”收款“这个名称",
		}
		param.Create()
		return err
	}

	var kind model.DictionaryType
	var actual model.DictionaryDetail
	err = global.DB.
		Where("name = '收付款的种类'").
		First(&kind).Error
	if err != nil {
		param := service.ErrorLogCreate{
			Detail: "在dictionary_type表中找不到”收付款的种类“这个名称",
		}
		param.Create()
		return err
	}
	err = global.DB.
		Where("dictionary_type_id =?", kind.Id).
		Where("name = '实际'").
		First(&actual).Error
	if err != nil {
		param := service.ErrorLogCreate{
			Detail: "在dictionary_detail表中找不到”实际“这个名称(dictionaryType为：收付款的种类)",
		}
		param.Create()
		return err
	}

	var affectedProjectIds []int64
	var affectedContractIds []int64

	for i := range records {
		if i > 0 && i%1000 == 0 {
			process, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(i)/float64(len(records))), 64)
			fmt.Println("已处理", i, "条实际收汇记录，当前进度：", fmt.Sprintf("%.0f", process*100), "%")
		} else if i == len(records)-1 {
			fmt.Println("已处理", i, "条实际收汇记录，当前进度：100 %")
		}

		var tempCount int64
		global.DB.Model(&model.IncomeAndExpenditure{}).
			Where("imported_approval_id = ?", records[i].BankSerialId+records[i].IOrd).
			Count(&tempCount)

		if tempCount > 0 {
			continue
		}

		records[i].Date = records[i].Date[:10]

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

		var detailedCurrency model.DictionaryDetail
		if records[i].Currency != "" {
			err = global.DB.
				Where("dictionary_type_id = ?", currency.Id).
				Where("name = ?", records[i].Currency).
				First(&detailedCurrency).Error
			if err != nil {
				param := service.ErrorLogCreate{
					Detail: "tabShouHui视图的记录中发现无法匹配的币种：" +
						records[i].Currency + "，银行流水号为：" + records[i].BankSerialId + "，IOrd为：" + records[i].IOrd,
				}
				param.Create()
			}
		}

		var project model.Project
		if records[i].ProjectCode != "" {
			err = global.DB.
				Where("code = ?", records[i].ProjectCode).
				First(&project).Error
			if err != nil {
				param := service.ErrorLogCreate{
					Detail: "tabShouHui视图的记录中发现无法匹配的项目编号：" +
						records[i].ProjectCode +
						"，银行流水号为：" + records[i].BankSerialId +
						"，IOrd为：" + records[i].IOrd,
				}
				param.Create()
			}

			affectedProjectIds = append(affectedProjectIds, project.Id)
		}

		var relatedParty model.RelatedParty
		if records[i].ImportedRelatedPartyName != "" {
			err = global.DB.
				Where("name = ?", strings.TrimSpace(records[i].ImportedRelatedPartyName)).
				First(&relatedParty).Error
			if err != nil {
				err = global.DB.
					Where("imported_original_name like ?", "%"+strings.TrimSpace(records[i].ImportedRelatedPartyName)+"%").
					First(&relatedParty).Error
				if err != nil {
					param := service.ErrorLogCreate{
						Detail: "tabShouHui视图的记录中发现无法匹配的相关方名称：" +
							records[i].ImportedRelatedPartyName + "，银行流水号为：" + records[i].BankSerialId + "，IOrd为：" + records[i].IOrd,
					}
					param.Create()
				}
			}
		}

		newRecord := service.IncomeAndExpenditureCreate{
			IgnoreUpdatingCumulativeIncomeAndExpenditure: true,
			UserId:             userId,
			ProjectId:          project.Id,
			Kind:               "实际",
			FundDirection:      "收款",
			Currency:           detailedCurrency.Id,
			Date:               records[i].Date,
			Amount:             &records[i].Amount,
			ImportedApprovalId: records[i].BankSerialId + records[i].IOrd,
			DataSource:         "收汇",
		}

		switch records[i].Currency {
		case "人民币":
			newRecord.ExchangeRate = model.Float64ToPointer(1)
		case "美元":
			newRecord.ExchangeRate = &global.Config.ExchangeRate.USD
		case "欧元":
			newRecord.ExchangeRate = &global.Config.ExchangeRate.EUR
		case "港币":
			newRecord.ExchangeRate = &global.Config.ExchangeRate.HKD
		case "新加坡元":
			newRecord.ExchangeRate = &global.Config.ExchangeRate.SGD
		case "马来西亚币":
			newRecord.ExchangeRate = &global.Config.ExchangeRate.MYR
		default:
			newRecord.ExchangeRate = model.Float64ToPointer(1)
		}

		errCode := newRecord.Create()

		if errCode != util.Success {
			param := service.ErrorLogCreate{
				Detail: "导入tabShouHui视图的记录时发生错误：" +
					util.GetErrorDescription(errCode) + "，银行流水id为：" +
					records[i].BankSerialId + "，iOrd为：" + records[i].IOrd,
			}
			param.Create()
		}
	}

	err = updateCumulativeIncome(userId, affectedProjectIds, affectedContractIds)
	if err != nil {
		return err
	}

	fmt.Println("★★★★★实际收汇记录处理完成......★★★★★")

	return nil
}

type tabShouPiao struct {
	Date                     string  `gorm:"column:F6446"`
	Amount                   float64 `gorm:"column:F6445"`
	ProjectCode              string  `gorm:"column:F10517"`
	ImportedRelatedPartyName string  `gorm:"column:F6443"`
	BankSerialId             string  `gorm:"column:F8544"`
}

func ImportActualIncomeFromTabShouPiao(userId int64) error {
	fmt.Println("★★★★★开始处理实际收票记录......★★★★★")

	//只处理“确认全部收票金额”的记录，别的状态可能会发生修改
	var records []tabShouPiao
	global.DBForLvmin.Table("tabShouPiao").
		Where("F6445 > 0").
		Where("F6442 = '确认全部收票金额'").
		Find(&records)

	var currency model.DictionaryType
	err := global.DB.
		Where("name = ?", "币种").
		First(&currency).Error
	if err != nil {
		param := service.ErrorLogCreate{
			Detail: "dictionary_type表中找不到”币种“这个名称",
		}
		param.Create()
		return err
	}
	var CNY model.DictionaryDetail
	err = global.DB.
		Where("dictionary_type_id = ?", currency.Id).
		Where("name = ?", "人民币").
		First(&CNY).Error
	if err != nil {
		param := service.ErrorLogCreate{
			Detail: "dictionary_detail表中找不到”人民币“这个名称",
		}
		param.Create()
		return err
	}

	var fundDirection model.DictionaryType
	var income model.DictionaryDetail
	err = global.DB.
		Where("name = '收付款的资金方向'").
		First(&fundDirection).Error
	if err != nil {
		param := service.ErrorLogCreate{
			Detail: "在dictionary_type表中找不到”收付款的资金方向“这个名称",
		}
		param.Create()
		return err
	}
	err = global.DB.
		Where("dictionary_type_id =?", fundDirection.Id).
		Where("name = '收款'").
		First(&income).Error
	if err != nil {
		param := service.ErrorLogCreate{
			Detail: "在dictionary_detail表中找不到”收款“这个名称",
		}
		param.Create()
		return err
	}

	var kind model.DictionaryType
	var actual model.DictionaryDetail
	err = global.DB.
		Where("name = '收付款的种类'").
		First(&kind).Error
	if err != nil {
		param := service.ErrorLogCreate{
			Detail: "在dictionary_type表中找不到”收付款的种类“这个名称",
		}
		param.Create()
		return err
	}
	err = global.DB.
		Where("dictionary_type_id =?", kind.Id).
		Where("name = '实际'").
		First(&actual).Error
	if err != nil {
		param := service.ErrorLogCreate{
			Detail: "在dictionary_detail表中找不到”实际“这个名称(dictionaryType为：收付款的种类)",
		}
		param.Create()
		return err
	}

	var affectedProjectIds []int64
	var affectedContractIds []int64

	for i := range records {
		if i > 0 && i%1000 == 0 {
			process, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(i)/float64(len(records))), 64)
			fmt.Println("已处理", i, "条实际收票记录，当前进度：", fmt.Sprintf("%.0f", process*100), "%")
		} else if i == len(records)-1 {
			fmt.Println("已处理", i, "条实际收票记录，当前进度：100 %")
		}

		var tempCount int64
		global.DB.Model(&model.IncomeAndExpenditure{}).
			Where("imported_approval_id = ?", records[i].BankSerialId).
			Count(&tempCount)

		if tempCount > 0 {
			continue
		}

		records[i].Date = records[i].Date[:10]

		var project model.Project
		if records[i].ProjectCode != "" {
			err = global.DB.
				Where("code = ?", records[i].ProjectCode).
				First(&project).Error
			if err != nil {
				param := service.ErrorLogCreate{
					Detail: "tabShouPiao视图的记录中发现无法匹配的项目编号：" +
						records[i].ProjectCode +
						"，银行单号为：" + records[i].BankSerialId,
				}
				param.Create()
			}

			affectedProjectIds = append(affectedProjectIds, project.Id)
		}

		var relatedParty model.RelatedParty
		if records[i].ImportedRelatedPartyName != "" {
			err = global.DB.
				Where("name = ?", strings.TrimSpace(records[i].ImportedRelatedPartyName)).
				First(&relatedParty).Error
			if err != nil {
				err = global.DB.
					Where("imported_original_name like ?", "%"+strings.TrimSpace(records[i].ImportedRelatedPartyName)+"%").
					First(&relatedParty).Error
				if err != nil {
					param := service.ErrorLogCreate{
						Detail: "tabShouHui视图的记录中发现无法匹配的相关方名称：" +
							records[i].ImportedRelatedPartyName + "，银行单号为：" + records[i].BankSerialId,
					}
					param.Create()
				}
			}
		}

		exchangeRate := 1.0

		newRecord := service.IncomeAndExpenditureCreate{
			IgnoreUpdatingCumulativeIncomeAndExpenditure: true,
			UserId:             userId,
			ProjectId:          project.Id,
			Kind:               "实际",
			FundDirection:      "收款",
			Currency:           CNY.Id,
			ExchangeRate:       &exchangeRate,
			Date:               records[i].Date,
			Amount:             &records[i].Amount,
			ImportedApprovalId: records[i].BankSerialId,
			DataSource:         "收票",
		}

		errCode := newRecord.Create()

		if errCode != util.Success {
			param := service.ErrorLogCreate{
				Detail: "导入tabShouPiao视图的记录时发生错误：" +
					util.GetErrorDescription(errCode) + "，银行单号为：" +
					records[i].BankSerialId,
			}
			param.Create()
		}
	}

	err = updateCumulativeIncome(userId, affectedProjectIds, affectedContractIds)
	if err != nil {
		return err
	}

	fmt.Println("★★★★★实际收票记录处理完成......★★★★★")

	return nil
}
