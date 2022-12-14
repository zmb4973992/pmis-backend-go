package dto

// OperationRecordGetDTO
// mapstructure用于list
// list方法中返回的数据为[]map，需要借助mapstructure转换为struct，再返回给前端
type OperationRecordGetDTO struct {
	BaseDTO    `mapstructure:",squash"` //这里是嵌套结构体，mapstructure必须加squash，否则无法匹配
	ProjectID  *int                     `json:"project_id" mapstructure:"project_id"`   //项目id
	OperatorID *int                     `json:"operator_id" mapstructure:"operator_id"` //操作人id
	Date       *string                  `json:"date" mapstructure:"-"`                  //日期
	Action     *string                  `json:"action" mapstructure:"action"`           //动作
	Detail     *string                  `json:"detail" mapstructure:"detail"`           //详情
}

type OperationRecordCreateOrUpdateDTO struct {
	BaseDTO
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
	ID         int     `form:"id"`
	ProjectID  *int    `form:"project_id"`
	OperatorID *int    `form:"operator_id"`
	DateGte    *string `form:"date_gte"`
	DateLte    *string `form:"date_lte"`
	Action     *string `form:"action"`
}
