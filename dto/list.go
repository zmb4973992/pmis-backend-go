package dto

//以下为入参

// ListInput 标准的list入参，几乎所有的list都要用到这些参数
type ListInput struct {
	PagingInput
	SortingInput
	SqlConditionInput
}

type PagingInput struct {
	Page     int `json:"page" binding:"omitempty,gte=1"`
	PageSize int `json:"page_size" binding:"omitempty,gte=1"`
}

type SortingInput struct {
	OrderBy string `json:"order_by"` //排序字段
	Desc    bool   `json:"desc"`     //是否为降序（从大到小）
}

type SqlConditionInput struct {
	SelectedColumns []string `json:"selected_columns"` //需要显示数据的列
}

// AuthInput 用于校验角色、分级显示的入参，按需导入
type AuthInput struct {
	IsShowedByRole bool `json:"is_showed_by_role,omitempty"` //根据角色来分级显示
	UserID         int
}

type AuthInputOld struct {
	UserInfoInput
	VerifyRole *bool `json:"verify_role"` //是否需要校验角色、分级显示
}

type UserInfoInput struct {
	RoleNames           []string //用户的角色名称数组
	BusinessDivisionIDs []int    //用户所属的事业部id数组
	DepartmentIDs       []int    //用户所属的部门id数组
}

//以下为出参

// ListOutput 为标准的出参形式，几乎所有的list都按这个标准来
type ListOutput struct {
	PagingOutput
	SortingOutput
}

type PagingOutput struct {
	Page            int `json:"page"`
	PageSize        int `json:"page_size"`
	NumberOfPages   int `json:"number_of_pages"`
	NumberOfRecords int `json:"number_of_records"`
}

type SortingOutput struct {
	OrderBy string `json:"order_by"` //排序字段
	Desc    bool   `json:"desc"`     //是否为降序（从大到小）
}
