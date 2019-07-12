package wrapper

import (
	"bytes"
	"encoding/gob"
)

type UserData struct {
	UserName             string		//"UN:xxxxxx"
	UserPasswordHash     []byte		//"xxxxxxxxx"
	ChangePasswordPolicy string		//"CPP:xxxxx"
	GetTipPolicy         string		//"GTP:xxxxx"
	GetTipMessage        string		//"GTM:xxxxx"
	UserAttributes       []string	//"xxxxxxxxx"
	SpecialAAId			 string		//"AA_1"
	PartSk               [][]byte   //nil
	Aid                  [][]byte	//nil
}

type CompanyData struct {
	CompanyName string //"COM:xxxxxxx"
	DataName string //"DN:xxxxxxx"
	Data string //"DATA:xxxxxxx"
	AuthPolicy string //"AP:xxxxxxx"
	SpecialAAId			 string		//"AA_1"
	PartSk               [][]byte   //nil
}

type AuthList struct {
	Auths []string
}

func (u *UserData) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(u)
	if err != nil {
		return []byte{}, err
	}

	return result.Bytes(), nil
}

func (u *CompanyData) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(u)
	if err != nil {
		return []byte{}, err
	}

	return result.Bytes(), nil
}

func DeserializeUserData(d []byte) (*UserData, error) {
	ud := new(UserData)

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&ud)
	if err != nil {
		return nil, err
	}
	return ud, nil
}

func (u *AuthList) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(u)
	if err != nil {
		return []byte{}, err
	}

	return result.Bytes(), nil
}

func DeserializeAuthList(d []byte) (*AuthList, error) {
	ud := new(AuthList)

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&ud)
	if err != nil {
		return nil, err
	}
	return ud, nil
}

func DeserializeCompanyData(d []byte) (*CompanyData, error) {
	ud := new(CompanyData)

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&ud)
	if err != nil {
		return nil, err
	}
	return ud, nil
}
