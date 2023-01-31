package util

// 自定义错误的code
const (
	Success = iota
	Error
	ErrorRecordNotFound
	ErrorNotEnoughParameters
	ErrorInvalidURIParameters
	ErrorInvalidFormDataParameters
	ErrorInvalidJSONParameters
	ErrorInvalidQueryParameters
	ErrorFailToDeleteRecord
	ErrorFileTooLarge
	ErrorFailToUploadFiles
	ErrorFailToTransferDataFromDtoToModel
	ErrorInvalidBaseController
	ErrorFailToCreateRecord
	ErrorFailToUpdateRecord
	ErrorInvalidUsernameOrPassword
	ErrorUsernameExist
	ErrorPasswordIncorrect
	ErrorAccessTokenInvalid
	ErrorAccessTokenNotFound
	ErrorPermissionDenied
	ErrorNeedAdminPrivilege
	ErrorFailToEncrypt
	ErrorInvalidRequest
	ErrorMethodNotAllowed
	ErrorInvalidColumns
	ErrorRequestFrequencyTooHigh
	ErrorFieldsToBeUpdatedNotFound
	ErrorSortingFieldDoesNotExist
)

// Message 自定义错误的message
var Message = map[int]string{
	Success: "成功",
	Error:   "错误",

	ErrorRecordNotFound:                   "未找到指定记录",
	ErrorNotEnoughParameters:              "没有足够的参数",
	ErrorInvalidURIParameters:             "URI参数无效",
	ErrorInvalidFormDataParameters:        "form-data参数无效",
	ErrorInvalidJSONParameters:            "JSON参数无效",
	ErrorInvalidQueryParameters:           "query参数无效",
	ErrorFailToDeleteRecord:               "删除记录失败",
	ErrorFileTooLarge:                     "文件过大",
	ErrorFailToUploadFiles:                "上传文件失败",
	ErrorFailToTransferDataFromDtoToModel: "dto转model失败，请检查service层",
	ErrorInvalidBaseController:            "BaseController配置错误，请检查",
	ErrorFailToCreateRecord:               "添加记录失败",
	ErrorFailToUpdateRecord:               "修改记录失败",

	ErrorInvalidUsernameOrPassword: "用户名或密码错误",
	ErrorUsernameExist:             "用户名已存在",
	ErrorPasswordIncorrect:         "密码错误",

	ErrorAccessTokenInvalid:  "access_token无效",
	ErrorAccessTokenNotFound: "缺少access_token",
	ErrorPermissionDenied:    "权限不足",
	ErrorNeedAdminPrivilege:  "权限不足，该操作需要管理员权限",

	ErrorFailToEncrypt:    "加密失败",
	ErrorInvalidRequest:   "请求路径错误",
	ErrorMethodNotAllowed: "请求方法错误",
	ErrorInvalidColumns:   "列名无效",

	ErrorRequestFrequencyTooHigh:   "请求频率过高，请稍后再试",
	ErrorFieldsToBeUpdatedNotFound: "未找到需要更新的字段",
	ErrorSortingFieldDoesNotExist:  "排序字段不存在",
}

func GetMessage(code int) string {
	message, ok := Message[code]
	if !ok {
		return "由于错误代码未定义返回信息，导致获取错误信息失败，" +
			"建议检查utils/code-and-message相关设置"
	}
	return message
}