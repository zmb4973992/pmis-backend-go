package util

import (
	"errors"
	"github.com/go-ldap/ldap/v3"
	"pmis-backend-go/global"
	"strings"
)

var attributes = []string{
	"cn",                //Common Name, 中文名，如：张三
	"distinguishedName", //DN, 区分名，如：[CN=张三,OU=综合管理和法律部,OU=中航国际北京公司,DC=avicbj,DC=ad]
	"sAMAccountName",    //登录名，如：x0020888
	"userPrincipalName", //UPN, 用户主体名称，如：x0020888@avicbj.ad
	"mail",              //邮箱，如：zhangsan@intl-bj.avic.com
}

type UserInfo struct {
	FullName   *string
	Email      *string
	Department *string
}

func LoginByLDAP(username, password string) (permitted bool, userInfo *UserInfo, err error) {
	ldapServer := global.Config.LDAPConfig.Server
	baseDN := global.Config.LDAPConfig.BaseDN
	filter := global.Config.LDAPConfig.Filter
	suffix := global.Config.LDAPConfig.Suffix

	l, err := ldap.DialURL(ldapServer)
	if err != nil {
		return false, nil, err
	}
	defer l.Close()

	err = l.Bind(username+suffix, password)
	if err != nil {
		return false, nil, errors.New("账号密码错误")
	}

	searchRequest := ldap.NewSearchRequest(
		baseDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases,
		0, 0, false, filter, attributes, nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return false, nil, errors.New("搜索失败")
	}

	for i := range sr.Entries {
		entry := sr.Entries[i]
		userPrincipalName := entry.GetAttributeValue("userPrincipalName")
		DN := entry.GetAttributeValue("distinguishedName")

		if username+suffix == userPrincipalName {
			allowedOUs := []string{
				"公司领导", "公司专务", "公司总监", "综合管理和法律部",
				"人力资源和海外机构事务部", "财务管理部", "党建和纪检审计部",
				"储运管理部", "事业部管理委员会和水泥工程事业部", "技术中心",
				"水泥工程市场一部", "水泥工程市场二部", "项目管理部", "工程项目执行部",
				"水泥延伸业务部", "进口部/航空技术部", "成套业务一部", "成套业务二部",
				"成套业务三部", "成套业务四部", "成套业务五部", "成套业务六部",
				"国内企业管理部"}
			for j := range allowedOUs {
				if strings.Contains(DN, allowedOUs[j]) {
					var userInfo UserInfo
					fullName := entry.GetAttributeValue("cn")
					if fullName != "" {
						userInfo.FullName = &fullName
					}

					email := entry.GetAttributeValue("mail")
					if email != "" {
						userInfo.Email = &email
					}

					userInfo.Department = &allowedOUs[j]

					return true, &userInfo, nil
				}
			}
		}
	}
	return false, nil, nil
}
