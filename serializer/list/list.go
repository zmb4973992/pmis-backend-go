package list

//以下为入参

// Input 标准的list入参，几乎所有的list都要用到这些参数
type Input struct {
	PagingInput
	SortingInput
}

type PagingInput struct {
	Page     int  `json:"page" binding:"omitempty,gte=1"`
	PageSize *int `json:"page_size"`
}

type SortingInput struct {
	OrderBy string `json:"order_by"` //排序字段
	Desc    bool   `json:"desc"`     //是否为降序（从大到小）
}

// DataScopeInput 用于校验角色、分级显示的入参，按需导入
type DataScopeInput struct {
	UserSnowID int64 `json:"-"`
	//LoadDataScopeByRole bool `json:"load_data_scope_by_role,omitempty"` //根据组织确定数据范围
}

//以下为出参

// Output 为标准的出参形式，几乎所有的list都按这个标准来
type Output struct {
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
