package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role     string `json:"role"`

}

func NewJWTManger(secretkey string,tokenduration time.Duration)*JWTManager{
	return &JWTManager{
		secretKey: secretkey,
		tokenDuration: tokenduration,
	}
}

func (manager *JWTManager)Generate(user *User)(string,error){
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(manager.tokenDuration).Unix(),
		},
		Username: user.Name,
		Role: user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	return token.SignedString([]byte(manager.secretKey))

}

func (manager *JWTManager)Verify(accessToken string)(*UserClaims,error){

	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func (token *jwt.Token)(interface{},error){
			_,ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil,fmt.Errorf("unexpected token signing Method")
			}
			return []byte(manager.secretKey),nil
		},
	)
	if err != nil {
		return nil,fmt.Errorf("invalid token: %v",err)
	}

	claims ,ok:= token.Claims.(*UserClaims)
	if !ok {
		return nil,fmt.Errorf("invalid token claims")
	}
	return claims,nil
}

