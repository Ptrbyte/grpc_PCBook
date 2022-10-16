package service

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name          string
	HashPassword  string
	Role          string
}

func NewUser(name string,password string,role string)(*User,error) {

	hashPassword,err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if err != nil {
		return nil,fmt.Errorf("cannot hash Password: %v",err)
	}

	user := &User{
		Name:name,
		HashPassword:string(hashPassword),
		Role:role,
	}
	return user,nil
}

func (user *User)IsCorrectPassWord(passsword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(passsword))
	return err == nil
}

func (user *User)Clone()*User{
	return &User{
		Name: user.Name,
		HashPassword: user.HashPassword,
		Role: user.Role,
	}
}
