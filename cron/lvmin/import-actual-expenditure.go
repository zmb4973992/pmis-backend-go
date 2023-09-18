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

type tabFukuan2a struct {
	Date                     string  `gorm:"column:F7854"`
	Amount                   float64 `gorm:"column:F7538"`
	Currency                 string  `gorm:"column:F7539"`
	ContractCode             string  `gorm:"column:F7535"`
	ProjectCode              string  `gorm:"column:F7546"`
	Type                     string  `gorm:"column:F7541"`
	Remarks                  string  `gorm:"column:F8666"`
	ImportedApprovalId       string  `gorm:"column:F7549"`
	ImportedRelatedPartyName string  `gorm:"column:F13591"`
}

func ImportActualExpenditure(userId int64) error {
	fmt.Println("★★★★★开始处理实际付款记录......★★★★★")

	var records []tabFukuan2a
	//F7538是金额，F7549是外部导入的付款审批id
	global.DBForLvmin.Table("tabFukuan2").
		Where("F7538 > 0").
		Where("F7549 is not null").
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
	}

	err = global.DB.
		Where("dictionary_type_id =?", fundDirection.Id).
		Where("name = '付款'").
		First(&expenditure).Error
	if err != nil {
		param := service.ErrorLogCreate{
			Detail: "在dictionary_detail表中找不到”付款“这个名称",
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
			fmt.Println("已处理", i, "条实际付款记录，当前进度：", fmt.Sprintf("%.0f", process*100), "%")
		} else if i == len(records)-1 {
			fmt.Println("已处理", i, "条实际付款记录，当前进度：100 %")
		}

		var tempCount int64
		global.DB.Model(&model.IncomeAndExpenditure{}).
			Where("imported_approval_id = ?", records[i].ImportedApprovalId).
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

		var specificCurrency model.DictionaryDetail
		if records[i].Currency != "" {
			err = global.DB.
				Where("dictionary_type_id = ?", currency.Id).
				Where("name = ?", records[i].Currency).
				First(&specificCurrency).Error
			if err != nil {
				param := service.ErrorLogCreate{
					Detail: "tabFukuan2视图的记录中发现无法匹配的币种：" +
						records[i].Currency + "，付款审批id为：" + records[i].ImportedApprovalId,
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
					Detail: "tabFukuan2视图的记录中发现无法匹配的合同编号：" +
						records[i].ContractCode + "，付款审批id为：" + records[i].ImportedApprovalId,
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
					Detail: "tabFukuan2视图的记录中发现无法匹配的项目编号：" +
						records[i].ProjectCode + "，付款审批id为：" + records[i].ImportedApprovalId,
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
						Detail: "tabFukuan2视图的记录中发现无法匹配的相关方名称：" +
							records[i].ImportedRelatedPartyName + "付款审批id为：" + records[i].ImportedApprovalId,
					}
					param.Create()
				}
			}
		}

		newRecord := service.IncomeAndExpenditureCreate{
			IgnoreUpdatingCumulativeIncomeAndExpenditure: true,
			UserId:             userId,
			ProjectId:          project.Id,
			ContractId:         contract.Id,
			Kind:               "实际",
			FundDirection:      "付款",
			Currency:           specificCurrency.Id,
			Date:               records[i].Date,
			Amount:             &records[i].Amount,
			Type:               records[i].Type,
			Remarks:            records[i].Remarks,
			ImportedApprovalId: records[i].ImportedApprovalId,
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
				Detail: "导入tabFukuan2视图的记录时发生错误：" +
					util.GetErrorDescription(errCode) +
					"，付款审批id为：" + records[i].ImportedApprovalId,
			}
			param.Create()
		}
	}

	err = updateCumulativeExpenditure(userId, affectedProjectIds, affectedContractIds)
	if err != nil {
		return err
	}

	fmt.Println("★★★★★所有实际付款记录处理完成......★★★★★")

	return nil
}
