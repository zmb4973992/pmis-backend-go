package lvmin

import (
	"errors"
	"fmt"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/service"
	"strconv"
	"strings"
)

type tabContract struct {
	Code                      string  `gorm:"column:F6100"`
	Name                      string  `gorm:"column:F6099"`
	Amount                    float64 `gorm:"column:F6101"`
	RelatedParty              string  `gorm:"column:F6102"`
	Type                      string  `gorm:"column:F6110"`
	ProjectCode               string  `gorm:"column:F6482"`
	Organization              string  `gorm:"column:F6484"`
	OurSignatory              string  `gorm:"column:F6488"`
	Currency                  string  `gorm:"column:F6525"`
	FinancialAccountingNumber string  `gorm:"column:F8449"`
	Content                   string  `gorm:"column:F6487"`
	FundDirection             string  `gorm:"column:F12338"`
}

func ImportContract(userID int64) error {
	fmt.Println("★★★★★开始处理合同记录......★★★★★")

	var records []tabContract
	//主合同的定义是项目
	global.DBForLvmin.Table("tabContract").Where("F6110 != '主合同'").
		Find(&records)

	var contractType model.DictionaryType
	err := global.DB.
		Where("name = ?", "合同类型").
		First(&contractType).Error
	if err != nil {
		param := service.ErrorLogCreate{
			Detail: "dictionaryType表中找不到”合同类型“这个名称",
		}
		param.Create()
		return err
	}

	var currency model.DictionaryType
	err = global.DB.
		Where("name = ?", "币种").
		First(&currency).Error
	if err != nil {
		param := service.ErrorLogCreate{
			Detail: "dictionaryType表中找不到”币种“这个名称",
		}
		param.Create()
		return err
	}

	var ContractFundDirection model.DictionaryType
	err = global.DB.
		Where("name = ?", "合同的资金方向").
		First(&ContractFundDirection).Error
	if err != nil {
		param := service.ErrorLogCreate{
			Detail: "dictionaryType表中找不到”合同的资金方向“这个名称",
		}
		param.Create()
		return err
	}

	var ourSignatory model.DictionaryType
	err = global.DB.
		Where("name = ?", "我方签约主体").
		First(&ourSignatory).Error
	if err != nil {
		param := service.ErrorLogCreate{
			Detail: "dictionaryType表中找不到”我方签约主体“这个名称",
		}
		param.Create()
		return err
	}

	for i := range records {
		if i > 0 && i%1000 == 0 {
			process, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(i)/float64(len(records))), 64)
			fmt.Println("已处理", i, "条合同记录，当前进度：", fmt.Sprintf("%.0f", process*100), "%")
		} else if i == len(records)-1 {
			fmt.Println("已处理", i, "条合同记录，当前进度：100 %")
		}

		var tempCount int64
		global.DB.Model(&model.Contract{}).
			Where("code = ?", records[i].Code).
			Count(&tempCount)
		if tempCount == 0 {
			var organization model.Organization
			if records[i].Organization != "" {
				switch records[i].Organization {
				case "机械车辆部":
					records[i].Organization = "成套业务一部"
				case "成套六部":
					records[i].Organization = "成套业务六部"
				case "事业部管理委员会和水泥工程事业部":
					records[i].Organization = "水泥工程事业部"
				case "项目管理及技术支持部":
					records[i].Organization = "项目管理部"
				case "综合管理部":
					records[i].Organization = "综合管理和法律部"
				case "工程项目管理部":
					records[i].Organization = "项目管理部"
				case "党建和纪检审计部":
					records[i].Organization = "党建文宣部"
				case "成套业务三部":
					records[i].Organization = "成套业务四部"
				}

				err = global.DB.
					Where("name = ?", records[i].Organization).
					First(&organization).Error
				if err != nil {
					param := service.ErrorLogCreate{
						Detail: "tabContract视图的记录中发现无法匹配的部门：" +
							records[i].Organization + "，合同编码为：" + records[i].Code,
					}
					param.Create()
					records[i].Organization = ""
				}
			}

			var specificContractType model.DictionaryDetail
			if records[i].Type != "" {
				switch records[i].Type {
				case "服务类合同":
					records[i].Type = "技服"
				case "工程类合同":
					records[i].Type = "工程"
				case "国内采购":
					records[i].Type = "采购"
				case "国内销售":
					records[i].Type = "销售"
				case "库存采购":
					records[i].Type = "采购"
				case "库存销售":
					records[i].Type = "销售"
				case "贸易类合同":
					records[i].Type = "贸易"
				case "特定采购订单":
					records[i].Type = "采购"
				case "延伸业务调价协议":
					records[i].Type = "其他"
				case "其它":
					records[i].Type = "其他"
				}

				err = global.DB.
					Where("dictionary_type_id = ?", contractType.ID).
					Where("name = ?", records[i].Type).
					First(&specificContractType).Error
				if err != nil {
					param := service.ErrorLogCreate{
						Detail: "tabContract视图的记录中发现无法匹配的合同类型：" +
							records[i].Type + "，合同编码为：" + records[i].Code,
					}
					param.Create()
					records[i].Type = ""
				}
			}

			var specificCurrency model.DictionaryDetail
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

				err = global.DB.
					Where("dictionary_type_id = ?", currency.ID).
					Where("name = ?", records[i].Currency).
					First(&specificCurrency).Error
				if err != nil {
					param := service.ErrorLogCreate{
						Detail: "tabContract视图的记录中发现无法匹配的币种：" +
							records[i].Currency + "，合同编码为：" + records[i].Code,
					}
					param.Create()
					records[i].Currency = ""
				}
			}

			var project model.Project
			if records[i].ProjectCode != "" {
				err = global.DB.
					Where("code = ?", records[i].ProjectCode).
					First(&project).Error
				if err != nil {
					param := service.ErrorLogCreate{
						Detail: "tabContract视图的记录中发现无法匹配的项目编号：" +
							records[i].ProjectCode + "，合同编码为：" + records[i].Code,
					}
					param.Create()
					records[i].ProjectCode = ""
				}
			}

			var relatedParty model.RelatedParty
			if records[i].RelatedParty != "" {
				err = global.DB.
					Where("name = ?", records[i].RelatedParty).
					First(&relatedParty).Error
				if err != nil {
					err = global.DB.
						Where("imported_original_name like ?", "%"+strings.TrimSpace(records[i].RelatedParty)+"%").
						First(&relatedParty).Error
					if err != nil {
						param := service.ErrorLogCreate{
							Detail: "tabContract视图的记录中发现无法匹配的相关方名称：" +
								records[i].RelatedParty + "，合同编码为：" + records[i].Code,
						}
						param.Create()
						records[i].RelatedParty = ""
					}
				}
			}

			var specificContractFundDirection model.DictionaryDetail
			if records[i].FundDirection != "" {
				switch records[i].FundDirection {
				case "收款":
					records[i].FundDirection = "收款合同"
				case "付款":
					records[i].FundDirection = "付款合同"
				}

				err = global.DB.
					Where("dictionary_type_id = ?", ContractFundDirection.ID).
					Where("name = ?", records[i].FundDirection).
					First(&specificContractFundDirection).Error
				if err != nil {
					param := service.ErrorLogCreate{
						Detail: "tabContract视图的记录中发现无法匹配的合同的资金方向：" +
							records[i].FundDirection + "，合同编码为：" + records[i].Code,
					}
					param.Create()
					records[i].Currency = ""
				}

			}

			var specificOurSignatory model.DictionaryDetail
			if records[i].FundDirection != "" {
				err = global.DB.
					Where("dictionary_type_id = ?", ourSignatory.ID).
					Where("name = ?", records[i].OurSignatory).
					First(&specificOurSignatory).Error
				if err != nil {
					param := service.ErrorLogCreate{
						Detail: "tabContract视图的记录中发现无法匹配的我方签约主体：" +
							records[i].OurSignatory + "，合同编码为：" + records[i].Code,
					}
					param.Create()
					records[i].Currency = ""
				}
			}

			newRecord := service.ContractCreate{
				UserID:         userID,
				ProjectID:      project.ID,
				OrganizationID: organization.ID,
				RelatedPartyID: relatedParty.ID,
				FundDirection:  specificContractFundDirection.ID,
				OurSignatory:   specificOurSignatory.ID,
				Currency:       specificCurrency.ID,
				Type:           specificContractType.ID,
				Amount:         &records[i].Amount,
				Name:           records[i].Name,
				Code:           records[i].Code,
				Content:        records[i].Content,
			}

			res := newRecord.Create()
			if res.Code != 0 {
				return errors.New(res.Message)
			}
		}
	}
	fmt.Println("★★★★★所有合同记录添加完成......★★★★★")

	return nil
}

func UpdateExchangeRageOfContract(userID int64) error {
	//处理合同的币种和汇率
	var contracts []model.Contract
	global.DB.
		Where("currency is not null").
		Where("exchange_rate is null").
		Find(&contracts)

	for i := range contracts {
		if i > 0 && i%1000 == 0 {
			process, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(i)/float64(len(contracts))), 64)
			fmt.Println("已处理", i, "条合同记录，当前进度：", fmt.Sprintf("%.0f", process*100), "%")
		} else if i == len(contracts)-1 {
			fmt.Println("已处理", i, "条合同记录，当前进度：100 %")
		}

		var specificCurrency model.DictionaryDetail
		err := global.DB.
			Where("id = ?", *contracts[i].Currency).
			First(&specificCurrency).Error
		if err != nil {
			continue
		}

		switch specificCurrency.Name {
		case "人民币":
			param := service.ContractUpdate{
				IgnoreDataAuthority: true,
				UserID:              userID,
				ContractID:          contracts[i].ID,
				ExchangeRate:        model.Float64ToPointer(1),
			}
			res := param.Update()
			if res.Code != 0 {
				return errors.New(res.Message)
			}

		case "美元":
			param := service.ContractUpdate{
				IgnoreDataAuthority: true,
				UserID:              userID,
				ContractID:          contracts[i].ID,
				ExchangeRate:        model.Float64ToPointer(7.2),
			}
			res := param.Update()
			if res.Code != 0 {
				return errors.New(res.Message)
			}

		case "欧元":
			param := service.ContractUpdate{
				IgnoreDataAuthority: true,
				UserID:              userID,
				ContractID:          contracts[i].ID,
				ExchangeRate:        model.Float64ToPointer(7.8),
			}
			res := param.Update()
			if res.Code != 0 {
				return errors.New(res.Message)
			}

		case "港币":
			param := service.ContractUpdate{
				IgnoreDataAuthority: true,
				UserID:              userID,
				ContractID:          contracts[i].ID,
				ExchangeRate:        model.Float64ToPointer(0.92),
			}
			res := param.Update()
			if res.Code != 0 {
				return errors.New(res.Message)
			}

		case "新加坡元":
			param := service.ContractUpdate{
				IgnoreDataAuthority: true,
				UserID:              userID,
				ContractID:          contracts[i].ID,
				ExchangeRate:        model.Float64ToPointer(5.3),
			}
			res := param.Update()
			if res.Code != 0 {
				return errors.New(res.Message)
			}

		case "马来西亚币":
			param := service.ContractUpdate{
				IgnoreDataAuthority: true,
				UserID:              userID,
				ContractID:          contracts[i].ID,
				ExchangeRate:        model.Float64ToPointer(1.5),
			}
			res := param.Update()
			if res.Code != 0 {
				return errors.New(res.Message)
			}
		}
	}

	fmt.Println("★★★★★所有合同的汇率修改完成......★★★★★")
	return nil
}
