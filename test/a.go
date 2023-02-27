package main

import (
	"github.com/go-ldap/ldap/v3"
)

func main() {
	l, err := ldap.DialURL("ldap://10.100.10.120:389")
	if err != nil {
		panic(err)
	}

	defer l.Close()

}
