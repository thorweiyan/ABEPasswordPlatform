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
		"UserName":             {"UN:roy"},
		"UserPasswordHash":     {hash},
		"ChangePasswordPolicy": {"(UN:roy AND ZS:shuai AND ZS:hei AND SFZ:678987000236787654 AND SJ:17317301908)"},
		"GetTipPolicy":         {"(UN:roy AND SFZ:678987000236787654 AND SJ:17317301908)"},
		"GetTipMessage":        {"shuai & hei"},
		"UserAttributes":       {"UN:roy,SFZ:678987000236787654,SJ:17317301908,YX:zry_nuaa@897.com,ZS:shuai,ZS:hei"},
	}

	res := Dopost(userdata,"http://localhost:8000/signup")
	fmt.Println(res)
}

func TestChangePasswordController_Post(t *testing.T) {

	hash := []byte(fmt.Sprint(sha256.Sum256([]byte("123456789"))))
	userdata := url.Values{
		"UserName":         {"UN:roy"},
		"UserPasswordHash": {string(hash)},
		"UserAttributes":   {"UN:roy,SFZ:678987000236787654,SJ:17317301908,ZS:shuai,ZS:hei"},
	}

	res := Dopost(userdata,"http://localhost:8000/changepassword")
	fmt.Println(res)
}

func TestGetTipController_Post(t *testing.T) {
	userdata := url.Values{
		"UserName":       {"UN:roy"},
		"UserAttributes": {"UN:roy,SFZ:678987000236787654,SJ:17317301908"},
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
