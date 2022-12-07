package dto

type ListDTO struct {
	PagingDTO
	SortingDTO
	SqlDTO
	AuthDTO
	UserInfoDTO
}

type PagingDTO struct {
	Page         int `json:"page" form:"page"`
	PageSize     int `json:"page_size" form:"page_size"`
	TotalPages   int `json:"total_pages"`
	TotalRecords int `json:"total_records"`
}

type SortingDTO struct {
	OrderBy string `json:"order_by" form:"order_by"` //排序字段
	Desc    bool   `json:"desc" form:"desc"`         //是否为降序（从大到小）
}

type SqlDTO struct {
	SelectedColumns []string `json:"selected_columns"` //需要显示数据的列
}

type AuthDTO struct {
	VerifyRole *bool `json:"verify_role"` //是否需要校验角色、分级显示
}

type UserInfoDTO struct {
	RoleNames           []string //用户的角色名称数组
	BusinessDivisionIDs []int    //用户所属的事业部id数组
	DepartmentIDs       []int    //用户所属的部门id数组
}
