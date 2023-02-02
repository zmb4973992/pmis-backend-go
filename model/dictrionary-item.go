package model

import "pmis-backend-go/global"

type DictionaryItem struct {
	BaseModel
	DictionaryTypeID int     //字典类型的id
	Name             string  //名称
	Sort             *int    //用于排序的值
	Remarks          *string //备注
}

// TableName 修改数据库的表名
func (*DictionaryItem) TableName() string {
	return "dictionary_item"
}

func generateDictionaryItems() error {
	var dictionaryItems []DictionaryItem

	var dictionaryTypeIDOfProvince int
	global.DB.Model(&DictionaryType{}).
		Where("name = '省份'").Select("id").First(&dictionaryTypeIDOfProvince)
	dictionaryItemsOfProvince := []DictionaryItem{
		{Name: "北京", DictionaryTypeID: dictionaryTypeIDOfProvince},
		{Name: "上海", DictionaryTypeID: dictionaryTypeIDOfProvince},
		{Name: "山东", DictionaryTypeID: dictionaryTypeIDOfProvince},
	}
	dictionaryItems = append(dictionaryItems, dictionaryItemsOfProvince...)

	var dictionaryTypeIDOfReceiptOrPaymentTerm int
	global.DB.Model(&DictionaryType{}).
		Where("name = '收付款方式'").Select("id").First(&dictionaryTypeIDOfReceiptOrPaymentTerm)
	dictionaryItemsOfReceiptOrPaymentTerm := []DictionaryItem{
		{Name: "现金", DictionaryTypeID: dictionaryTypeIDOfReceiptOrPaymentTerm},
		{Name: "汇票", DictionaryTypeID: dictionaryTypeIDOfReceiptOrPaymentTerm},
		{Name: "信用证", DictionaryTypeID: dictionaryTypeIDOfReceiptOrPaymentTerm},
	}
	dictionaryItems = append(dictionaryItems, dictionaryItemsOfReceiptOrPaymentTerm...)

	var dictionaryTypeIDOfCurrency int
	global.DB.Model(&DictionaryType{}).
		Where("name = '币种'").Select("id").First(&dictionaryTypeIDOfCurrency)
	dictionaryItemsOfCurrency := []DictionaryItem{
		{Name: "人民币", DictionaryTypeID: dictionaryTypeIDOfCurrency},
		{Name: "美元", DictionaryTypeID: dictionaryTypeIDOfCurrency},
		{Name: "欧元", DictionaryTypeID: dictionaryTypeIDOfCurrency},
	}
	dictionaryItems = append(dictionaryItems, dictionaryItemsOfCurrency...)

	var dictionaryTypeIDOfContractType int
	global.DB.Model(&DictionaryType{}).
		Where("name = '合同类型'").Select("id").First(&dictionaryTypeIDOfContractType)
	dictionaryItemsOfContractType := []DictionaryItem{
		{Name: "采购", DictionaryTypeID: dictionaryTypeIDOfContractType},
		{Name: "销售", DictionaryTypeID: dictionaryTypeIDOfContractType},
		{Name: "代理", DictionaryTypeID: dictionaryTypeIDOfContractType},
	}
	dictionaryItems = append(dictionaryItems, dictionaryItemsOfContractType...)

	var dictionaryTypeIDOfProjectType int
	global.DB.Model(&DictionaryType{}).
		Where("name = '项目类型'").Select("id").First(&dictionaryTypeIDOfProjectType)
	dictionaryItemsOfProjectType := []DictionaryItem{
		{Name: "EP", DictionaryTypeID: dictionaryTypeIDOfProjectType},
		{Name: "EPC", DictionaryTypeID: dictionaryTypeIDOfProjectType},
		{Name: "分销", DictionaryTypeID: dictionaryTypeIDOfProjectType},
	}
	dictionaryItems = append(dictionaryItems, dictionaryItemsOfProjectType...)

	var dictionaryTypeIDOfProjectStatus int
	global.DB.Model(&DictionaryType{}).
		Where("name = '项目状态'").Select("id").First(&dictionaryTypeIDOfProjectStatus)
	dictionaryItemsOfProjectStatus := []DictionaryItem{
		{Name: "EP", DictionaryTypeID: dictionaryTypeIDOfProjectStatus},
		{Name: "EPC", DictionaryTypeID: dictionaryTypeIDOfProjectStatus},
		{Name: "分销", DictionaryTypeID: dictionaryTypeIDOfProjectStatus},
	}
	dictionaryItems = append(dictionaryItems, dictionaryItemsOfProjectStatus...)

	var dictionaryTypeIDOfBankName int
	global.DB.Model(&DictionaryType{}).
		Where("name = '银行名称'").Select("id").First(&dictionaryTypeIDOfBankName)
	dictionaryItemsOfBankName := []DictionaryItem{
		{Name: "中国银行", DictionaryTypeID: dictionaryTypeIDOfBankName},
		{Name: "工商银行", DictionaryTypeID: dictionaryTypeIDOfBankName},
		{Name: "交通银行", DictionaryTypeID: dictionaryTypeIDOfBankName},
	}
	dictionaryItems = append(dictionaryItems, dictionaryItemsOfBankName...)

	var dictionaryTypeIDOfContractFundDirection int
	global.DB.Model(&DictionaryType{}).
		Where("name = '银行名称'").Select("id").First(&dictionaryTypeIDOfContractFundDirection)
	dictionaryItemsOfContractFundDirection := []DictionaryItem{
		{Name: "中国银行", DictionaryTypeID: dictionaryTypeIDOfContractFundDirection},
		{Name: "工商银行", DictionaryTypeID: dictionaryTypeIDOfContractFundDirection},
		{Name: "交通银行", DictionaryTypeID: dictionaryTypeIDOfContractFundDirection},
	}
	dictionaryItems = append(dictionaryItems, dictionaryItemsOfContractFundDirection...)

	var dictionaryTypeIDOfOurSignatory int
	global.DB.Model(&DictionaryType{}).
		Where("name = '我方签约主体'").Select("id").First(&dictionaryTypeIDOfOurSignatory)
	dictionaryItemsOfOurSignatory := []DictionaryItem{
		{Name: "北京公司", DictionaryTypeID: dictionaryTypeIDOfOurSignatory},
		{Name: "凯昌", DictionaryTypeID: dictionaryTypeIDOfOurSignatory},
		{Name: "凯祥", DictionaryTypeID: dictionaryTypeIDOfOurSignatory},
	}
	dictionaryItems = append(dictionaryItems, dictionaryItemsOfOurSignatory...)

	var dictionaryTypeIDOfSensitiveWord int
	global.DB.Model(&DictionaryType{}).
		Where("name = '敏感词'").Select("id").First(&dictionaryTypeIDOfSensitiveWord)
	dictionaryItemsOfSensitiveWord := []DictionaryItem{
		{Name: "伊朗", DictionaryTypeID: dictionaryTypeIDOfSensitiveWord},
		{Name: "委内瑞拉", DictionaryTypeID: dictionaryTypeIDOfSensitiveWord},
	}
	dictionaryItems = append(dictionaryItems, dictionaryItemsOfSensitiveWord...)

	for _, dictionaryItem := range dictionaryItems {
		err := global.DB.FirstOrCreate(&DictionaryItem{}, dictionaryItem).Error
		if err != nil {
			global.SugaredLogger.Errorln(err)
			return err
		}
	}

	return nil
}
