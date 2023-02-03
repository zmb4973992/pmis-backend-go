package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	password, err := bcrypt.GenerateFromPassword([]byte("kjsdf"), 20)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(password)
	fmt.Println("3kj")
}
