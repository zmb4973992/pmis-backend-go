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

type tabOriPayPlan struct {
	Date               string  `gorm:"column:F12390"`
	Amount             float64 `gorm:"column:F12387"`
	Currency           string  `gorm:"column:F12388"`
	ContractCode       string  `gorm:"column:F12383"`
	Type               string  `gorm:"column:F12384"`
	Remarks            string  `gorm:"column:F12385"`
	ImportedApprovalID string  `gorm:"column:F12395"`
}

func ImportPlannedExpenditure(userID int64) error {
	fmt.Println("★★★★★开始处理计划付款记录......★★★★★")

	var records []tabOriPayPlan
	//F12387是金额，F12395是外部导入的付款审批id
	global.DBForLvmin.Table("tabOriPayPlan").
		Where("F12387 > 0").
		Where("F12395 is not null").
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
	var expenditure model.DictionaryDetail
	err = global.DB.
		Where("name = '收付款的资金方向'").
		First(&fundDirection).Error
	if err != nil {
		param := service.ErrorLogCreate{
			Detail: "在dictionary_type表中找不到”收付款的资金方向“这个名称",
		}
		param.Create()
		return err
	} else {
		err = global.DB.
			Where("dictionary_type_id =?", fundDirection.ID).
			Where("name = '付款'").
			First(&expenditure).Error
		if err != nil {
			param := service.ErrorLogCreate{
				Detail: "在dictionary_detail表中找不到”付款“这个名称",
			}
			param.Create()
			return err
		}
	}

	var kind model.DictionaryType
	var forecasted model.DictionaryDetail
	err = global.DB.
		Where("name = '收付款的种类'").
		First(&kind).Error
	if err != nil {
		param := service.ErrorLogCreate{
			Detail: "在dictionary_type表中找不到”收付款的种类“这个名称",
		}
		param.Create()
		return err
	} else {
		err = global.DB.
			Where("dictionary_type_id =?", kind.ID).
			Where("name = '计划'").
			First(&forecasted).Error
		if err != nil {
			param := service.ErrorLogCreate{
				Detail: "在dictionary_detail表中找不到”计划“这个名称(dictionaryType为：收付款的种类)",
			}
			param.Create()
			return err
		}
	}

	var affectedProjectIDs []int64
	var affectedContractIDs []int64

	for i := range records {
		if i > 0 && i%1000 == 0 {
			process, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(i)/float64(len(records))), 64)
			fmt.Println("已处理", i, "条计划付款记录，当前进度：", fmt.Sprintf("%.0f", process*100), "%")
		} else if i == len(records)-1 {
			fmt.Println("已处理", i, "条计划付款记录，当前进度：100 %")
		}

		var tempCount int64
		global.DB.Model(&model.IncomeAndExpenditure{}).
			Where("imported_approval_id = ?", records[i].ImportedApprovalID).
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
				Where("dictionary_type_id = ?", currency.ID).
				Where("name = ?", records[i].Currency).
				First(&detailedCurrency).Error
			if err != nil {
				param := service.ErrorLogCreate{
					Detail: "tabOriPayPlan视图的记录中发现无法匹配的币种：" +
						records[i].Currency + "，付款审批ID为：" + records[i].ImportedApprovalID,
				}
				param.Create()
			}
		}

		var contract model.Contract
		if records[i].ContractCode != "" {
			err = global.DB.
				Where("code = ?", records[i].ContractCode).
				First(&contract).Error
			if err != nil {
				param := service.ErrorLogCreate{
					Detail: "tabOriPayPlan视图的记录中发现无法匹配的合同编号：" +
						records[i].ContractCode + "，付款审批id为：" + records[i].ImportedApprovalID,
				}
				param.Create()
			}

			affectedContractIDs = append(affectedContractIDs, contract.ID)
		}

		var projectID int64
		if contract.ProjectID != nil {
			projectID = *contract.ProjectID
			affectedProjectIDs = append(affectedProjectIDs, projectID)
		}

		switch {
		case strings.Contains(records[i].Type, "预付"):
			records[i].Type = "预付款"
		case strings.Contains(records[i].Type, "定金"):
			records[i].Type = "定金"
		case strings.Contains(records[i].Type, "发货款"):
			records[i].Type = "发货款"
		case strings.Contains(records[i].Type, "发运款"):
			records[i].Type = "发货款"
		case strings.Contains(records[i].Type, "货款"):
			records[i].Type = "发货款"
		case strings.Contains(records[i].Type, "港杂费"):
			records[i].Type = "港杂费"
		case strings.Contains(records[i].Type, "进度款"):
			records[i].Type = "进度款"
		case strings.Contains(records[i].Type, "调试款"):
			records[i].Type = "调试款"
		case strings.Contains(records[i].Type, "杂费"):
			records[i].Type = "杂费"
		case strings.Contains(records[i].Type, "租金"):
			records[i].Type = "租金"
		case strings.Contains(records[i].Type, "服务费"):
			records[i].Type = "服务费"
		case strings.Contains(records[i].Type, "保证金"):
			records[i].Type = "保证金"
		case strings.Contains(records[i].Type, "保费"):
			records[i].Type = "保费"
		case strings.Contains(records[i].Type, "尾款"):
			records[i].Type = "尾款"
		default:
			records[i].Type = ""
		}

		newRecord := service.IncomeAndExpenditureCreate{
			IgnoreUpdatingCumulativeIncomeAndExpenditure: true,
			UserID:             userID,
			ProjectID:          projectID,
			ContractID:         contract.ID,
			Kind:               "计划",
			FundDirection:      "付款",
			Currency:           detailedCurrency.ID,
			Date:               records[i].Date,
			Amount:             &records[i].Amount,
			Type:               records[i].Type,
			Remarks:            records[i].Remarks,
			ImportedApprovalID: records[i].ImportedApprovalID,
		}

		switch records[i].Currency {
		case "人民币":
			newRecord.ExchangeRate = model.Float64ToPointer(1)
		case "美元":
			newRecord.ExchangeRate = model.Float64ToPointer(7.2)
		case "欧元":
			newRecord.ExchangeRate = model.Float64ToPointer(7.8)
		case "港币":
			newRecord.ExchangeRate = model.Float64ToPointer(0.9)
		case "新加坡元":
			newRecord.ExchangeRate = model.Float64ToPointer(5.2)
		case "马来西亚币":
			newRecord.ExchangeRate = model.Float64ToPointer(1.5)
		default:
			newRecord.ExchangeRate = model.Float64ToPointer(1)
		}

		errCode := newRecord.Create()

		if errCode != util.Success {
			param := service.ErrorLogCreate{
				Detail: "导入tabOriPayPlan视图的记录时发生错误：" +
					util.GetErrorDescription(errCode) + "，付款审批ID为：" +
					records[i].ImportedApprovalID,
			}
			param.Create()
		}
	}

	err = updateCumulativeExpenditure(userID, affectedProjectIDs, affectedContractIDs)
	if err != nil {
		return err
	}

	fmt.Println("★★★★★计划付款记录处理完成......★★★★★")

	return nil
}
