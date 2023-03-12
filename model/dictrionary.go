package model

type Dictionary struct {
	BasicModel
	Country                         *string //国家
	Province                        *string //省份或地区
	ReceiptOrPaymentTerm            *string //收付款方式
	Currency                        *string //币种
	OrderForCurrency                *int    //币种排序
	ContractType                    *string //合同类型
	ProjectType                     *string //项目类型
	ProjectStatus                   *string //项目状态
	OrderForStatusOfProject         *int    //项目状态排序
	NameOfBank                      *string //银行名称
	FundDirectionOfContract         *string //合同资金方向
	OrderForFundDirectionOfContract *int    //合同资金方向的排序
	OurSignatory                    *string //我方签约主体
	SensitiveWord                   *string //敏感词
}

// TableName 修改数据库的表名
func (*Dictionary) TableName() string {
	return "dictionary"
}
