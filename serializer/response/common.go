package response

import (
	"pmis-backend-go/dto"
	"pmis-backend-go/util"
)

type Common struct {
	Data    any    `json:"data"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// List 这里不直接嵌套response，是为了调整字段显示顺序，另外嵌套多层也会导致出参结果需要嵌套，略麻烦
type List struct {
	Data    any            `json:"data"`
	Paging  *dto.PagingDTO `json:"paging"`
	Code    int            `json:"code"`
	Message string         `json:"message"`
}

func Success() Common {
	return Common{
		Data:    nil,
		Code:    util.Success,
		Message: util.GetMessage(util.Success),
	}
}

func SuccessWithData(data any) Common {
	return Common{
		Data:    data,
		Code:    util.Success,
		Message: util.GetMessage(util.Success),
	}
}

func Failure(errCode int) Common {
	return Common{
		Data:    nil,
		Code:    errCode,
		Message: util.GetMessage(errCode),
	}
}

func FailureForList(errCode int) List {
	return List{
		Data:    nil,
		Paging:  nil,
		Code:    errCode,
		Message: util.GetMessage(errCode),
	}
}
