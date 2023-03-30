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
	ErrorRoleInfoNotFound
	ErrorRolePermissionDenied
	ErrorNeedAdminPrivilege
	ErrorFailToEncrypt
	ErrorInvalidRequest
	ErrorMethodNotAllowed
	ErrorInvalidColumns
	ErrorRequestFrequencyTooHigh
	ErrorFieldsToBeCreatedNotFound
	ErrorFieldsToBeUpdatedNotFound
	ErrorSortingFieldDoesNotExist
	ErrorUserIDDoesNotExist
	ErrorFailToGenerateCaptcha
	ErrorWrongCaptcha
	ErrorFailToGenerateToken
	ErrorFailToDeleteFiles
	ErrorDictionaryTypeNameNotFound
	ErrorInvalidDateFormat
	ErrorFileNotFound
	ErrorDuplicateRecord
	ErrorFailToCalculateSelfProgress
	ErrorFailToCalculateSuperiorProgress
	ErrorFailToCalculateSelfAndSuperiorProgress
	ErrorWrongSuperiorInformation
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
	ErrorRoleInfoNotFound:          "未找到用户的角色信息",

	ErrorAccessTokenInvalid:   "access_token无效",
	ErrorAccessTokenNotFound:  "缺少access_token",
	ErrorRolePermissionDenied: "用户所属角色的权限不足",
	ErrorNeedAdminPrivilege:   "权限不足，该操作需要管理员权限",
	ErrorUserIDDoesNotExist:   "用户ID不存在",

	ErrorFailToEncrypt:    "加密失败",
	ErrorInvalidRequest:   "请求路径错误",
	ErrorMethodNotAllowed: "请求方法错误",
	ErrorInvalidColumns:   "列名无效",

	ErrorRequestFrequencyTooHigh:                "请求频率过高，请稍后再试",
	ErrorFieldsToBeCreatedNotFound:              "未找到需要新增的字段",
	ErrorFieldsToBeUpdatedNotFound:              "未找到需要更新的字段",
	ErrorSortingFieldDoesNotExist:               "排序字段不存在",
	ErrorFailToGenerateCaptcha:                  "生成验证码失败",
	ErrorWrongCaptcha:                           "验证码错误",
	ErrorFailToGenerateToken:                    "生成token失败",
	ErrorFailToDeleteFiles:                      "删除文件失败",
	ErrorDictionaryTypeNameNotFound:             "字典名称未找到",
	ErrorInvalidDateFormat:                      "日期格式无效",
	ErrorFileNotFound:                           "文件未找到",
	ErrorDuplicateRecord:                        "系统已存在该日期、该类型的记录。请修改后再提交",
	ErrorFailToCalculateSelfProgress:            "计算自身进度失败，错误详情请查看系统日志文件",
	ErrorFailToCalculateSuperiorProgress:        "计算上级进度失败，错误详情请查看系统日志文件",
	ErrorFailToCalculateSelfAndSuperiorProgress: "计算自身和上级进度失败，错误详情请查看系统日志文件",
	ErrorWrongSuperiorInformation:               "上级信息错误，可能缺失项目id或层级",
}

func GetMessage(code int) string {
	message, ok := Message[code]
	if !ok {
		return "由于错误代码未定义返回信息，导致获取错误信息失败，" +
			"建议检查utils/code-and-message相关设置"
	}
	return message
}
