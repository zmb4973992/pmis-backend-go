package util

import (
	"errors"
	"github.com/go-ldap/ldap/v3"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
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
	ID           int
	FullName     *string
	Email        *string
	Organization *string
}

func LoginByLDAP(username, password string) (permitted bool, err error) {
	//以下这段为测试专用，记得删除
	{
		if username == "a" && password == "a" {
			var user model.User
			err = global.DB.Model(&model.User{}).Where(model.User{Username: "a"}).
				First(&user).Error
			if err != nil {
				global.SugaredLogger.Errorln(err)
				return false, err
			}
			return true, nil
		}
	}
	//以上为测试专用，记得删除

	ldapServer := global.Config.LDAPConfig.Server
	baseDN := global.Config.LDAPConfig.BaseDN
	filter := global.Config.LDAPConfig.Filter
	suffix := global.Config.LDAPConfig.Suffix

	l, err := ldap.DialURL(ldapServer)
	if err != nil {
		return false, err
	}
	defer l.Close()

	err = l.Bind(username+suffix, password)
	if err != nil {
		return false, errors.New("账号密码错误")
	}

	searchRequest := ldap.NewSearchRequest(
		baseDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases,
		0, 0, false, filter, attributes, nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return false, errors.New("搜索失败")
	}

	for i := range sr.Entries {
		entry := sr.Entries[i]
		userPrincipalName := entry.GetAttributeValue("userPrincipalName")
		DN := entry.GetAttributeValue("distinguishedName")

		if username+suffix == userPrincipalName {
			permittedOUs := global.Config.LDAPConfig.PermittedOUs
			for j := range permittedOUs {
				if strings.Contains(DN, permittedOUs[j]) {
					//var user UserInfo
					//fullName := entry.GetAttributeValue("cn")
					//if fullName != "" {
					//	user.FullName = &fullName
					//}
					//
					//email := entry.GetAttributeValue("mail")
					//if email != "" {
					//	user.Email = &email
					//}
					//
					//user.Organization = &permittedOUs[j]

					return true, nil
				}
			}
		}
	}
	return false, nil
}
