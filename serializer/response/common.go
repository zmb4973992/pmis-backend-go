package response

import (
	"pmis-backend-go/serializer/list"
	"pmis-backend-go/util"
)

type Common struct {
	Data    any    `json:"data"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// List 这里不直接嵌套response，是为了调整字段显示顺序，
// 另外嵌套多层也会导致出参结果需要嵌套，比较麻烦
type List struct {
	Data    any                `json:"data"`
	Paging  *list.PagingOutput `json:"paging"`
	Code    int                `json:"code"`
	Message string             `json:"message"`
}

func GenerateCommon(data any, errCode int) Common {
	return Common{
		Data:    data,
		Code:    errCode,
		Message: util.GetErrorDescription(errCode),
	}
}

func GenerateList(dataList any, errCode int, paging *list.PagingOutput) List {
	return List{
		Data:    dataList,
		Paging:  paging,
		Code:    errCode,
		Message: util.GetErrorDescription(errCode),
	}
}
