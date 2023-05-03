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

type RBACUpdate struct {
	LastModifier int
	RoleIDs      []int `json:"role_ids,omitempty"`
	MenuIDs      []int `json:"menu_ids,omitempty"`
	ApiIDs       []int `json:"api_ids,omitempty"`
}

func (r *RBACUpdate) Update() error {
	//如果需要更新角色id的权限
	err := updateRBACRulesByRoleIDs(r.RoleIDs)
	if err != nil {
		return err
	}

	//如果需要更新菜单的权限
	err = updateRBACRulesByMenuIDs(r.MenuIDs)
	if err != nil {
		return err
	}

	//如果需要更新api的权限
	err = updateRBACRulesByApiIDs(r.ApiIDs)
	if err != nil {
		return err
	}

	return nil
}

func updateRBACRulesByRoleIDs(roleIDs []int) error {
	if len(roleIDs) == 0 {
		return nil
	}

	cachedEnforcer, err := util.NewCachedEnforcer()
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return err
	}

	for _, roleID := range roleIDs {
		subject := strconv.Itoa(roleID)
		_, err = cachedEnforcer.RemoveFilteredPolicy(0, subject)
		if err != nil {
			global.SugaredLogger.Errorln(err)
			return err
		}

		//找到角色拥有的菜单
		var menuIDs []int
		global.DB.Model(&model.RoleAndMenu{}).Where("role_id = ?", roleID).
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
	}
	return nil
}

func updateRBACRulesByMenuIDs(menuIDs []int) error {
	if len(menuIDs) == 0 {
		return nil
	}

	//先找到菜单关联的角色id
	var roleIDs []int
	global.DB.Model(&model.RoleAndMenu{}).Where("menu_id in ?", menuIDs).
		Select("role_id").Find(&roleIDs)
	err := updateRBACRulesByRoleIDs(roleIDs)
	if err != nil {
		return err
	}
	return nil
}

func updateRBACRulesByApiIDs(apiIDs []int) error {
	if len(apiIDs) == 0 {
		return nil
	}

	//先找到api关联的菜单id
	var menuIDs []int
	global.DB.Model(&model.MenuAndApi{}).Where("api_id in ?", apiIDs).
		Select("menu_id").Find(&menuIDs)
	err := updateRBACRulesByMenuIDs(menuIDs)
	if err != nil {
		return err
	}
	return nil
}
