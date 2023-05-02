package util

import (
	"github.com/casbin/casbin/v2"
	gormAdapter "github.com/casbin/gorm-adapter/v3"
	"pmis-backend-go/global"
	"time"
)

// 待删
func NewEnforcer() (enforcer *casbin.Enforcer, err error) {
	adapter, err := gormAdapter.NewAdapterByDB(global.DB)
	if err != nil {
		return nil, err
	}
	enforcer, err = casbin.NewEnforcer("./config/casbin-model.conf", adapter)
	if err != nil {
		return nil, err
	}
	err = enforcer.LoadPolicy()
	if err != nil {
		return nil, err
	}
	return enforcer, nil
}

func NewCachedEnforcer() (cachedEnforcer *casbin.CachedEnforcer, err error) {
	adapter, err := gormAdapter.NewAdapterByDB(global.DB)
	if err != nil {
		return nil, err
	}
	cachedEnforcer, err = casbin.NewCachedEnforcer("./config/casbin-model.conf", adapter)
	if err != nil {
		return nil, err
	}
	cachedEnforcer.SetExpireTime(24 * time.Hour)
	err = cachedEnforcer.LoadPolicy()
	if err != nil {
		return nil, err
	}
	return cachedEnforcer, nil
}
