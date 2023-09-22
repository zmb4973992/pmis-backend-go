package model

import "time"

type Contract struct {
	BasicModel
	//连接其他表的id
	ProjectId      *int64 //项目id
	OrganizationId *int64 //组织id
	RelatedPartyId *int64 //相关方id
	//连接dictionary_item表的id
	FundDirection *int64 //资金方向
	OurSignatory  *int64 //我方签约主体
	Currency      *int64 //币种
	Type          *int64 //类型(总包、采购、结算单等)
	//日期
	ApprovalDate      *time.Time `gorm:"type:date"` //过审日期、批准日期
	SigningDate       *time.Time `gorm:"type:date"` //签约日期
	EffectiveDate     *time.Time `gorm:"type:date"` //生效日期
	CommissioningDate *time.Time `gorm:"type:date"` //调试日期
	CompletionDate    *time.Time `gorm:"type:date"` //完工日期
	//数字(允许为0、nil)
	Amount             *float64 //金额
	ExchangeRate       *float64 //汇率
	ConstructionPeriod *int     //工期，天
	//数字(不允许为0、nil，必须有值)，暂无

	//字符串(允许为nil)
	Name        *string //合同名称
	Code        *string //合同编码
	Content     *string //工作内容
	Deliverable *string //交付物
	PenaltyRule *string //罚则
	Operator    *string //经办人
	FileIds     *string //附件

	//字符串(不允许为nil，必须有值)，暂无
}

// TableName 修改表名
func (c *Contract) TableName() string {
	return "contract"
}
