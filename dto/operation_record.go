package dto

type OperationRecordGetDTO struct {
	ProjectID  *int    `json:"project_id"`  //项目id
	OperatorID *int    `json:"operator_id"` //操作人id
	Date       *string `json:"date"`        //日期
	Action     *string `json:"action"`      //动作
	Detail     *string `json:"detail"`      //详情
}

func (OperationRecordGetDTO) GetDTO() {}

type Test interface {
	GetDTO()
}

type OperationRecordCreateAndUpdateDTO struct {
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
	ID         int     `json:"id"`
	ProjectID  *int    `json:"project_id"`
	OperatorID *int    `json:"operator_id"`
	DateGte    *string `json:"date_gte"`
	DateLte    *string `json:"date_lte"`
	Action     *string `json:"action"`

	ListDTO
}
