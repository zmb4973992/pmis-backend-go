package model

type Temp struct {
	BasicModel
	BatchID int64 //批次id，批量添加数据时确定批次数据的唯一性
	//连接其他表的id
	OrganizationID         *int64
	RelatedPartyID         *int64
	ProjectID              *int64
	ContractID             *int64
	IncomeAndExpenditureID *int64
	//连接dictionary_item表的id

}

// TableName 将表名改为project
func (t *Temp) TableName() string {
	return "temp"
}
