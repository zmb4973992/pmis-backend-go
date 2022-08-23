package dto

type OperationRecordCreateOrUpdateDTO struct {
	BaseDTO
	ID         int     `json:"id"`
	ProjectID  *int    `json:"project_id" binding:"required"`  //项目id
	OperatorID *int    `json:"operator_id" binding:"required"` //操作人id
	Date       *string `json:"date" binding:"required"`        //日期
	Action     *string `json:"action" binding:"required"`      //动作
	Detail     *string `json:"detail" binding:"required"`      //详情
}

// OperationRecordListDTO 是list查询的过滤器
// 在dto传递给sqlCondition时，空值会被忽略
type OperationRecordListDTO struct {
	ListDTO
	ID         int     `json:"id"`
	ProjectID  *int    `json:"project_id"`
	OperatorID *int    `json:"operator_id"`
	DateGte    *string `json:"date_gte"`
	DateLte    *string `json:"date_lte"`
	Action     *string `json:"action"`
}

// OperationRecordGetDTO
// mapstructure用于list
// list方法中返回的数据为[]map，需要借助mapstructure转换为struct，再返回给前端
type OperationRecordGetDTO struct {
	ProjectID  *int    `json:"project_id" mapstructure:"project_id"`   //项目id
	OperatorID *int    `json:"operator_id" mapstructure:"operator_id"` //操作人id
	Date       *string `json:"date" mapstructure:"-"`                  //日期
	Action     *string `json:"action" mapstructure:"action"`           //动作
	Detail     *string `json:"detail" mapstructure:"detail"`           //详情
}
