package model

import (
	"gorm.io/gorm"
	"time"
)

type Progress struct {
	BasicModel
	DisassemblyID *int64     //拆解情况ID
	Date          *time.Time `gorm:"type:date"`
	Type          *int64
	Value         *float64
	DataSource    *int64
	Remarks       *string

	//PlannedProgress             *float64 //初始计划进度
	//RemarkOfPlannedProgress     *string  //初始计划进度的备注
	//DataSourceOfPlannedProgress *string  //初始计划进度的数据来源

	//ActualProgress                    *float64 //实际进度
	//RemarkOfActualProgress            *string  //实际进度的备注
	//DataSourceOfActualProgress        *string  //实际进度的数据来源
	//ActualProgressIsTemporarilyFilled *bool    //实际进度是否为临时填充

	//PredictedProgress             *float64 //预测进度
	//RemarkOfPredictedProgress     *string  //预测进度的备注
	//DataSourceOfPredictedProgress *string  //预测进度的数据来源
}

// TableName 修改表名
func (*Progress) TableName() string {
	return "progress"
}

func (d *Progress) BeforeDelete(tx *gorm.DB) error {
	if d.ID == 0 {
		return nil
	}

	//删除相关的子表记录
	//先find、再delete，可以激活相关的钩子函数

	//fmt.Println("进度id：", d.ID)

	return nil
}
