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

type rbacUpdatePolicyByRoleId struct {
	RoleId int64
}

type rbacUpdatePolicyByMenuId struct {
	MenuId int64
}
type rbacUpdatePolicyByApiId struct {
	ApiId int64
}

type rbacUpdateGroupingPolicyByGroup struct {
	Group   string
	Members []string
}

type rbacUpdateGroupingPolicyByMember struct {
	Member string
	Groups []string
}

func (r *rbacUpdatePolicyByRoleId) Update() error {
	if r.RoleId == 0 {
		return nil
	}

	cachedEnforcer, err := util.NewCachedEnforcer()
	if err != nil {
		return err
	}

	subject := strconv.FormatInt(r.RoleId, 10)

	_, err = cachedEnforcer.RemoveFilteredPolicy(0, subject)
	if err != nil {
		return err
	}

	//找到角色拥有的菜单
	var menuIds []int64
	global.DB.Model(&model.RoleAndMenu{}).Where("role_id = ?", r.RoleId).
		Select("menu_id").Find(&menuIds)

	//找到菜单拥有的api
	var apiIds []int64
	global.DB.Model(&model.MenuAndApi{}).Where("menu_id in ?", menuIds).
		Select("api_id").Find(&apiIds)

	//找到api详细信息
	var rbacRules [][]string
	for _, apiId := range apiIds {
		var api model.Api
		err = global.DB.Where("id = ?", apiId).First(&api).Error
		if err != nil {
			continue
		}
		rbacRules = append(rbacRules, []string{subject, api.Path, api.Method})
	}

	if len(rbacRules) > 0 {
		_, err = cachedEnforcer.AddPolicies(rbacRules)
		if err != nil {
			return err
		}
	}

	//修改了policy以后，因为用的是cachedEnforcer，所以要清除缓存
	err = cachedEnforcer.InvalidateCache()
	if err != nil {
		return err
	}

	return nil
}

func (r *rbacUpdatePolicyByMenuId) Update() error {
	if r.MenuId == 0 {
		return nil
	}

	//先找到菜单关联的角色id
	var roleIds []int64
	global.DB.Model(&model.RoleAndMenu{}).Where("menu_id = ?", r.MenuId).
		Select("role_id").Find(&roleIds)

	for _, roleId := range roleIds {
		param := rbacUpdatePolicyByRoleId{RoleId: roleId}
		err := param.Update()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *rbacUpdatePolicyByApiId) Update() error {
	if r.ApiId == 0 {
		return nil
	}

	//先找到api关联的菜单id
	var menuIds []int64
	global.DB.Model(&model.MenuAndApi{}).Where("api_id in ?", r.ApiId).
		Select("menu_id").Find(&menuIds)

	for _, menuId := range menuIds {
		param := rbacUpdatePolicyByMenuId{MenuId: menuId}
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
