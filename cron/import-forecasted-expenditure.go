package cron

import (
	"errors"
	"fmt"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/service"
	"strconv"
	"strings"
	"time"
)

type tabFukuan struct {
	Date               string  `gorm:"column:F6685"`
	Amount             float64 `gorm:"column:F7117"`
	Currency           string  `gorm:"column:F7128"`
	ContractCode       string  `gorm:"column:F8016"`
	Type               string  `gorm:"column:F6681"`
	Remarks            string  `gorm:"column:F6682"`
	ImportedApprovalID string  `gorm:"column:F7120"`
	ExchangeRate       float64
}

func importForecastedExpenditure() error {
	fmt.Println("★★★★★开始处理预测付款记录......★★★★★")

	var records []tabFukuan
	//F7117是金额，F7120是外部导入的付款审批id
	global.DB2.Table("tabFukuan").Where("F7117 > 0").
		Where("F7120 is not null").
		Find(&records)

	for i := range records {
		if i > 0 && i%100 == 0 {
			process, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(i)/float64(len(records))), 64)
			fmt.Println("已处理", i, "条预测付款记录，当前进度：", fmt.Sprintf("%.0f", process*100), "%")
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
						Date:   time.Now().Format("2006-01-02"),
					}
					param.Create()
					records[i].Currency = ""
				} else {
					err = global.DB.Model(&model.DictionaryDetail{}).
						Where("name = ?", records[i].Currency).Select("id").
						First(&currencyID).Error
					if err != nil {
						param := service.ErrorLogCreate{
							Detail: "tabFukuan视图的记录中发现无法匹配的币种：" +
								records[i].Currency + "，付款审批ID为：" + records[i].ImportedApprovalID,
							Date: time.Now().Format("2006-01-02"),
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
						Detail: "tabFukuan视图的记录中发现无法匹配的合同编号：" +
							records[i].ContractCode + "，付款审批id为：" + records[i].ImportedApprovalID,
						Date: time.Now().Format("2006-01-02"),
					}
					param.Create()
					records[i].ContractCode = ""
				}
			}

			var projectID int64
			if records[i].ContractCode != "" {
				err := global.DB.Model(&model.Contract{}).
					Where("code = ?", records[i].ContractCode).Select("project_id").
					First(&projectID).Error
				if err != nil {
					param := service.ErrorLogCreate{
						Detail: "tabFukuan视图的记录中发现该合同号找不到对应的项目号：" +
							records[i].ContractCode + "，付款审批id为：" + records[i].ImportedApprovalID,
						Date: time.Now().Format("2006-01-02"),
					}
					param.Create()
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
						Date:   time.Now().Format("2006-01-02"),
					}
					param.Create()
					records[i].Type = ""
				} else {
					err = global.DB.Model(&model.DictionaryDetail{}).
						Where("name = ?", records[i].Type).Select("id").
						First(&typeID).Error
					if err != nil {
						param := service.ErrorLogCreate{
							Detail: "tabFukuan视图的记录中发现无法匹配的款项类型：" +
								records[i].Type + "，付款审批ID为：" + records[i].ImportedApprovalID,
							Date: time.Now().Format("2006-01-02"),
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
				Kind:               "预测",
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
				return errors.New(res.Message)
			}

			if res.Code != 0 {
				param := service.ErrorLogCreate{
					Detail: "导入tabFukuan视图的记录时发生错误：" +
						res.Message + "，付款审批ID为：" + records[i].ImportedApprovalID,
					Date: time.Now().Format("2006-01-02"),
				}
				res = param.Create()
				if res.Code != 0 {
					return errors.New(res.Message)
				}
			}
		}
	}

	return nil
}
