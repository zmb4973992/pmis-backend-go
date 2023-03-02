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
	attributes := []string{
		"cn",                //Common Name, 中文名，如：张三
		"distinguishedName", //DN, 区分名，如：[CN=张三,OU=综合管理和法律部,OU=中航国际北京公司,DC=avicbj,DC=ad]
		"sAMAccountName",    //登录名，如：x0020888、zhangsan
		"userPrincipalName", //UPN, 用户主体名称，如：x0020888@avicbj.ad
		"mail",              //邮箱，如：zhangsan@intl-bj.avic.com
	}

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
		allowedOUs := []string{
			"公司领导", "公司专务", "公司总监", "综合管理和法律部",
			"人力资源和海外机构事务部", "财务管理部", "党建和纪检审计部",
			"储运管理部", "事业部管理委员会和水泥工程事业部", "技术中心",
			"水泥工程市场一部", "水泥工程市场二部", "项目管理部", "工程项目执行部",
			"水泥延伸业务部", "进口部/航空技术部", "成套业务一部", "成套业务二部",
			"成套业务三部", "成套业务四部", "成套业务五部", "成套业务六部", "国内企业管理部",
		}
		for j := range allowedOUs {
			if strings.Contains(DN, allowedOUs[j]) {
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
					global.DB.Model(&model.Department{}).Where("name = ?", "北京公司").
						Select("id").First(&departmentID)
					var departmentAndUser model.DepartmentAndUser
					departmentAndUser.UserID = &user.ID
					departmentAndUser.DepartmentID = &departmentID
					global.DB.Model(&model.DepartmentAndUser{}).
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
					global.DB.Model(&model.Department{}).Where("name = ?", "水泥工程事业部").
						Select("id").First(&departmentID)
					var departmentAndUser model.DepartmentAndUser
					departmentAndUser.UserID = &user.ID
					departmentAndUser.DepartmentID = &departmentID
					global.DB.Model(&model.DepartmentAndUser{}).
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
					global.DB.Model(&model.Department{}).Where("name = ?", allowedOUs[j]).
						Select("id").First(&departmentID)
					var departmentAndUser model.DepartmentAndUser
					departmentAndUser.UserID = &user.ID
					departmentAndUser.DepartmentID = &departmentID
					global.DB.Model(&model.DepartmentAndUser{}).
						Where("department_id = ?", departmentID).
						Where("user_id = ?", user.ID).
						FirstOrCreate(&departmentAndUser)
				}

			}
		}
	}
	//fmt.Println("用户更新完毕")
}
