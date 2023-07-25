package cron

import (
	"errors"
	"fmt"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/service"
	"strconv"
	"time"
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

func importProject() error {
	fmt.Println("★★★★★开始处理项目记录......★★★★★")

	var records []tabProject
	global.DB2.Table("tabProject").Find(&records)

	for i := range records {
		if i > 0 && i%100 == 0 {
			process, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(i)/float64(len(records))), 64)
			fmt.Println("已处理", i, "条项目记录，当前进度：", fmt.Sprintf("%.0f", process*100), "%")
		}

		var tempCount int64
		global.DB.Model(&model.Project{}).
			Where("code = ?", records[i].Code).
			Count(&tempCount)
		if tempCount == 0 {
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
				}

				err := global.DB.Model(&model.Organization{}).
					Where("name = ?", records[i].Organization).Select("id").
					First(&organizationID).Error
				if err != nil {
					param := service.ErrorLogCreate{
						Detail: "tabProject视图的记录中发现无法匹配的部门：" +
							records[i].Organization,
						Date: time.Now().Format("2006-01-02"),
					}
					param.Create()
					records[i].Organization = ""
				}
			}

			var countryID int64
			if records[i].Country != "" {
				switch records[i].Country {
				case "AC":
					records[i].Country = "英国" //阿森松岛
				case "AD":
					records[i].Country = "安道尔"
				case "AE":
					records[i].Country = "阿联酋"
				case "AF":
					records[i].Country = "阿富汗"
				case "AL":
					records[i].Country = "阿尔巴尼亚"
				case "AM":
					records[i].Country = "亚美尼亚"
				case "AO":
					records[i].Country = "安哥拉"
				case "AR":
					records[i].Country = "阿根廷"
				case "AT":
					records[i].Country = "奥地利"
				case "AU":
					records[i].Country = "澳大利亚"
				case "AZ":
					records[i].Country = "阿塞拜疆"

				case "BA":
					records[i].Country = "波黑"
				case "BD":
					records[i].Country = "孟加拉"
				case "BE":
					records[i].Country = "比利时"
				case "BF":
					records[i].Country = "布基纳法索"
				case "BG":
					records[i].Country = "保加利亚"
				case "BH":
					records[i].Country = "巴林"
				case "BJ":
					records[i].Country = "贝宁"
				case "BO":
					records[i].Country = "玻利维亚"
				case "BR":
					records[i].Country = "巴西"
				case "BT":
					records[i].Country = "不丹"
				case "BY":
					records[i].Country = "白俄罗斯"

				case "CA":
					records[i].Country = "加拿大"
				case "CD":
					records[i].Country = "刚果民主共和国"
				case "CF":
					records[i].Country = "中非"
				case "CG":
					records[i].Country = "刚果共和国"
				case "CI":
					records[i].Country = "科特迪瓦"
				case "CL":
					records[i].Country = "智利"
				case "CM":
					records[i].Country = "喀麦隆"
				case "CN":
					records[i].Country = "智利"
				case "CO":
					records[i].Country = "哥伦比亚"
				case "CR":
					records[i].Country = "哥斯达黎加"
				case "CU":
					records[i].Country = "古巴"
				case "CY":
					records[i].Country = "塞浦路斯"
				case "CZ":
					records[i].Country = "捷克"

				case "DE":
					records[i].Country = "德国"
				case "DJ":
					records[i].Country = "吉布提"
				case "DK":
					records[i].Country = "丹麦"
				case "DO":
					records[i].Country = "多米尼加"
				case "DZ":
					records[i].Country = "阿尔及利亚"

				case "EC":
					records[i].Country = "厄瓜多尔"
				case "EG":
					records[i].Country = "埃及"
				case "ES":
					records[i].Country = "西班牙"
				case "ET":
					records[i].Country = "埃塞俄比亚"

				case "FR":
					records[i].Country = "法国"

				case "GA":
					records[i].Country = "加蓬"
				case "GH":
					records[i].Country = "加纳"
				case "GJ":
					records[i].Country = "格鲁吉亚"
				case "GM":
					records[i].Country = "冈比亚"
				case "GN":
					records[i].Country = "几内亚"
				case "GQ":
					records[i].Country = "赤道几内亚"
				case "GR":
					records[i].Country = "希腊"
				case "GT":
					records[i].Country = "危地马拉"
				case "GW":
					records[i].Country = "几内亚比绍"

				case "HK":
					records[i].Country = "中国"
				case "HR":
					records[i].Country = "克罗地亚"
				case "HU":
					records[i].Country = "匈牙利"

				case "ID":
					records[i].Country = "印度尼西亚"
				case "IL":
					records[i].Country = "以色列"
				case "IN":
					records[i].Country = "印度"
				case "IQ":
					records[i].Country = "伊拉克"
				case "IR":
					records[i].Country = "伊朗"
				case "IT":
					records[i].Country = "意大利"

				case "JO":
					records[i].Country = "约旦"

				case "KE":
					records[i].Country = "肯尼亚"
				case "KG":
					records[i].Country = "吉尔吉斯斯坦"
				case "KH":
					records[i].Country = "柬埔寨"
				case "KR":
					records[i].Country = "韩国"
				case "KZ":
					records[i].Country = "哈萨克斯坦"

				case "LA":
					records[i].Country = "老挝"
				case "LR":
					records[i].Country = "利比里亚"
				case "LY":
					records[i].Country = "利比亚"

				case "MA":
					records[i].Country = "摩洛哥"
				case "MM":
					records[i].Country = "缅甸"
				case "MN":
					records[i].Country = "蒙古"
				case "MX":
					records[i].Country = "墨西哥"
				case "MY":
					records[i].Country = "马来西亚"
				case "MZ":
					records[i].Country = "莫桑比克"

				case "NG":
					records[i].Country = "尼日利亚"
				case "NP":
					records[i].Country = "尼泊尔"
				case "NZ":
					records[i].Country = "新西兰"

				case "OM":
					records[i].Country = "阿曼"

				case "PE":
					records[i].Country = "秘鲁"
				case "PH":
					records[i].Country = "菲律宾"
				case "PK":
					records[i].Country = "巴基斯坦"
				case "PL":
					records[i].Country = "波兰"

				case "RU":
					records[i].Country = "俄罗斯"

				case "SA":
					records[i].Country = "沙特阿拉伯"
				case "SD":
					records[i].Country = "苏丹"
				case "SG":
					records[i].Country = "新加坡"
				case "SN":
					records[i].Country = "塞内加尔"
				case "SO":
					records[i].Country = "索马里"
				case "SY":
					records[i].Country = "叙利亚"

				case "TD":
					records[i].Country = "乍得"
				case "TH":
					records[i].Country = "泰国"
				case "TJ":
					records[i].Country = "塔吉克斯坦"
				case "TR":
					records[i].Country = "土耳其"
				case "TZ":
					records[i].Country = "坦桑尼亚"

				case "UA":
					records[i].Country = "乌克兰"
				case "US":
					records[i].Country = "美国"
				case "UZ":
					records[i].Country = "乌兹别克斯坦"

				case "VE":
					records[i].Country = "委内瑞拉"
				case "VN":
					records[i].Country = "越南"

				case "ZA":
					records[i].Country = "南非"
				case "ZM":
					records[i].Country = "赞比亚"
				case "ZW":
					records[i].Country = "津巴布韦"
				}

				err := global.DB.Model(&model.DictionaryDetail{}).
					Where("name = ?", records[i].Country).Select("id").
					First(&countryID).Error
				if err != nil {
					param := service.ErrorLogCreate{
						Detail: "tabProject视图的记录中发现无法匹配的国别:" +
							records[i].Country,
						Date: time.Now().Format("2006-01-02"),
					}
					param.Create()
					records[i].Country = ""
				}
			}

			var typeID int64
			if records[i].Type != "" {
				switch records[i].Type {
				case "C":
					records[i].Type = "工程EPC"
				case "T":
					records[i].Type = "贸易"
				case "S":
					records[i].Type = "服务"
				case "I":
					records[i].Type = "投资"
				case "O":
					records[i].Type = "其他"
				}

				err := global.DB.Model(&model.DictionaryDetail{}).
					Where("name = ?", records[i].Type).Select("id").
					First(&typeID).Error
				if err != nil {
					param := service.ErrorLogCreate{
						Detail: "tabProject视图的记录中发现无法匹配的项目类型：" +
							records[i].Type,
						Date: time.Now().Format("2006-01-02"),
					}
					param.Create()
					records[i].Type = ""
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

				err := global.DB.Model(&model.DictionaryDetail{}).
					Where("name = ?", records[i].Currency).Select("id").
					First(&currencyID).Error
				if err != nil {
					param := service.ErrorLogCreate{
						Detail: "tabProject视图的记录中发现无法匹配的币种：" +
							records[i].Currency,
						Date: time.Now().Format("2006-01-02"),
					}
					param.Create()
					records[i].Currency = ""
				}
			}

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

				err := global.DB.Model(&model.DictionaryDetail{}).
					Where("name = ?", records[i].Currency).Select("id").
					First(&currencyID).Error
				if err != nil {
					param := service.ErrorLogCreate{
						Detail: "tabProject视图的记录中发现无法匹配的币种：" +
							records[i].Currency,
						Date: time.Now().Format("2006-01-02"),
					}
					param.Create()
					records[i].Currency = ""
				}
			}

			newRecord := service.ProjectCreate{
				OrganizationID:     organizationID,
				RelatedPartyID:     0,
				Country:            countryID,
				Type:               typeID,
				DetailedType:       0,
				Currency:           currencyID,
				Status:             0,
				OurSignatory:       0,
				SigningDate:        "",
				EffectiveDate:      "",
				CommissioningDate:  "",
				Amount:             &records[i].Amount,
				ExchangeRate:       nil,
				ConstructionPeriod: nil,
				Code:               records[i].Code,
				Name:               records[i].Name,
				Content:            "",
			}

			var count int64
			global.DB.Model(&model.Project{}).
				Where("code = ?", records[i].Code).
				Count(&count)
			if count == 0 {
				res := newRecord.Create()
				if res.Code != 0 {
					return errors.New(res.Message)
				}
			}
		}
	}

	var projects []model.Project
	global.DB.Find(&projects)
	for j := range projects {
		if projects[j].Currency != nil && projects[j].ExchangeRate == nil {
			var currencyName string
			err := global.DB.Model(&model.DictionaryDetail{}).
				Where("id = ?", *projects[j].Currency).
				Select("name").First(&currencyName).Error
			if err != nil {
				continue
			}

			switch currencyName {
			case "人民币":
				param := service.ProjectUpdate{
					ID:           projects[j].ID,
					ExchangeRate: model.Float64ToPointer(1),
				}
				res := param.Update()
				if res.Code != 0 {
					return errors.New(res.Message)
				}

			case "美元":
				param := service.ProjectUpdate{
					ID:           projects[j].ID,
					ExchangeRate: model.Float64ToPointer(7.2),
				}
				res := param.Update()
				if res.Code != 0 {
					return errors.New(res.Message)
				}

			case "欧元":
				param := service.ProjectUpdate{
					ID:           projects[j].ID,
					ExchangeRate: model.Float64ToPointer(7.8),
				}
				res := param.Update()
				if res.Code != 0 {
					return errors.New(res.Message)
				}

			case "港币":
				param := service.ProjectUpdate{
					ID:           projects[j].ID,
					ExchangeRate: model.Float64ToPointer(0.92),
				}
				res := param.Update()
				if res.Code != 0 {
					return errors.New(res.Message)
				}

			case "新加坡元":
				param := service.ProjectUpdate{
					ID:           projects[j].ID,
					ExchangeRate: model.Float64ToPointer(5.3),
				}
				res := param.Update()
				if res.Code != 0 {
					return errors.New(res.Message)
				}
			}
		}
	}

	return nil
}
