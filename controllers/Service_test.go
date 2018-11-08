package controllers

import (
	"crypto/sha256"
	"fmt"
	"net/url"
	"testing"
)

func TestLoginController_Post(t *testing.T) {
	hash := []byte(fmt.Sprint(sha256.Sum256([]byte("123456789"))))

	userdata := url.Values{
		"UserName":       {"UN:roy"},
		"UserPasswordHash": {string(hash)},
	}
	res := Dopost(userdata,"http://localhost:8000/login")
	fmt.Println(res)
}
