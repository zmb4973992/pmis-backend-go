package windows_ad

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

func UpdateUsersForCron() {
	err := UpdateUsersByLDAP()
	if err != nil {
		param := service.ErrorLogCreate{
			Detail: err.Error(),
		}
		param.Create()
	}
}

func UpdateUsersByLDAP() error {
	ldapServer := global.Config.LDAPConfig.Server
	baseDN := global.Config.LDAPConfig.BaseDN
	filter := global.Config.LDAPConfig.Filter
	suffix := global.Config.LDAPConfig.Suffix
	account := global.Config.LDAPConfig.Account
	password := global.Config.LDAPConfig.Password
	permittedOUs := global.Config.LDAPConfig.PermittedOUs
	attributes := global.Config.LDAPConfig.Attributes

	l, err := ldap.DialURL(ldapServer)
	if err != nil {
		return err
	}

	defer l.Close()

	err = l.Bind(account+suffix, password)
	if err != nil {
		return err
	}

	searchRequest := ldap.NewSearchRequest(
		baseDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases,
		0, 0, false, filter, attributes, nil,
	)

	searchResult, err1 := l.Search(searchRequest)
	if err1 != nil {
		return err1
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
					global.SugaredLogger.Errorln(err)
					param := service.ErrorLogCreate{Detail: err.Error()}
					param.Create()
					continue
				}

				//部门信息不从LDAP导入，因为LDAP不是严格按照部门进行设置的

				if permittedOU == "公司领导" || permittedOU == "公司总监" ||
					permittedOU == "公司专务" {
					//添加用户和组织的关联
					var organization model.Organization
					err = global.DB.Where("name = '公司领导'").First(&organization).Error
					if err != nil {
						global.SugaredLogger.Errorln(err)
						param := service.ErrorLogCreate{
							Detail: "在组织机构表中找不到：公司领导",
						}
						param.Create()
						continue
					}

					record1 := model.OrganizationAndUser{
						UserID:         user.ID,
						OrganizationID: organization.ID,
						ImportedByLDAP: model.BoolToPointer(true),
					}
					err = global.DB.
						Where("organization_id = ?", organization.ID).
						Where("user_id = ?", user.ID).
						FirstOrCreate(&record1).Error
					if err != nil {
						global.SugaredLogger.Errorln(err)
						param := service.ErrorLogCreate{
							Detail: err.Error(),
						}
						param.Create()
						continue
					}

					//添加用户和数据权限的关联
					var dataAuthority model.DataAuthority
					err = global.DB.Where("name = '所有部门'").First(&dataAuthority).Error
					if err != nil {
						global.SugaredLogger.Errorln(err)
						param := service.ErrorLogCreate{
							Detail: "在数据权限表中找不到：所有部门",
						}
						param.Create()
						continue
					}

					record2 := model.UserAndDataAuthority{
						UserID:          user.ID,
						DataAuthorityID: dataAuthority.ID,
						ImportedByLDAP:  model.BoolToPointer(true),
					}
					err = global.DB.
						Where("data_authority_id = ?", dataAuthority.ID).
						Where("user_id = ?", user.ID).
						FirstOrCreate(&record2).Error
					if err != nil {
						global.SugaredLogger.Errorln(err)
						param := service.ErrorLogCreate{
							Detail: err.Error(),
						}
						param.Create()
						continue
					}

				} else if permittedOU == "事业部管理委员会和水泥工程事业部" {
					//添加用户和组织的关联
					var organization model.Organization
					err = global.DB.Where("name = ?", "水泥工程事业部").First(&organization).Error
					if err != nil {
						global.SugaredLogger.Errorln(err)
						param := service.ErrorLogCreate{
							Detail: "在组织机构表中找不到：水泥工程事业部",
						}
						param.Create()
						continue
					}
					record := model.OrganizationAndUser{
						UserID:         user.ID,
						OrganizationID: organization.ID,
						ImportedByLDAP: model.BoolToPointer(true),
					}
					err = global.DB.
						Where("organization_id = ?", organization.ID).
						Where("user_id = ?", user.ID).
						FirstOrCreate(&record).Error
					if err != nil {
						global.SugaredLogger.Errorln(err)
						param := service.ErrorLogCreate{
							Detail: err.Error(),
						}
						param.Create()
						continue
					}

					//添加用户和数据权限的关联
					var dataAuthority model.DataAuthority
					err = global.DB.Where("name = '所属部门和子部门'").First(&dataAuthority).Error
					if err != nil {
						global.SugaredLogger.Errorln(err)
						param := service.ErrorLogCreate{
							Detail: "在数据权限表中找不到：所属部门和子部门",
						}
						param.Create()
						continue
					}

					record2 := model.UserAndDataAuthority{
						UserID:          user.ID,
						DataAuthorityID: dataAuthority.ID,
						ImportedByLDAP:  model.BoolToPointer(true),
					}
					err = global.DB.
						Where("data_authority_id = ?", dataAuthority.ID).
						Where("user_id = ?", user.ID).
						FirstOrCreate(&record2).Error
					if err != nil {
						global.SugaredLogger.Errorln(err)
						param := service.ErrorLogCreate{
							Detail: err.Error(),
						}
						param.Create()
						continue
					}

				} else {
					//添加用户和组织的关联
					var organization model.Organization
					err = global.DB.Where("name = ?", permittedOU).First(&organization).Error
					if err != nil {
						global.SugaredLogger.Errorln(err)
						param := service.ErrorLogCreate{Detail: "在组织机构表中找不到：" + permittedOU}
						param.Create()
						continue
					}
					record1 := model.OrganizationAndUser{
						UserID:         user.ID,
						OrganizationID: organization.ID,
						ImportedByLDAP: model.BoolToPointer(true),
					}
					err = global.DB.
						Where("organization_id = ?", organization.ID).
						Where("user_id = ?", user.ID).
						FirstOrCreate(&record1).Error
					if err != nil {
						global.SugaredLogger.Errorln(err)
						param := service.ErrorLogCreate{
							Detail: err.Error(),
						}
						param.Create()
						continue
					}

					//添加用户和数据权限的关联
					var dataAuthority model.DataAuthority
					err = global.DB.Where("name = '所属部门和子部门'").First(&dataAuthority).Error
					if err != nil {
						global.SugaredLogger.Errorln(err)
						param := service.ErrorLogCreate{
							Detail: "在数据权限表中找不到：所属部门和子部门",
						}
						param.Create()
						continue
					}

					record2 := model.UserAndDataAuthority{
						UserID:          user.ID,
						DataAuthorityID: dataAuthority.ID,
						ImportedByLDAP:  model.BoolToPointer(true),
					}
					err = global.DB.
						Where("data_authority_id = ?", dataAuthority.ID).
						Where("user_id = ?", user.ID).
						FirstOrCreate(&record2).Error
					if err != nil {
						global.SugaredLogger.Errorln(err)
						param := service.ErrorLogCreate{
							Detail: err.Error(),
						}
						param.Create()
						continue
					}
				}
			}
		}
	}

	err = setAdmin()
	if err != nil {
		return err
	}

	return nil
}

func setAdmin() error {
	var user model.User
	err := global.DB.Where("username = 'z0030975'").
		First(&user).Error
	if err != nil {
		return err
	}

	fullName := "周梦斌"
	user.FullName = &fullName
	err = global.DB.Save(&user).Error
	if err != nil {
		return err
	}

	var dataAuthority model.DataAuthority
	err = global.DB.Where("name = '所有部门'").
		First(&dataAuthority).Error
	if err != nil {
		return err
	}

	global.DB.Where("user_id = ?", user.ID).
		Delete(&model.UserAndDataAuthority{})

	var userAndAuthority model.UserAndDataAuthority
	userAndAuthority.UserID = user.ID
	userAndAuthority.DataAuthorityID = dataAuthority.ID
	userAndAuthority.ImportedByLDAP = model.BoolToPointer(true)

	err = global.DB.Create(&userAndAuthority).Error
	if err != nil {
		return err
	}

	return nil
}
