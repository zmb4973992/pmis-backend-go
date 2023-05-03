package service

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
	"strconv"
)

type roleAndUser struct{}

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

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

func (r1 *RoleAndUserUpdateByRoleID) Update() response.Common {
	//先删掉原始记录
	err := global.DB.Where("role_id = ?", r1.RoleID).Delete(&model.UserAndRole{}).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToDeleteRecord)
	}

	//再增加新的记录
	var paramOut []model.UserAndRole
	for _, userID := range r1.UserIDs {
		var record model.UserAndRole
		if r1.Creator > 0 {
			record.Creator = &r1.Creator
		}
		if r1.LastModifier > 0 {
			record.LastModifier = &r1.LastModifier
		}

		record.RoleID = r1.RoleID
		record.UserID = userID
		paramOut = append(paramOut, record)
	}

	for i := range paramOut {
		//计算有修改值的字段数，分别进行不同处理
		tempParamOut, err := util.StructToMap(paramOut[i])
		if err != nil {
			return response.Failure(util.ErrorFailToUpdateRecord)
		}
		paramOutForCounting := util.MapCopy(tempParamOut,
			"Creator", "LastModifier", "CreateAt", "UpdatedAt")

		if len(paramOutForCounting) == 0 {
			return response.Failure(util.ErrorFieldsToBeCreatedNotFound)
		}
	}

	err = global.DB.Create(&paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToCreateRecord)
	}

	//更新casbin的rbac分组规则
	var param1 rbacUpdateGroupingPolicyByFather
	param1.Father = strconv.Itoa(r1.RoleID)
	for _, userID := range r1.UserIDs {
		param1.Sons = append(param1.Sons, strconv.Itoa(userID))
	}
	err = param1.UpdateGroupingPolicyByFather()
	if err != nil {
		return response.Failure(util.ErrorFailToUpdateRBACGroupingPolicies)
	}

	return response.Success()
}

func (r2 *RoleAndUserUpdateByUserID) Update() response.Common {
	//先删掉原始记录
	err := global.DB.Where("user_id = ?", r2.UserID).Delete(&model.UserAndRole{}).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToUpdateRecord)
	}

	//再增加新的记录
	var paramOut []model.UserAndRole
	for _, roleID := range r2.RoleIDs {
		var record model.UserAndRole
		if r2.Creator > 0 {
			record.Creator = &r2.Creator
		}
		if r2.LastModifier > 0 {
			record.LastModifier = &r2.LastModifier
		}

		record.UserID = r2.UserID
		record.RoleID = roleID
		paramOut = append(paramOut, record)
	}

	for i := range paramOut {
		//计算有修改值的字段数，分别进行不同处理
		tempParamOut, err := util.StructToMap(paramOut[i])
		if err != nil {
			return response.Failure(util.ErrorFailToUpdateRecord)
		}
		paramOutForCounting := util.MapCopy(tempParamOut,
			"Creator", "LastModifier", "CreateAt", "UpdatedAt")

		if len(paramOutForCounting) == 0 {
			return response.Failure(util.ErrorFieldsToBeCreatedNotFound)
		}
	}

	err = global.DB.Create(&paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToUpdateRecord)
	}

	//更新casbin的rbac分组规则
	var param1 rbacUpdateGroupingPolicyBySon
	param1.Son = strconv.Itoa(r2.UserID)
	for _, roleID := range r2.RoleIDs {
		param1.Fathers = append(param1.Fathers, strconv.Itoa(roleID))
	}
	err = param1.UpdateGroupingPolicyBySon()
	if err != nil {
		return response.Failure(util.ErrorFailToUpdateRBACGroupingPolicies)
	}
	return response.Success()
}
