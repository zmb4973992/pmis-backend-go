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

type RoleAndUserUpdateByRoleID struct {
	Creator      int
	LastModifier int

	RoleID  int   `json:"-"`
	UserIDs []int `json:"user_ids,omitempty"`
}

type RoleAndUserUpdateByUserID struct {
	Creator      int
	LastModifier int

	UserID  int   `json:"-"`
	RoleIDs []int `json:"role_ids,omitempty"`
}

func (r *RoleAndUserUpdateByRoleID) Update() response.Common {
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		//先删掉原始记录
		err := tx.Where("role_id = ?", r.RoleID).Delete(&model.UserAndRole{}).Error
		if err != nil {
			global.SugaredLogger.Errorln(err)
			return ErrorFailToDeleteRecord
		}

		//再增加新的记录
		var paramOut []model.UserAndRole
		for _, userID := range r.UserIDs {
			var record model.UserAndRole
			if r.Creator > 0 {
				record.Creator = &r.Creator
			}
			if r.LastModifier > 0 {
				record.LastModifier = &r.LastModifier
			}

			record.RoleSnowID = r.RoleID
			record.UserSnowID = userID
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
		param1.Group = strconv.Itoa(r.RoleID)
		for _, userID := range r.UserIDs {
			param1.Members = append(param1.Members, strconv.Itoa(userID))
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

func (r *RoleAndUserUpdateByUserID) Update() response.Common {
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		//先删掉原始记录
		err := tx.Where("user_id = ?", r.UserID).Delete(&model.UserAndRole{}).Error
		if err != nil {
			global.SugaredLogger.Errorln(err)
			return ErrorFailToDeleteRecord
		}

		//再增加新的记录
		var paramOut []model.UserAndRole
		for _, roleID := range r.RoleIDs {
			var record model.UserAndRole
			if r.Creator > 0 {
				record.Creator = &r.Creator
			}
			if r.LastModifier > 0 {
				record.LastModifier = &r.LastModifier
			}

			record.UserSnowID = r.UserID
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
		param1.Member = strconv.Itoa(r.UserID)
		for _, roleID := range r.RoleIDs {
			param1.Groups = append(param1.Groups, strconv.Itoa(roleID))
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
