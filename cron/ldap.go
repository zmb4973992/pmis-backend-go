package cron

import (
	"github.com/go-ldap/ldap/v3"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"strings"
)

func updateUser() {
	ldapServer := global.Config.LDAPConfig.Server
	baseDN := global.Config.LDAPConfig.BaseDN
	filter := global.Config.LDAPConfig.Filter
	permittedOUs := global.Config.LDAPConfig.PermittedOUs
	attributes := global.Config.LDAPConfig.Attributes

	l, err := ldap.DialURL(ldapServer)
	if err != nil {
		global.SugaredLogger.Errorln(err)
	}
	defer l.Close()

	err = l.Bind("z0030975@avicbj.ad", "Bfsu028912")
	if err != nil {
		global.SugaredLogger.Errorln(err)
	}

	searchRequest := ldap.NewSearchRequest(
		baseDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases,
		0, 0, false, filter, attributes, nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		global.SugaredLogger.Errorln(err)
	}

	for i := range sr.Entries {
		entry := sr.Entries[i]
		DN := entry.GetAttributeValue("distinguishedName")
		for j := range permittedOUs {
			if strings.Contains(DN, permittedOUs[j]) {
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

				res := global.DB.Where("username = ?", user.Username).
					FirstOrCreate(&user)
				if res.Error != nil {
					global.SugaredLogger.Errorln(err)
				}

				//给公司领导
				if strings.Contains(DN, "公司领导") ||
					strings.Contains(DN, "公司专务") ||
					strings.Contains(DN, "公司总监") {
					//创建角色
					var roleID int
					global.DB.Model(&model.Role{}).Where("name = ?", "公司级").
						Select("id").First(&roleID)
					var roleAndUser model.RoleAndUser
					roleAndUser.UserID = &user.ID
					roleAndUser.RoleID = &roleID
					global.DB.Where("role_id = ?", roleID).
						Where("user_id = ?", user.ID).
						FirstOrCreate(&roleAndUser)

					//创建部门
					var departmentID int
					global.DB.Model(&model.Organization{}).Where("name = ?", "北京公司").
						Select("id").First(&departmentID)
					var departmentAndUser model.OrganizationAndUser
					departmentAndUser.UserID = &user.ID
					departmentAndUser.OrganizationID = &departmentID
					global.DB.Model(&model.OrganizationAndUser{}).
						Where("department_id = ?", departmentID).
						Where("user_id = ?", user.ID).
						FirstOrCreate(&departmentAndUser)
					//给事业部领导
				} else if strings.Contains(DN, "事业部管理委员会和水泥工程事业部") {
					//创建角色
					var roleID int
					global.DB.Model(&model.Role{}).Where("name = ?", "事业部级").
						Select("id").First(&roleID)
					var roleAndUser model.RoleAndUser
					roleAndUser.UserID = &user.ID
					roleAndUser.RoleID = &roleID
					global.DB.Where("role_id = ?", roleID).
						Where("user_id = ?", user.ID).
						FirstOrCreate(&roleAndUser)

					//创建部门
					var departmentID int
					global.DB.Model(&model.Organization{}).Where("name = ?", "水泥工程事业部").
						Select("id").First(&departmentID)
					var departmentAndUser model.OrganizationAndUser
					departmentAndUser.UserID = &user.ID
					departmentAndUser.OrganizationID = &departmentID
					global.DB.Model(&model.OrganizationAndUser{}).
						Where("department_id = ?", departmentID).
						Where("user_id = ?", user.ID).
						FirstOrCreate(&departmentAndUser)
				} else {
					//创建角色
					var roleID int
					global.DB.Model(&model.Role{}).Where("name = ?", "部门级").
						Select("id").First(&roleID)
					var roleAndUser model.RoleAndUser
					roleAndUser.UserID = &user.ID
					roleAndUser.RoleID = &roleID
					global.DB.Where("role_id = ?", roleID).
						Where("user_id = ?", user.ID).
						FirstOrCreate(&roleAndUser)

					//创建部门
					var departmentID int
					global.DB.Model(&model.Organization{}).Where("name = ?", permittedOUs[j]).
						Select("id").First(&departmentID)
					var departmentAndUser model.OrganizationAndUser
					departmentAndUser.UserID = &user.ID
					departmentAndUser.OrganizationID = &departmentID
					global.DB.Model(&model.OrganizationAndUser{}).
						Where("department_id = ?", departmentID).
						Where("user_id = ?", user.ID).
						FirstOrCreate(&departmentAndUser)
				}
			}
		}
	}
}
