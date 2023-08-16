package lvmin

import (
	"errors"
	"fmt"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/service"
	"strconv"
)

type tabProject struct {
	Code         string  `gorm:"column:F769"`
	Name         string  `gorm:"column:F770"`
	Organization string  `gorm:"column:F3937"`
	Type         string  `gorm:"column:F5029"`
	Country      string  `gorm:"column:F5030"`
	Amount       float64 `gorm:"column:F5031"`
	Currency     string  `gorm:"column:F5032"`
}

func ImportProject(userID int64) error {
	fmt.Println("★★★★★开始导入项目记录......★★★★★")

	var projects []tabProject
	global.DBForLvmin.Table("tabProject").Find(&projects)

	var country model.DictionaryType
	err := global.DB.
		Where("name = ?", "国家").
		First(&country).Error
	if err != nil {
		param := service.ErrorLogCreate{
			Detail: "dictionaryType表中找不到”国家“这个名称",
		}
		param.Create()
		return err
	}

	var projectType model.DictionaryType
	err = global.DB.
		Where("name = ?", "项目类型").
		First(&projectType).Error
	if err != nil {
		param := service.ErrorLogCreate{
			Detail: "dictionaryType表中找不到”项目类型“这个名称",
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

	for i := range projects {
		if i > 0 && i%1000 == 0 {
			process, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(i)/float64(len(projects))), 64)
			fmt.Println("已处理", i, "条项目记录，当前进度：", fmt.Sprintf("%.0f", process*100), "%")
		} else if i == len(projects)-1 {
			fmt.Println("已处理", i, "条项目记录，当前进度：100 %")
		}

		var tempCount int64
		global.DB.Model(&model.Project{}).
			Where("code = ?", projects[i].Code).
			Count(&tempCount)
		if tempCount == 0 {
			var organization model.Organization
			if projects[i].Organization != "" {
				switch projects[i].Organization {
				case "机械车辆部":
					projects[i].Organization = "成套业务一部"
				case "成套六部":
					projects[i].Organization = "成套业务六部"
				case "事业部管理委员会和水泥工程事业部":
					projects[i].Organization = "水泥工程事业部"
				case "项目管理及技术支持部":
					projects[i].Organization = "项目管理部"
				case "综合管理部":
					projects[i].Organization = "综合管理和法律部"
				case "党建和纪检审计部":
					projects[i].Organization = "党建文宣部"
				case "成套业务三部":
					projects[i].Organization = "成套业务四部"
				}

				err = global.DB.
					Where("name = ?", projects[i].Organization).
					First(&organization).Error
				if err != nil {
					param := service.ErrorLogCreate{
						Detail: "tabProject视图的记录中发现无法匹配的部门：" +
							projects[i].Organization,
					}
					param.Create()
					projects[i].Organization = ""
				}
			}

			var specificCountry model.DictionaryDetail
			if projects[i].Country != "" {
				switch projects[i].Country {
				case "AC":
					projects[i].Country = "英国" //阿森松岛
				case "AD":
					projects[i].Country = "安道尔"
				case "AE":
					projects[i].Country = "阿联酋"
				case "AF":
					projects[i].Country = "阿富汗"
				case "AL":
					projects[i].Country = "阿尔巴尼亚"
				case "AM":
					projects[i].Country = "亚美尼亚"
				case "AO":
					projects[i].Country = "安哥拉"
				case "AR":
					projects[i].Country = "阿根廷"
				case "AT":
					projects[i].Country = "奥地利"
				case "AU":
					projects[i].Country = "澳大利亚"
				case "AZ":
					projects[i].Country = "阿塞拜疆"

				case "BA":
					projects[i].Country = "波黑"
				case "BD":
					projects[i].Country = "孟加拉"
				case "BE":
					projects[i].Country = "比利时"
				case "BF":
					projects[i].Country = "布基纳法索"
				case "BG":
					projects[i].Country = "保加利亚"
				case "BH":
					projects[i].Country = "巴林"
				case "BI":
					projects[i].Country = "布隆迪"
				case "BJ":
					projects[i].Country = "贝宁"
				case "BO":
					projects[i].Country = "玻利维亚"
				case "BR":
					projects[i].Country = "巴西"
				case "BT":
					projects[i].Country = "不丹"
				case "BW":
					projects[i].Country = "博茨瓦纳"
				case "BY":
					projects[i].Country = "白俄罗斯"

				case "CA":
					projects[i].Country = "加拿大"
				case "CD":
					projects[i].Country = "刚果民主共和国"
				case "CF":
					projects[i].Country = "中非"
				case "CG":
					projects[i].Country = "刚果共和国"
				case "CH":
					projects[i].Country = "瑞士"
				case "CI":
					projects[i].Country = "科特迪瓦"
				case "CL":
					projects[i].Country = "智利"
				case "CM":
					projects[i].Country = "喀麦隆"
				case "CN":
					projects[i].Country = "中国"
				case "CO":
					projects[i].Country = "哥伦比亚"
				case "CR":
					projects[i].Country = "哥斯达黎加"
				case "CU":
					projects[i].Country = "古巴"
				case "CY":
					projects[i].Country = "塞浦路斯"
				case "CZ":
					projects[i].Country = "捷克"

				case "DE":
					projects[i].Country = "德国"
				case "DJ":
					projects[i].Country = "吉布提"
				case "DK":
					projects[i].Country = "丹麦"
				case "DM": //两个都是多米尼加
					projects[i].Country = "多米尼加"
				case "DO": //两个都是多米尼加
					projects[i].Country = "多米尼加"
				case "DZ":
					projects[i].Country = "阿尔及利亚"

				case "EC":
					projects[i].Country = "厄瓜多尔"
				case "EE":
					projects[i].Country = "爱沙尼亚"
				case "EG":
					projects[i].Country = "埃及"
				case "ER":
					projects[i].Country = "厄立特里亚"
				case "ES":
					projects[i].Country = "西班牙"
				case "ET":
					projects[i].Country = "埃塞俄比亚"

				case "FI":
					projects[i].Country = "芬兰"
				case "FJ":
					projects[i].Country = "斐济"
				case "FR":
					projects[i].Country = "法国"

				case "GA":
					projects[i].Country = "加蓬"
				case "GB":
					projects[i].Country = "英国"
				case "GH":
					projects[i].Country = "加纳"
				case "GJ":
					projects[i].Country = "格鲁吉亚"
				case "GM":
					projects[i].Country = "冈比亚"
				case "GN":
					projects[i].Country = "几内亚"
				case "GQ":
					projects[i].Country = "赤道几内亚"
				case "GR":
					projects[i].Country = "希腊"
				case "GT":
					projects[i].Country = "危地马拉"
				case "GW":
					projects[i].Country = "几内亚比绍"
				case "GY":
					projects[i].Country = "圭亚那"

				case "HK":
					projects[i].Country = "中国"
				case "HN":
					projects[i].Country = "洪都拉斯"
				case "HR":
					projects[i].Country = "克罗地亚"
				case "HT":
					projects[i].Country = "海地"
				case "HU":
					projects[i].Country = "匈牙利"

				case "ID":
					projects[i].Country = "印度尼西亚"
				case "IE":
					projects[i].Country = "爱尔兰"
				case "IL":
					projects[i].Country = "以色列"
				case "IN":
					projects[i].Country = "印度"
				case "IQ":
					projects[i].Country = "伊拉克"
				case "IR":
					projects[i].Country = "伊朗"
				case "IS":
					projects[i].Country = "冰岛"
				case "IT":
					projects[i].Country = "意大利"

				case "JM":
					projects[i].Country = "牙买加"
				case "JO":
					projects[i].Country = "约旦"
				case "JP":
					projects[i].Country = "日本"

				case "KE":
					projects[i].Country = "肯尼亚"
				case "KG":
					projects[i].Country = "吉尔吉斯斯坦"
				case "KH":
					projects[i].Country = "柬埔寨"
				case "KM":
					projects[i].Country = "科摩罗"
				case "KP":
					projects[i].Country = "朝鲜"
				case "KR":
					projects[i].Country = "韩国"
				case "KW":
					projects[i].Country = "科威特"
				case "KZ":
					projects[i].Country = "哈萨克斯坦"

				case "LA":
					projects[i].Country = "老挝"
				case "LB":
					projects[i].Country = "黎巴嫩"
				case "LI":
					projects[i].Country = "列支敦士登"
				case "LK":
					projects[i].Country = "斯里兰卡"
				case "LR":
					projects[i].Country = "利比里亚"
				case "LS":
					projects[i].Country = "莱索托"
				case "LT":
					projects[i].Country = "立陶宛"
				case "LU":
					projects[i].Country = "卢森堡"
				case "LV":
					projects[i].Country = "拉脱维亚"
				case "LY":
					projects[i].Country = "利比亚"

				case "MA":
					projects[i].Country = "摩洛哥"
				case "MC":
					projects[i].Country = "摩纳哥"
				case "MD":
					projects[i].Country = "摩尔多瓦"
				case "ME":
					projects[i].Country = "黑山"
				case "MG":
					projects[i].Country = "马达加斯加"
				case "MH":
					projects[i].Country = "马绍尔"
				case "MK":
					projects[i].Country = "马其顿"
				case "ML":
					projects[i].Country = "马里"
				case "MM":
					projects[i].Country = "缅甸"
				case "MN":
					projects[i].Country = "蒙古"
				case "MR":
					projects[i].Country = "毛里塔尼亚"
				case "MT":
					projects[i].Country = "马耳他"
				case "MU":
					projects[i].Country = "毛里求斯"
				case "MV":
					projects[i].Country = "马尔代夫"
				case "MW":
					projects[i].Country = "马拉维"
				case "MX":
					projects[i].Country = "墨西哥"
				case "MY":
					projects[i].Country = "马来西亚"
				case "MZ":
					projects[i].Country = "莫桑比克"

				case "NE":
					projects[i].Country = "尼日尔"
				case "NG":
					projects[i].Country = "尼日利亚"
				case "NI":
					projects[i].Country = "尼加拉瓜"
				case "NL":
					projects[i].Country = "荷兰"
				case "NO":
					projects[i].Country = "挪威"
				case "NP":
					projects[i].Country = "尼泊尔"
				case "NR":
					projects[i].Country = "尼泊尔"
				case "NZ":
					projects[i].Country = "新西兰"

				case "OM":
					projects[i].Country = "阿曼"

				case "PA":
					projects[i].Country = "巴拿马"
				case "PE":
					projects[i].Country = "秘鲁"
				case "PG":
					projects[i].Country = "巴布亚新几内亚"
				case "PH":
					projects[i].Country = "菲律宾"
				case "PK":
					projects[i].Country = "巴基斯坦"
				case "PL":
					projects[i].Country = "波兰"
				case "PR":
					projects[i].Country = "波多黎各"
				case "PT":
					projects[i].Country = "葡萄牙"
				case "PW":
					projects[i].Country = "帕劳"
				case "PY":
					projects[i].Country = "巴拉圭"

				case "QA":
					projects[i].Country = "卡塔尔"

				case "RO":
					projects[i].Country = "罗马尼亚"
				case "RS":
					projects[i].Country = "塞尔维亚"
				case "RU":
					projects[i].Country = "俄罗斯"
				case "RW":
					projects[i].Country = "卢旺达"

				case "SA":
					projects[i].Country = "沙特阿拉伯"
				case "SB":
					projects[i].Country = "所罗门群岛"
				case "SC":
					projects[i].Country = "塞舌尔"
				case "SD":
					projects[i].Country = "苏丹"
				case "SE":
					projects[i].Country = "瑞典"
				case "SG":
					projects[i].Country = "新加坡"
				case "SI":
					projects[i].Country = "斯洛文尼亚"
				case "SK":
					projects[i].Country = "斯洛伐克"
				case "SL":
					projects[i].Country = "塞拉利昂"
				case "SM":
					projects[i].Country = "圣马力诺"
				case "SN":
					projects[i].Country = "塞内加尔"
				case "SO":
					projects[i].Country = "索马里"
				case "SR":
					projects[i].Country = "苏里南"
				case "SV":
					projects[i].Country = "萨尔瓦多"
				case "SY":
					projects[i].Country = "叙利亚"
				case "SZ":
					projects[i].Country = "斯威士兰"

				case "TD":
					projects[i].Country = "乍得"
				case "TG":
					projects[i].Country = "多哥"
				case "TH":
					projects[i].Country = "泰国"
				case "TJ":
					projects[i].Country = "塔吉克斯坦"
				case "TL":
					projects[i].Country = "东帝汶"
				case "TM":
					projects[i].Country = "土库曼斯坦"
				case "TN":
					projects[i].Country = "突尼斯"
				case "TO":
					projects[i].Country = "汤加"
				case "TR":
					projects[i].Country = "土耳其"
				case "TT":
					projects[i].Country = "特立尼达和多巴哥"
				case "TV":
					projects[i].Country = "图瓦卢"
				case "TZ":
					projects[i].Country = "坦桑尼亚"

				case "UA":
					projects[i].Country = "乌克兰"
				case "UG":
					projects[i].Country = "乌干达"
				case "UK":
					projects[i].Country = "英国"
				case "US":
					projects[i].Country = "美国"
				case "UY":
					projects[i].Country = "乌拉圭"
				case "UZ":
					projects[i].Country = "乌兹别克斯坦"

				case "VC":
					projects[i].Country = "圣文森特和格林纳丁斯"
				case "VE":
					projects[i].Country = "委内瑞拉"
				case "VN":
					projects[i].Country = "越南"
				case "VU":
					projects[i].Country = "瓦努阿图"

				case "WS":
					projects[i].Country = "萨摩亚"

				case "YE":
					projects[i].Country = "也门"

				case "ZA":
					projects[i].Country = "南非"
				case "ZM":
					projects[i].Country = "赞比亚"
				case "ZW":
					projects[i].Country = "津巴布韦"
				}

				err = global.DB.Model(&model.DictionaryDetail{}).
					Where("dictionary_type_id = ?", country.ID).
					Where("name = ?", projects[i].Country).
					First(&specificCountry).Error
				if err != nil {
					param := service.ErrorLogCreate{
						Detail: "tabProject视图的记录中发现无法匹配的国别:" +
							projects[i].Country,
					}
					param.Create()
					projects[i].Country = ""
				}
			}

			var specificProjectType model.DictionaryDetail
			if projects[i].Type != "" {
				switch projects[i].Type {
				case "C":
					projects[i].Type = "工程EPC"
				case "T":
					projects[i].Type = "贸易"
				case "S":
					projects[i].Type = "服务"
				case "I":
					projects[i].Type = "投资"
				case "O":
					projects[i].Type = "其他"
				}

				err = global.DB.
					Where("dictionary_type_id = ?", projectType.ID).
					Where("name = ?", projects[i].Type).
					First(&specificProjectType).Error
				if err != nil {
					param := service.ErrorLogCreate{
						Detail: "tabProject视图的记录中发现无法匹配的项目类型：" +
							projects[i].Type,
					}
					param.Create()
					projects[i].Type = ""
				}
			}

			var specificCurrency model.DictionaryDetail
			if projects[i].Currency != "" {
				switch projects[i].Currency {
				case "RMB":
					projects[i].Currency = "人民币"
				case "1":
					projects[i].Currency = "人民币"
				case "2":
					projects[i].Currency = "美元"
				case "3":
					projects[i].Currency = "欧元"
				}

				err = global.DB.
					Where("dictionary_type_id = ?", currency.ID).
					Where("name = ?", projects[i].Currency).
					First(&specificCurrency).Error
				if err != nil {
					param := service.ErrorLogCreate{
						Detail: "tabProject视图的记录中发现无法匹配的币种：" +
							projects[i].Currency,
					}
					param.Create()
					projects[i].Currency = ""
				}
			}

			newRecord := service.ProjectCreate{
				UserID:             userID,
				OrganizationID:     organization.ID,
				RelatedPartyID:     0,
				Country:            specificCountry.ID,
				Type:               specificProjectType.ID,
				DetailedType:       0,
				Currency:           specificCurrency.ID,
				Status:             0,
				OurSignatory:       0,
				SigningDate:        "",
				EffectiveDate:      "",
				CommissioningDate:  "",
				Amount:             &projects[i].Amount,
				ExchangeRate:       nil,
				ConstructionPeriod: nil,
				Code:               projects[i].Code,
				Name:               projects[i].Name,
				Content:            "",
			}

			var count int64
			global.DB.Model(&model.Project{}).
				Where("code = ?", projects[i].Code).
				Count(&count)
			if count == 0 {
				res := newRecord.Create()
				if res.Code != 0 {
					return errors.New(res.Message)
				}
			}
		}
	}

	return nil
}

func UpdateExchangeRageOfProject(userID int64) error {
	fmt.Println("★★★★★开始更正所有项目的汇率......★★★★★")

	var projects []model.Project
	global.DB.
		Where("currency is not null").
		Where("exchange_rate is null").
		Find(&projects)

	for i := range projects {
		if i > 0 && i%1000 == 0 {
			process, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(i)/float64(len(projects))), 64)
			fmt.Println("已处理", i, "条项目记录，当前进度：", fmt.Sprintf("%.0f", process*100), "%")
		} else if i == len(projects)-1 {
			fmt.Println("已处理", i, "条项目记录，当前进度：100 %")
		}

		var specificCurrency model.DictionaryDetail
		err := global.DB.
			Where("id = ?", *projects[i].Currency).
			First(&specificCurrency).Error
		if err != nil {
			continue
		}

		switch specificCurrency.Name {
		case "人民币":
			param := service.ProjectUpdate{
				IgnoreDataAuthority: true,
				UserID:              userID,
				ID:                  projects[i].ID,
				ExchangeRate:        model.Float64ToPointer(1),
			}
			res := param.Update()
			if res.Code != 0 {
				return errors.New(res.Message)
			}

		case "美元":
			param := service.ProjectUpdate{
				IgnoreDataAuthority: true,
				UserID:              userID,
				ID:                  projects[i].ID,
				ExchangeRate:        model.Float64ToPointer(7.2),
			}
			res := param.Update()
			if res.Code != 0 {
				return errors.New(res.Message)
			}

		case "欧元":
			param := service.ProjectUpdate{
				IgnoreDataAuthority: true,
				UserID:              userID,
				ID:                  projects[i].ID,
				ExchangeRate:        model.Float64ToPointer(7.8),
			}
			res := param.Update()
			if res.Code != 0 {
				return errors.New(res.Message)
			}

		case "港币":
			param := service.ProjectUpdate{
				IgnoreDataAuthority: true,
				UserID:              userID,
				ID:                  projects[i].ID,
				ExchangeRate:        model.Float64ToPointer(0.92),
			}
			res := param.Update()
			if res.Code != 0 {
				return errors.New(res.Message)
			}

		case "新加坡元":
			param := service.ProjectUpdate{
				IgnoreDataAuthority: true,
				UserID:              userID,
				ID:                  projects[i].ID,
				ExchangeRate:        model.Float64ToPointer(5.3),
			}
			res := param.Update()
			if res.Code != 0 {
				return errors.New(res.Message)
			}
		}
	}

	fmt.Println("★★★★★所有项目的汇率修改完成......★★★★★")
	return nil
}
