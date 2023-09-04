package service

import (
	"errors"
	"pmis-backend-go/util"
)

var (
	ErrorFailToDeleteRecord               = errors.New(util.GetErrorDescription(util.ErrorFailToDeleteRecord))
	ErrorFailToCreateRecord               = errors.New(util.GetErrorDescription(util.ErrorFailToCreateRecord))
	ErrorFailToUpdateRecord               = errors.New(util.GetErrorDescription(util.ErrorFailToUpdateRecord))
	ErrorFieldsToBeCreatedNotFound        = errors.New(util.GetErrorDescription(util.ErrorFieldsToBeCreatedNotFound))
	ErrorFailToUpdateRBACGroupingPolicies = errors.New(util.GetErrorDescription(util.ErrorFailToUpdateRBACGroupingPolicies))
)
