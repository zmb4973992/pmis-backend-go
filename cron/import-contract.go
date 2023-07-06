package cron

import (
	"fmt"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/service"
	"time"
)

type tabContract struct {
	Code                      string  `gorm:"column:F6100"`
	Name                      string  `gorm:"column:F6099"`
	Amount                    float64 `gorm:"column:F5031"`
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

func ImportContract() {
	var records []tabContract
	//主合同的定义是项目
	global.DB2.Table("tabContract").Where("F6110 != '主合同'").
		Find(&records)

	for i := range records {
		var organizationID int64
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
			}

			err := global.DB.Model(&model.Organization{}).
				Where("name = ?", records[i].Organization).Select("id").
				First(&organizationID).Error
			if err != nil {
				param := service.ErrorLogCreate{
					Detail: "tabContract视图的记录中发现无法匹配的部门：" +
						records[i].Organization,
					Date: time.Now().Format("2006-01-02"),
				}
				param.Create()
				records[i].Organization = ""
			}
		}

		var typeID int64
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

			var dictionaryTypeID int64
			err := global.DB.Model(&model.DictionaryType{}).
				Where("name = ?", "合同类型").Select("id").
				First(&dictionaryTypeID).Error
			if err != nil {
				param := service.ErrorLogCreate{
					Detail: "dictionaryType表中找不到”合同类型“这个名称",
					Date:   time.Now().Format("2006-01-02"),
				}
				param.Create()
				records[i].Type = ""
			} else {
				err = global.DB.Model(&model.DictionaryDetail{}).
					Where("dictionary_type_id = ?", dictionaryTypeID).
					Where("name = ?", records[i].Type).Select("id").
					First(&typeID).Error
				if err != nil {
					param := service.ErrorLogCreate{
						Detail: "tabContract视图的记录中发现无法匹配的合同类型：" +
							records[i].Type,
						Date: time.Now().Format("2006-01-02"),
					}
					param.Create()
					records[i].Type = ""
				}
			}

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
						Detail: "tabContract视图的记录中发现无法匹配的币种：" +
							records[i].Currency,
						Date: time.Now().Format("2006-01-02"),
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
					Detail: "tabContract视图的记录中发现无法匹配的项目编号：" +
						records[i].ProjectCode,
					Date: time.Now().Format("2006-01-02"),
				}
				param.Create()
				records[i].ProjectCode = ""
			}
		}

		var relatedPartyID int64
		if records[i].RelatedParty != "" {
			err := global.DB.Model(&model.RelatedParty{}).
				Where("name = ?", records[i].RelatedParty).Select("id").
				First(&relatedPartyID).Error
			if err != nil {
				fmt.Println(333)
				err = global.DB.Model(&model.RelatedParty{}).
					Where("imported_original_name like ?", "%"+records[i].RelatedParty+"%").
					Select("id").
					First(&relatedPartyID).Error
				if err != nil {
					param := service.ErrorLogCreate{
						Detail: "tabContract视图的记录中发现无法匹配的相关方名称：" +
							records[i].RelatedParty,
						Date: time.Now().Format("2006-01-02"),
					}
					param.Create()
					records[i].RelatedParty = ""
				}
			}
		}

		//newRecord := service.ContractCreate{
		//	ProjectID:          projectID,
		//	OrganizationID:     organizationID,
		//	RelatedPartyID:     relatedPartyID,
		//	FundDirection:      0,
		//	OurSignatory:       0,
		//	Currency:           0,
		//	Type:               0,
		//	SigningDate:        "",
		//	EffectiveDate:      "",
		//	CommissioningDate:  "",
		//	CompletionDate:     "",
		//	Amount:             nil,
		//	ExchangeRate:       nil,
		//	ConstructionPeriod: nil,
		//	Name:               "",
		//	Code:               "",
		//	Content:            "",
		//	Deliverable:        "",
		//	PenaltyRule:        "",
		//	Attachment:         "",
		//	Operator:           "",
		//}
		//
		//var count int64
		//global.DB.Model(&model.Project{}).
		//	Where("code = ?", records[i].Code).
		//	Count(&count)
		//if count == 0 {
		//	newRecord.Create()
		//}
		//
		var contracts []model.Contract
		global.DB.Find(&contracts)
		for j := range contracts {
			if contracts[j].Currency != nil && contracts[j].ExchangeRate == nil {
				var currencyName string
				err := global.DB.Model(&model.DictionaryDetail{}).
					Where("id = ?", *contracts[j].Currency).
					Select("name").First(&currencyName).Error
				if err != nil {
					continue
				}

				switch currencyName {
				case "人民币":
					param := service.ProjectUpdate{
						ID:           contracts[j].ID,
						ExchangeRate: model.Float64ToPointer(1),
					}
					param.Update()

				case "美元":
					param := service.ProjectUpdate{ID: contracts[j].ID,
						ExchangeRate: model.Float64ToPointer(7.2),
					}
					param.Update()

				case "欧元":
					param := service.ProjectUpdate{ID: contracts[j].ID,
						ExchangeRate: model.Float64ToPointer(7.8),
					}
					param.Update()

				case "港币":
					param := service.ProjectUpdate{ID: contracts[j].ID,
						ExchangeRate: model.Float64ToPointer(0.92),
					}
					param.Update()

				case "新加坡元":
					param := service.ProjectUpdate{ID: contracts[j].ID,
						ExchangeRate: model.Float64ToPointer(5.3),
					}
					param.Update()

				case "马来西亚币":
					param := service.ProjectUpdate{ID: contracts[j].ID,
						ExchangeRate: model.Float64ToPointer(1.5),
					}
					param.Update()
				}
			}

		}
	}
}
