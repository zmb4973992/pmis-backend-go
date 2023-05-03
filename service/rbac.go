package service

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/util"
	"strconv"
)

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type rbacUpdatePolicyByRoleID struct {
	RoleID int
}

type rbacUpdatePolicyByMenuID struct {
	MenuID int
}
type rbacUpdatePolicyByApiID struct {
	ApiID int
}

type rbacUpdateGroupingPolicyByGroup struct {
	Group   string
	Members []string
}

type rbacUpdateGroupingPolicyByMember struct {
	Member string
	Groups []string
}

func (r *rbacUpdatePolicyByRoleID) Update() error {
	if r.RoleID == 0 {
		return nil
	}

	cachedEnforcer, err := util.NewCachedEnforcer()
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return err
	}

	subject := strconv.Itoa(r.RoleID)
	_, err = cachedEnforcer.RemoveFilteredPolicy(0, subject)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return err
	}

	//找到角色拥有的菜单
	var menuIDs []int
	global.DB.Model(&model.RoleAndMenu{}).Where("role_id = ?", r.RoleID).
		Select("menu_id").Find(&menuIDs)

	//找到菜单拥有的api
	var apiIDs []int
	global.DB.Model(&model.MenuAndApi{}).Where("menu_id in ?", menuIDs).
		Select("api_id").Find(&apiIDs)

	//找到api详细信息
	var rbacRules [][]string
	for _, apiID := range apiIDs {
		var api model.Api
		err = global.DB.Where("id = ?", apiID).First(&api).Error
		if err != nil {
			continue
		}
		//如果api带param参数，那么rbac规则要带上正则，否则无法放行
		if api.WithParam {
			api.Path += "/*"
		}
		rbacRules = append(rbacRules, []string{subject, api.Path, api.Method})
	}

	_, err = cachedEnforcer.AddPolicies(rbacRules)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return err
	}

	//修改了policy以后，因为用的是cachedEnforcer，所以要清除缓存
	err = cachedEnforcer.InvalidateCache()
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return err
	}
	return nil
}

func (r *rbacUpdatePolicyByMenuID) Update() error {
	if r.MenuID == 0 {
		return nil
	}

	//先找到菜单关联的角色id
	var roleIDs []int
	global.DB.Model(&model.RoleAndMenu{}).Where("menu_id = ?", r.MenuID).
		Select("role_id").Find(&roleIDs)

	for _, roleID := range roleIDs {
		param := rbacUpdatePolicyByRoleID{RoleID: roleID}
		err := param.Update()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *rbacUpdatePolicyByApiID) Update() error {
	if r.ApiID == 0 {
		return nil
	}

	//先找到api关联的菜单id
	var menuIDs []int
	global.DB.Model(&model.MenuAndApi{}).Where("api_id in ?", r.ApiID).
		Select("menu_id").Find(&menuIDs)

	for _, menuID := range menuIDs {
		param := rbacUpdatePolicyByMenuID{MenuID: menuID}
		err := param.Update()
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *rbacUpdateGroupingPolicyByGroup) Update() error {
	if len(u.Members) == 0 {
		return nil
	}

	cachedEnforcer, err := util.NewCachedEnforcer()
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return err
	}

	_, err = cachedEnforcer.RemoveFilteredGroupingPolicy(1, u.Group)
	if err != nil {
		return err
	}

	for _, member := range u.Members {
		_, err = cachedEnforcer.AddGroupingPolicy([]string{member, u.Group})
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *rbacUpdateGroupingPolicyByMember) Update() error {
	if len(u.Groups) == 0 {
		return nil
	}

	cachedEnforcer, err := util.NewCachedEnforcer()
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return err
	}

	_, err = cachedEnforcer.RemoveFilteredGroupingPolicy(0, u.Member)
	if err != nil {
		return err
	}

	for _, group := range u.Groups {
		_, err = cachedEnforcer.AddGroupingPolicy([]string{u.Member, group})
		if err != nil {
			return err
		}
	}

	return nil
}
