package dto

//以下为入参

type ListInput struct {
	PagingInput
	SortingInput
	SqlConditionInput
	AuthInput
	UserInfoInput //department service会用，以后改到服务里，不要放到公共区
}

type PagingInput struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type SortingInput struct {
	OrderBy string `json:"order_by"` //排序字段
	Desc    bool   `json:"desc"`     //是否为降序（从大到小）
}

type SqlConditionInput struct {
	SelectedColumns []string `json:"selected_columns"` //需要显示数据的列
}

type AuthInput struct {
	VerifyRole *bool `json:"verify_role"` //是否需要校验角色、分级显示
}

type UserInfoInput struct {
	RoleNames           []string //用户的角色名称数组
	BusinessDivisionIDs []int    //用户所属的事业部id数组
	DepartmentIDs       []int    //用户所属的部门id数组
}

//以下为出参

type ListOutput struct {
	PagingOutput
	SortingOutput
}

type PagingOutput struct {
	Page         int `json:"page"`
	PageSize     int `json:"page_size"`
	TotalPages   int `json:"total_pages"`
	TotalRecords int `json:"total_records"`
}

type SortingOutput struct {
	OrderBy string `json:"order_by"` //排序字段
	Desc    bool   `json:"desc"`     //是否为降序（从大到小）
}
