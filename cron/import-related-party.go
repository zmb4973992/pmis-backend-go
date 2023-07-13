package cron

import (
	"fmt"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/service"
	"pmis-backend-go/util"
	"strconv"
	"strings"
)

func importRelatedParty() {
	fmt.Println("★★★★★开始处理相关方记录......★★★★★")
	importRelatedPartyFromTabSupplier()
	importRelatedPartyFromTabContract()
	importRelatedPartyFromTabFukuan2()
	importRelatedPartyFromTabShouKuan()
	importRelatedPartyFromTabShouHui()

}

type tabSupplier struct {
	Name                    string `gorm:"column:F10627"`
	Address                 string `gorm:"column:F10628"`
	UniformSocialCreditCode string `gorm:"column:F10632"`
}

func importRelatedPartyFromTabSupplier() {
	fmt.Println("正在从tabSupplier导入相关方数据......")
	var records []tabSupplier
	global.DB2.Table("tabSupplier").Find(&records)

	var existedNames []string

	for i := range records {
		if i > 0 && i%100 == 0 {
			process, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(i)/float64(len(records))), 64)
			fmt.Println("已处理", i, "条相关方记录，当前进度：", fmt.Sprintf("%.0f", process*100), "%")
		}

		//初筛，基本能过滤掉95%以上的重复数据
		var tempCount int64
		global.DB.Model(&model.RelatedParty{}).
			Where("name = ?", strings.TrimSpace(records[i].Name)).
			Count(&tempCount)

		//如果通过初筛、没有重复记录，才执行细筛
		if tempCount == 0 {
			var relatedParties []model.RelatedParty
			global.DB.Model(&model.RelatedParty{}).Find(&relatedParties)

			for j := range relatedParties {
				if relatedParties[j].Name != nil {
					existedNames = append(existedNames, *relatedParties[j].Name)
				}
				if relatedParties[j].EnglishName != nil {
					existedNames = append(existedNames, *relatedParties[j].EnglishName)
				}
				if relatedParties[j].ImportedOriginalName != nil {
					importedOriginalNames := strings.Split(*relatedParties[j].ImportedOriginalName, "|")
					existedNames = append(existedNames, importedOriginalNames...)
				}
			}

			if util.SliceIncludes(existedNames, strings.TrimSpace(records[i].Name)) {
				continue
			}

			param := service.RelatedPartyCreate{
				Name:                    strings.TrimSpace(records[i].Name),
				Address:                 records[i].Address,
				UniformSocialCreditCode: records[i].UniformSocialCreditCode,
				ImportedOriginalName:    records[i].Name + "|",
			}
			param.Create()
		}

	}
}

type tabContract2 struct {
	Name string `gorm:"column:F6102"`
}

func importRelatedPartyFromTabContract() {
	fmt.Println("正在从tabContract导入相关方数据......")

	var records []tabContract2
	global.DB2.Table("tabContract").Find(&records)

	var existedNames []string

	for i := range records {
		if i > 0 && i%100 == 0 {
			process, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(i)/float64(len(records))), 64)
			fmt.Println("已处理", i, "条相关方记录，当前进度：", fmt.Sprintf("%.0f", process*100), "%")
		}

		//初筛，基本能过滤掉95%以上的重复数据
		var tempCount int64
		global.DB.Model(&model.RelatedParty{}).
			Where("name = ?", strings.TrimSpace(records[i].Name)).
			Count(&tempCount)

		//如果通过初筛、没有重复记录，才执行细筛
		if tempCount == 0 {
			var relatedParties []model.RelatedParty
			global.DB.Model(&model.RelatedParty{}).Find(&relatedParties)

			for j := range relatedParties {
				if relatedParties[j].Name != nil {
					existedNames = append(existedNames, *relatedParties[j].Name)
				}
				if relatedParties[j].EnglishName != nil {
					existedNames = append(existedNames, *relatedParties[j].EnglishName)
				}
				if relatedParties[j].ImportedOriginalName != nil {
					importedOriginalNames := strings.Split(*relatedParties[j].ImportedOriginalName, "|")
					existedNames = append(existedNames, importedOriginalNames...)
				}
			}

			if util.SliceIncludes(existedNames, strings.TrimSpace(records[i].Name)) {
				continue
			}

			param := service.RelatedPartyCreate{
				Name:                 strings.TrimSpace(records[i].Name),
				ImportedOriginalName: records[i].Name + "|",
			}
			param.Create()
		}
	}
}

type tabFukuan2 struct {
	Name string `gorm:"column:F13591"`
}

func importRelatedPartyFromTabFukuan2() {
	fmt.Println("正在从tabFukuan2导入相关方数据......")

	var records []tabFukuan2
	global.DB2.Table("tabFukuan2").Find(&records)

	var existedNames []string

	for i := range records {
		if i > 0 && i%100 == 0 {
			process, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(i)/float64(len(records))), 64)
			fmt.Println("已处理", i, "条相关方记录，当前进度：", fmt.Sprintf("%.0f", process*100), "%")
		}

		//初筛，基本能过滤掉95%以上的重复数据
		var tempCount int64
		global.DB.Model(&model.RelatedParty{}).
			Where("name = ?", strings.TrimSpace(records[i].Name)).
			Count(&tempCount)

		//如果通过初筛、没有重复记录，才执行细筛
		if tempCount == 0 {
			var relatedParties []model.RelatedParty
			global.DB.Model(&model.RelatedParty{}).Find(&relatedParties)

			for j := range relatedParties {
				if relatedParties[j].Name != nil {
					existedNames = append(existedNames, *relatedParties[j].Name)
				}
				if relatedParties[j].EnglishName != nil {
					existedNames = append(existedNames, *relatedParties[j].EnglishName)
				}
				if relatedParties[j].ImportedOriginalName != nil {
					importedOriginalNames := strings.Split(*relatedParties[j].ImportedOriginalName, "|")
					existedNames = append(existedNames, importedOriginalNames...)
				}
			}

			if util.SliceIncludes(existedNames, strings.TrimSpace(records[i].Name)) {
				continue
			}

			param := service.RelatedPartyCreate{
				Name:                 strings.TrimSpace(records[i].Name),
				ImportedOriginalName: records[i].Name + "|"}
			param.Create()
		}
	}
}

type tabShouKuanA struct {
	RelatedPartyName string `gorm:"column:F10851"`
}

func importRelatedPartyFromTabShouKuan() {
	fmt.Println("正在从tabShouKuan导入相关方数据......")

	var records []tabShouKuanA
	global.DB2.Table("tabShouKuan").Find(&records)

	var existedNames []string

	for i := range records {
		if i > 0 && i%100 == 0 {
			process, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(i)/float64(len(records))), 64)
			fmt.Println("已处理", i, "条相关方记录，当前进度：", fmt.Sprintf("%.0f", process*100), "%")
		}

		//初筛，基本能过滤掉95%以上的重复数据
		var tempCount int64
		global.DB.Model(&model.RelatedParty{}).
			Where("name = ?", strings.TrimSpace(records[i].RelatedPartyName)).
			Count(&tempCount)

		//如果通过初筛、没有重复记录，才执行细筛
		if tempCount == 0 {
			var relatedParties []model.RelatedParty
			global.DB.Model(&model.RelatedParty{}).Find(&relatedParties)

			for j := range relatedParties {
				if relatedParties[j].Name != nil {
					existedNames = append(existedNames, *relatedParties[j].Name)
				}
				if relatedParties[j].EnglishName != nil {
					existedNames = append(existedNames, *relatedParties[j].EnglishName)
				}
				if relatedParties[j].ImportedOriginalName != nil {
					importedOriginalNames := strings.Split(*relatedParties[j].ImportedOriginalName, "|")
					existedNames = append(existedNames, importedOriginalNames...)
				}
			}

			if util.SliceIncludes(existedNames, strings.TrimSpace(records[i].RelatedPartyName)) {
				continue
			}

			param := service.RelatedPartyCreate{
				Name:                 strings.TrimSpace(records[i].RelatedPartyName),
				ImportedOriginalName: records[i].RelatedPartyName + "|"}
			param.Create()
		}
	}
}

type tabShouHuiA struct {
	RelatedPartyName string `gorm:"column:F14394"`
}

func importRelatedPartyFromTabShouHui() {
	fmt.Println("正在从tabShouHui导入相关方数据......")

	var records []tabShouHuiA
	global.DB2.Table("tabShouHui").Find(&records)

	var existedNames []string

	for i := range records {
		if i > 0 && i%100 == 0 {
			process, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(i)/float64(len(records))), 64)
			fmt.Println("已处理", i, "条相关方记录，当前进度：", fmt.Sprintf("%.0f", process*100), "%")
		}

		//初筛，基本能过滤掉95%以上的重复数据
		var tempCount int64
		global.DB.Model(&model.RelatedParty{}).
			Where("name = ?", strings.TrimSpace(records[i].RelatedPartyName)).
			Count(&tempCount)

		//如果通过初筛、没有重复记录，才执行细筛
		if tempCount == 0 {
			var relatedParties []model.RelatedParty
			global.DB.Model(&model.RelatedParty{}).Find(&relatedParties)

			for j := range relatedParties {
				if relatedParties[j].Name != nil {
					existedNames = append(existedNames, *relatedParties[j].Name)
				}
				if relatedParties[j].EnglishName != nil {
					existedNames = append(existedNames, *relatedParties[j].EnglishName)
				}
				if relatedParties[j].ImportedOriginalName != nil {
					importedOriginalNames := strings.Split(*relatedParties[j].ImportedOriginalName, "|")
					existedNames = append(existedNames, importedOriginalNames...)
				}
			}

			if util.SliceIncludes(existedNames, strings.TrimSpace(records[i].RelatedPartyName)) {
				continue
			}

			param := service.RelatedPartyCreate{
				Name:                 strings.TrimSpace(records[i].RelatedPartyName),
				ImportedOriginalName: records[i].RelatedPartyName + "|"}
			param.Create()
		}
	}
}
