package main

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte("MySigningKey")

func GenerateJwt() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256) //Generating the Token
	claims := token.Claims.(jwt.MapClaims)   // We are getting the claims out of it

	claims["authorized"] = true
	claims["user"] = "Elliot Forbes"
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey) // signing the token via signing Key private
	if err != nil {
		fmt.Errorf("Somthing went Wrong!! %s", err.Error())
		return "", err
	}
	return tokenString, err
}
func main() {
	fmt.Println("Hello Worlds from Rudraksh")
	tokenString, err := GenerateJwt()
	if err != nil {
		fmt.Println("Error while generating the Token")
	}
	fmt.Println(tokenString)
}
