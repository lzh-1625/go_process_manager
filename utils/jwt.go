package utils

import (
	"errors"
	"time"

	"github.com/lzh-1625/go_process_manager/config"

	"github.com/golang-jwt/jwt"
)

var mySecret []byte

func SetSecret(secret []byte) {
	mySecret = secret
}

func keyFunc(_ *jwt.Token) (i interface{}, err error) {
	return mySecret, nil
}

type MyClaims struct {
	UserName string `json:"user_name"`
	jwt.StandardClaims
}

func GenToken(UserName string) (string, error) {
	// 创建一个我们自己的声明的数据
	c := MyClaims{
		UserName,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(
				time.Duration(config.CF.TokenExpirationTime) * time.Hour).Unix(), // 过期时间
			Issuer: "jwt", // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(mySecret)
}

func ParseToken(tokenString string) (*MyClaims, error) {
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, keyFunc)
	if err != nil {
		return nil, err
	}
	if token.Valid {
		return mc, nil
	}
	return nil, errors.New("invalid token")
}

func RefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
	if _, err = jwt.Parse(rToken, keyFunc); err != nil {
		return
	}
	var claims MyClaims
	_, err = jwt.ParseWithClaims(aToken, &claims, keyFunc)
	v, _ := err.(*jwt.ValidationError)

	if v.Errors == jwt.ValidationErrorExpired {
		token, _ := GenToken(claims.UserName)
		return token, "", nil
	}
	return
}
