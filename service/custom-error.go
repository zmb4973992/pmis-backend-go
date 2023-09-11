package service

import (
	"pmis-backend-go/util"
)

var (
	ErrorFailToDeleteRecord               = util.GenerateCustomError(util.ErrorFailToDeleteRecord)
	ErrorFailToCreateRecord               = util.GenerateCustomError(util.ErrorFailToCreateRecord)
	ErrorFailToUpdateRecord               = util.GenerateCustomError(util.ErrorFailToUpdateRecord)
	ErrorFieldsToBeCreatedNotFound        = util.GenerateCustomError(util.ErrorFieldsToBeCreatedNotFound)
	ErrorFailToUpdateRBACGroupingPolicies = util.GenerateCustomError(util.ErrorFailToUpdateRBACGroupingPolicies)
)
