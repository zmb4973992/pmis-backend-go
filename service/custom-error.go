package service

import (
	"errors"
	"pmis-backend-go/util"
)

var (
	ErrorFailToDeleteRecord               = errors.New(util.GetMessage(util.ErrorFailToDeleteRecord))
	ErrorFailToCreateRecord               = errors.New(util.GetMessage(util.ErrorFailToCreateRecord))
	ErrorFailToUpdateRecord               = errors.New(util.GetMessage(util.ErrorFailToUpdateRecord))
	ErrorFieldsToBeCreatedNotFound        = errors.New(util.GetMessage(util.ErrorFieldsToBeCreatedNotFound))
	ErrorFailToUpdateRBACGroupingPolicies = errors.New(util.GetMessage(util.ErrorFailToUpdateRBACGroupingPolicies))
)
