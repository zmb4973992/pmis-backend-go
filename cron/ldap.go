package cron

import (
	"github.com/go-ldap/ldap/v3"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/service"
	"strings"
)

//目前存在一个问题，就是导入用户后，用户会一直有效。
//后期需要增加用户有效性校验的函数。
//可以把从ldap读取到的用户列表放到临时表，然后把现在的用户表和临时表进行比对

func updateUsersByLDAP() {
	ldapServer := global.Config.LDAPConfig.Server
	baseDN := global.Config.LDAPConfig.BaseDN
	filter := global.Config.LDAPConfig.Filter
	suffix := global.Config.LDAPConfig.Suffix
	account := global.Config.LDAPConfig.Account
	password := global.Config.LDAPConfig.Password
	permittedOUs := global.Config.LDAPConfig.PermittedOUs
	attributes := global.Config.LDAPConfig.Attributes

	//fmt.Println("ldapServer:", ldapServer)
	//fmt.Println("baseDN:", baseDN)
	//fmt.Println("filter:", filter)
	//fmt.Println("suffix:", suffix)
	//fmt.Println("account:", account)
	//fmt.Println("password:", password)
	//fmt.Println("permittedOUs:", permittedOUs)
	//fmt.Println("attributes:", attributes)

	l, err := ldap.DialURL(ldapServer)
	if err != nil {
		param := service.ErrorLogCreate{
			Detail: err.Error(),
		}
		param.Create()
		return
	}

	defer l.Close()

	err = l.Bind(account+suffix, password)
	if err != nil {
		param := service.ErrorLogCreate{
			Detail: err.Error(),
		}
		param.Create()
		return
	}

	searchRequest := ldap.NewSearchRequest(
		baseDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases,
		0, 0, false, filter, attributes, nil,
	)

	searchResult, err1 := l.Search(searchRequest)
	if err1 != nil {
		param := service.ErrorLogCreate{
			Detail: err.Error(),
		}
		param.Create()
		return
	}

	//把组织-用户中间表中的ldap导入数据删除
	global.DB.Where("imported_by_ldap = ?", 1).Delete(&model.OrganizationAndUser{})

	for i := range searchResult.Entries {
		entry := searchResult.Entries[i]
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
					FirstOrCreate(&user).Error

				if err != nil {
					param := service.ErrorLogCreate{
						Detail: err.Error(),
					}
					param.Create()
					return
				}

				//部门信息不从LDAP导入，因为LDAP不是严格按照部门进行设置的

				//添加组织机构和用户的关联
				if permittedOU == "公司领导" ||
					permittedOU == "公司总监" ||
					permittedOU == "公司专务" {
					var organization model.Organization
					err = global.DB.Where("name = ?", "北京公司").First(&organization).Error
					if err != nil {
						param := service.ErrorLogCreate{
							Detail: err.Error(),
						}
						param.Create()
						return
					}
					record := model.OrganizationAndUser{
						UserID:         user.ID,
						OrganizationID: organization.ID,
					}
					err = global.DB.Where("organization_id = ?", organization.ID).
						Where("user_id = ?", user.ID).
						Attrs(&model.OrganizationAndUser{
							ImportedByLDAP: model.BoolToPointer(true),
						}).
						FirstOrCreate(&record).Error
					if err != nil {
						param := service.ErrorLogCreate{
							Detail: err.Error(),
						}
						param.Create()
						return
					}

				} else if permittedOU == "事业部管理委员会和水泥工程事业部" {
					var organization model.Organization
					err = global.DB.Where("name = ?", "水泥工程事业部").First(&organization).Error
					if err != nil {
						param := service.ErrorLogCreate{
							Detail: err.Error(),
						}
						param.Create()
						return
					}
					record := model.OrganizationAndUser{
						UserID:         user.ID,
						OrganizationID: organization.ID,
					}
					err = global.DB.Where("organization_id = ?", organization.ID).
						Where("user_id = ?", user.ID).
						Attrs(&model.OrganizationAndUser{
							ImportedByLDAP: model.BoolToPointer(true),
						}).
						FirstOrCreate(&record).Error
					if err != nil {
						param := service.ErrorLogCreate{
							Detail: err.Error(),
						}
						param.Create()
						return
					}

				} else {
					var organization model.Organization
					err = global.DB.Where("name = ?", permittedOU).First(&organization).Error
					if err != nil {
						global.SugaredLogger.Errorln(err)
						return
					}
					record := model.OrganizationAndUser{
						UserID:         user.ID,
						OrganizationID: organization.ID,
					}
					err = global.DB.Where("organization_id = ?", organization.ID).
						Where("user_id = ?", user.ID).
						Attrs(&model.OrganizationAndUser{
							ImportedByLDAP: model.BoolToPointer(true),
						}).
						FirstOrCreate(&record).Error
					if err != nil {
						global.SugaredLogger.Errorln(err)
						return
					}
				}
			}
		}
	}

	return
}
