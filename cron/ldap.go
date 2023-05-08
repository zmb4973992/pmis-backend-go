package cron

import (
	"github.com/go-ldap/ldap/v3"
	"github.com/yitter/idgenerator-go/idgen"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"strings"
)

//目前存在一个问题，就是导入用户后，用户会一直有效。
//后期需要增加用户有效性校验的函数。
//可以把从ldap读取到的用户列表放到临时表，然后把现在的用户表和临时表进行比对

func updateUsers() error {
	ldapServer := global.Config.LDAPConfig.Server
	baseDN := global.Config.LDAPConfig.BaseDN
	filter := global.Config.LDAPConfig.Filter
	suffix := global.Config.LDAPConfig.Suffix
	account := global.Config.LDAPConfig.Account
	password := global.Config.LDAPConfig.Password
	permittedOUs := global.Config.LDAPConfig.PermittedOUs
	attributes := global.Config.LDAPConfig.Attributes

	l, err := ldap.DialURL(ldapServer)
	defer l.Close()
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return err
	}

	err = l.Bind(account+suffix, password)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return err
	}

	searchRequest := ldap.NewSearchRequest(
		baseDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases,
		0, 0, false, filter, attributes, nil,
	)

	sr, err1 := l.Search(searchRequest)
	if err1 != nil {
		global.SugaredLogger.Errorln(err)
		return err
	}

	//把组织-用户中间表中的ldap导入数据删除
	global.DB.Where("imported_by_ldap = ?", 1).Delete(&model.OrganizationAndUser{})

	for i := range sr.Entries {
		entry := sr.Entries[i]
		DN := entry.GetAttributeValue("distinguishedName")
		for _, permittedOU := range permittedOUs {
			if strings.Contains(DN, permittedOU) {
				//添加用户
				var user model.User
				user.Username = entry.GetAttributeValue("sAMAccountName")

				isValid := true
				user.IsValid = &isValid

				fullName := entry.GetAttributeValue("cn")
				if fullName != "" {
					user.FullName = &fullName
				}

				email := entry.GetAttributeValue("mail")
				if email != "" {
					user.EmailAddress = &email
				}

				//添加用户信息
				err = global.DB.Where("username = ?", user.Username).
					Attrs(&model.User{
						BasicModel: model.BasicModel{
							SnowID: idgen.NextId(),
						}}).
					FirstOrCreate(&user).Error

				if err != nil {
					global.SugaredLogger.Errorln(err)
					return err
				}

				//部门信息不从LDAP导入，因为LDAP不是严格按照部门进行设置的

				//添加组织机构和用户的关联
				if permittedOU == "公司领导" ||
					permittedOU == "公司总监" ||
					permittedOU == "公司专务" {
					var organization model.Organization
					err = global.DB.Where("name = ?", "北京公司").First(&organization).Error
					if err != nil {
						global.SugaredLogger.Errorln(err)
						return err
					}
					record := model.OrganizationAndUser{
						UserSnowID:         user.SnowID,
						OrganizationSnowID: organization.SnowID,
					}
					err = global.DB.Where("organization_snow_id = ?", organization.SnowID).
						Where("user_snow_id = ?", user.SnowID).
						Attrs(&model.OrganizationAndUser{
							ImportedByLDAP: model.BoolToPointer(true),
							BasicModel: model.BasicModel{
								SnowID: idgen.NextId(),
							}}).
						FirstOrCreate(&record).Error
					if err != nil {
						global.SugaredLogger.Errorln(err)
						return err
					}

				} else if permittedOU == "事业部管理委员会和水泥工程事业部" {
					var organization model.Organization
					err = global.DB.Where("name = ?", "水泥工程事业部").First(&organization).Error
					if err != nil {
						global.SugaredLogger.Errorln(err)
						return err
					}
					record := model.OrganizationAndUser{
						UserSnowID:         user.SnowID,
						OrganizationSnowID: organization.SnowID,
					}
					err = global.DB.Where("organization_snow_id = ?", organization.SnowID).
						Where("user_snow_id = ?", user.SnowID).
						Attrs(&model.OrganizationAndUser{
							ImportedByLDAP: model.BoolToPointer(true),
							BasicModel: model.BasicModel{
								SnowID: idgen.NextId(),
							}}).
						FirstOrCreate(&record).Error
					if err != nil {
						global.SugaredLogger.Errorln(err)
						return err
					}

				} else {
					var organization model.Organization
					err = global.DB.Where("name = ?", permittedOU).First(&organization).Error
					if err != nil {
						global.SugaredLogger.Errorln(err)
						return err
					}
					record := model.OrganizationAndUser{
						UserSnowID:         user.SnowID,
						OrganizationSnowID: organization.SnowID,
					}
					err = global.DB.Where("organization_snow_id = ?", organization.SnowID).
						Where("user_snow_id = ?", user.SnowID).
						Attrs(&model.OrganizationAndUser{
							ImportedByLDAP: model.BoolToPointer(true),
							BasicModel: model.BasicModel{
								SnowID: idgen.NextId(),
							}}).
						FirstOrCreate(&record).Error
					if err != nil {
						global.SugaredLogger.Errorln(err)
						return err
					}
				}
			}
		}
	}
	return nil
}
