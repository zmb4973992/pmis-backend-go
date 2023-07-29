package cron

import (
	"fmt"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/service"
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
	ImportedApprovalID       string  `gorm:"column:F7549"`
	ImportedRelatedPartyName string  `gorm:"column:F13591"`
	ExchangeRate             float64
}

func importActualExpenditure() error {
	fmt.Println("★★★★★开始处理实际付款记录......★★★★★")

	var records []tabFukuan2a
	//F7538是金额，F7549是外部导入的付款审批id
	global.DB2.Table("tabFukuan2").Where("F7538 > 0").
		Where("F7549 is not null").
		Find(&records)

	for i := range records {
		if i > 0 && i%100 == 0 {
			process, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(i)/float64(len(records))), 64)
			fmt.Println("已处理", i, "条实际付款记录，当前进度：", fmt.Sprintf("%.0f", process*100), "%")
		} else if i == len(records)-1 {
			fmt.Println("已处理", i, "条实际付款记录，当前进度：100 %")
		}

		var tempCount int64
		global.DB.Model(&model.IncomeAndExpenditure{}).
			Where("imported_approval_id = ?", records[i].ImportedApprovalID).
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
				case "1":
					records[i].Currency = "人民币"
				case "2":
					records[i].Currency = "美元"
				case "3":
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
							Detail: "tabFukuan2视图的记录中发现无法匹配的币种：" +
								records[i].Currency + "，付款审批ID为：" + records[i].ImportedApprovalID,
						}
						param.Create()
						records[i].Currency = ""
					}
				}
			}

			var contractID int64
			if records[i].ContractCode != "" {
				err := global.DB.Model(&model.Contract{}).
					Where("code = ?", records[i].ContractCode).Select("id").
					First(&contractID).Error
				if err != nil {
					param := service.ErrorLogCreate{
						Detail: "tabFukuan2视图的记录中发现无法匹配的合同编号：" +
							records[i].ContractCode + "，合同编码为：" + records[i].ContractCode,
					}
					param.Create()
					records[i].ContractCode = ""
				}
			}

			var projectID int64
			if records[i].ProjectCode != "" {
				err := global.DB.Model(&model.Project{}).
					Where("code = ?", records[i].ProjectCode).Select("id").
					First(&projectID).Error
				if err != nil {
					param := service.ErrorLogCreate{
						Detail: "tabFukuan2视图的记录中发现无法匹配的项目编号：" +
							records[i].ProjectCode + "，项目编号为：" + records[i].ProjectCode,
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
							Detail: "tabFukuan2视图的记录中发现无法匹配的相关方名称：" +
								records[i].ImportedRelatedPartyName + "付款审批ID为：" + records[i].ImportedApprovalID,
						}
						param.Create()
						records[i].ImportedRelatedPartyName = ""
					}
				}
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

			var typeID int64
			if records[i].Type != "" {
				var dictionaryTypeID int64
				err := global.DB.Model(&model.DictionaryType{}).
					Where("name = ?", "款项类型").Select("id").
					First(&dictionaryTypeID).Error
				if err != nil {
					param := service.ErrorLogCreate{
						Detail: "dictionaryType表中找不到”款项类型“这个名称",
					}
					param.Create()
					records[i].Type = ""
				} else {
					err = global.DB.Model(&model.DictionaryDetail{}).
						Where("name = ?", records[i].Type).Select("id").
						First(&typeID).Error
					if err != nil {
						param := service.ErrorLogCreate{
							Detail: "tabFukuan2视图的记录中发现无法匹配的款项类型：" +
								records[i].Type + "，付款审批ID为：" + records[i].ImportedApprovalID,
						}
						param.Create()
						records[i].Type = ""
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
				ContractID:         contractID,
				FundDirection:      "付款",
				Currency:           currencyID,
				Kind:               "实际",
				Date:               records[i].Date,
				Amount:             &records[i].Amount,
				ExchangeRate:       &records[i].ExchangeRate,
				Type:               typeID,
				Term:               0,
				Remarks:            records[i].Remarks,
				Attachment:         "",
				ImportedApprovalID: records[i].ImportedApprovalID,
			}

			res := newRecord.Create()

			if res.Code != 0 {
				param := service.ErrorLogCreate{
					Detail: "导入tabFukuan2视图的记录时发生错误：" +
						res.Message + "，付款审批ID为：" + records[i].ImportedApprovalID,
				}
				param.Create()
			}
		}
	}

	return nil
}
