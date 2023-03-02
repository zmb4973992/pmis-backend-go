package main

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"strings"
)

const (
	//ldapServer = "ldap://10.100.10.120:389"
	//baseDN     = "ou=中航国际北京公司,dc=avicbj,dc=ad"
	//filter     = "(&(objectClass=user))"

	ldapServer = "ldap://192.168.172.129:389"
	baseDN     = "ou=中航国际北京公司,dc=avicbj,dc=ad1"
	filter     = "(&(objectClass=user))"
	username   = "zhangsan@avicbj.ad1"
	password   = "Bfsu028912"
)

var attributes = []string{
	"cn",                //Common Name, 中文名，如：张三
	"distinguishedName", //DN, 区分名，如：[CN=张三,OU=综合管理和法律部,OU=中航国际北京公司,DC=avicbj,DC=ad]
	"sAMAccountName",    //登录名，如：x0020888
	"userPrincipalName", //UPN, 用户主体名称，如：x0020888@avicbj.ad
	"mail",              //邮箱，如：zhangsan@intl-bj.avic.com
}

func main() {
	l, err := ldap.DialURL(ldapServer)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	err = l.Bind(username, password)
	if err != nil {
		panic(err)
	}

	searchRequest := ldap.NewSearchRequest(
		baseDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases,
		0, 0, false, filter, attributes, nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		panic(err)
	}

	for i := range sr.Entries {
		entry := sr.Entries[i]
		userPrincipalName := entry.GetAttributeValue("userPrincipalName")
		DN := entry.GetAttributeValue("distinguishedName")
		if username == userPrincipalName {
			allowedOUs := []string{"综合管理和法律部", "一部"}
			var permitted bool
			for j := range allowedOUs {
				allowedOU := allowedOUs[j]
				if strings.Contains(DN, allowedOU) {
					permitted = true
					break
				}
			}
			if permitted {
				fmt.Println("登录成功")
			} else {
				fmt.Println("登录失败")
			}
		} else {
			fmt.Println("登录失败")
		}

	}

}
