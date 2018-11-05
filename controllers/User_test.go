package controllers

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"testing"
)

type data struct {
	UserName string
	UserPasswordHash string
	ChangePasswordPolicy string
	GetTipPolicy string
	GetTipMessage string
	UserAttributes string
}

func TestRand2(t *testing.T)  {
	for i:=0 ;i<10;i++  {
		fmt.Println(Rand2(big.NewInt(2)).Int64())
	}
}

func TestSignUpController_Post(t *testing.T) {

	hash := string([]byte(fmt.Sprint(sha256.Sum256([]byte("123456")))))

	userdata := url.Values{
		"UserName":             {"roy"},
		"UserPasswordHash":     {hash},
		"ChangePasswordPolicy": {"roy AND shuai AND hei AND 0123456789 AND 9876543210"},
		"GetTipPolicy":         {"roy AND 0123456789 AND 9876543210"},
		"GetTipMessage":        {"shuai & hei"},
		"UserAttributes":       {"roy,678987000236787654,17317301908,zry_nuaa@897.com,shuai,hei"},
	}

	res := Dopost(userdata,"http://localhost:8000/signup")
	fmt.Println(res)
}

func TestChangePasswordController_Post(t *testing.T) {

	hash := []byte(fmt.Sprint(sha256.Sum256([]byte("123456789"))))
	userdata := url.Values{
		"UserName":         {"roy"},
		"UserPasswordHash": {string(hash)},
		"UserAttributes":   {"roy,678987000236787654,17317301908,shuai,hei"},
	}

	res := Dopost(userdata,"http://localhost:8000/changepassword")
	fmt.Println(res)
}

func TestGetTipController_Post(t *testing.T) {
	userdata := url.Values{
		"UserName":       {"roy"},
		"UserAttributes": {"shuai,hei"},
	}
	res := Dopost(userdata,"http://localhost:8000/gettip")
	fmt.Println(res)
}

func Dopost(postValue url.Values, urlstr string) (content string) {
	resp, err := http.PostForm(urlstr, postValue)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	return  string(body)
}
