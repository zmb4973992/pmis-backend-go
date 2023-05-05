package service

import (
	"errors"
	"gorm.io/gorm"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
	"strconv"
)

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

var (
	ErrorFailToDeleteRecord               = errors.New(util.GetMessage(util.ErrorFailToDeleteRecord))
	ErrorFailToCreateRecord               = errors.New(util.GetMessage(util.ErrorFailToCreateRecord))
	ErrorFailToUpdateRecord               = errors.New(util.GetMessage(util.ErrorFailToUpdateRecord))
	ErrorFieldsToBeCreatedNotFound        = errors.New(util.GetMessage(util.ErrorFieldsToBeCreatedNotFound))
	ErrorFailToUpdateRBACGroupingPolicies = errors.New(util.GetMessage(util.ErrorFailToUpdateRBACGroupingPolicies))
)

type RoleAndUserUpdateByRoleSnowID struct {
	Creator      int64
	LastModifier int64

	RoleSnowID  int64   `json:"-"`
	UserSnowIDs []int64 `json:"user_snow_ids,omitempty"`
}

type RoleAndUserUpdateByUserSnowID struct {
	Creator      int64
	LastModifier int64

	UserSnowID  int64   `json:"-"`
	RoleSnowIDs []int64 `json:"role_snow_ids,omitempty"`
}

func (r *RoleAndUserUpdateByRoleSnowID) Update() response.Common {
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		//先删掉原始记录
		err := tx.Where("role_snow_id = ?", r.RoleSnowID).Delete(&model.UserAndRole{}).Error
		if err != nil {
			global.SugaredLogger.Errorln(err)
			return ErrorFailToDeleteRecord
		}

		//再增加新的记录
		var paramOut []model.UserAndRole
		for _, userSnowID := range r.UserSnowIDs {
			var record model.UserAndRole
			if r.Creator > 0 {
				record.Creator = &r.Creator
			}
			if r.LastModifier > 0 {
				record.LastModifier = &r.LastModifier
			}

			record.RoleSnowID = r.RoleSnowID
			record.UserSnowID = userSnowID
			paramOut = append(paramOut, record)
		}

		for i := range paramOut {
			//计算有修改值的字段数，分别进行不同处理
			tempParamOut, err := util.StructToMap(paramOut[i])
			if err != nil {
				return ErrorFailToUpdateRecord
			}
			paramOutForCounting := util.MapCopy(tempParamOut,
				"Creator", "LastModifier", "CreateAt", "UpdatedAt")

			if len(paramOutForCounting) == 0 {
				return ErrorFieldsToBeCreatedNotFound
			}
		}

		err = global.DB.Create(&paramOut).Error
		if err != nil {
			global.SugaredLogger.Errorln(err)
			return ErrorFailToCreateRecord
		}

		//更新casbin的rbac分组规则
		var param1 rbacUpdateGroupingPolicyByGroup
		param1.Group = strconv.FormatInt(r.RoleSnowID, 10)
		for _, userSnowID := range r.UserSnowIDs {
			param1.Members = append(param1.Members, strconv.FormatInt(userSnowID, 10))
		}
		err = param1.Update()
		if err != nil {
			return ErrorFailToUpdateRBACGroupingPolicies
		}

		return nil
	})

	switch err {
	case nil:
		return response.Success()
	case ErrorFailToCreateRecord:
		return response.Failure(util.ErrorFailToCreateRecord)
	case ErrorFailToDeleteRecord:
		return response.Failure(util.ErrorFailToDeleteRecord)
	case ErrorFieldsToBeCreatedNotFound:
		return response.Failure(util.ErrorFieldsToBeCreatedNotFound)
	case ErrorFailToUpdateRBACGroupingPolicies:
		return response.Failure(util.ErrorFailToUpdateRBACGroupingPolicies)
	default:
		return response.Failure(util.ErrorFailToUpdateRecord)
	}
}

func (r *RoleAndUserUpdateByUserSnowID) Update() response.Common {
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		//先删掉原始记录
		err := tx.Where("user_snow_id = ?", r.UserSnowID).Delete(&model.UserAndRole{}).Error
		if err != nil {
			global.SugaredLogger.Errorln(err)
			return ErrorFailToDeleteRecord
		}

		//再增加新的记录
		var paramOut []model.UserAndRole
		for _, roleID := range r.RoleSnowIDs {
			var record model.UserAndRole
			if r.Creator > 0 {
				record.Creator = &r.Creator
			}
			if r.LastModifier > 0 {
				record.LastModifier = &r.LastModifier
			}

			record.UserSnowID = r.UserSnowID
			record.RoleSnowID = roleID
			paramOut = append(paramOut, record)
		}

		for i := range paramOut {
			//计算有修改值的字段数，分别进行不同处理
			tempParamOut, err := util.StructToMap(paramOut[i])
			if err != nil {
				return ErrorFailToUpdateRecord
			}
			paramOutForCounting := util.MapCopy(tempParamOut,
				"Creator", "LastModifier", "CreateAt", "UpdatedAt")

			if len(paramOutForCounting) == 0 {
				return ErrorFieldsToBeCreatedNotFound
			}
		}

		err = global.DB.Create(&paramOut).Error
		if err != nil {
			global.SugaredLogger.Errorln(err)
			return ErrorFailToCreateRecord
		}

		//更新casbin的rbac分组规则
		var param1 rbacUpdateGroupingPolicyByMember
		param1.Member = strconv.FormatInt(r.UserSnowID, 10)
		for _, roleSnowID := range r.RoleSnowIDs {
			param1.Groups = append(param1.Groups, strconv.FormatInt(roleSnowID, 10))
		}
		err = param1.Update()
		if err != nil {
			return ErrorFailToUpdateRBACGroupingPolicies
		}

		return nil
	})

	switch err {
	case nil:
		return response.Success()
	case ErrorFailToCreateRecord:
		return response.Failure(util.ErrorFailToCreateRecord)
	case ErrorFailToDeleteRecord:
		return response.Failure(util.ErrorFailToDeleteRecord)
	case ErrorFieldsToBeCreatedNotFound:
		return response.Failure(util.ErrorFieldsToBeCreatedNotFound)
	case ErrorFailToUpdateRBACGroupingPolicies:
		return response.Failure(util.ErrorFailToUpdateRBACGroupingPolicies)
	default:
		return response.Failure(util.ErrorFailToUpdateRecord)
	}
}
