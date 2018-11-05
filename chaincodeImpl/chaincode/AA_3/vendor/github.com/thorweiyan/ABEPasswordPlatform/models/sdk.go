package models

import (
	"fmt"
	"encoding/hex"
	"github.com/thorweiyan/fabric_go_sdk"
	"os"
	"strconv"
	"github.com/thorweiyan/ABEPasswordPlatform/chaincodeImpl/wrapper"
)


var fSetup = fabric_go_sdk.FabricSetup{
	// Network parameters
	OrdererID: "orderer.fudan.edu.cn",
	OrgID: "org1.fudan.edu.cn",

	// Channel parameters
	ChannelID:     "fudanfabric",
	ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/thorweiyan/fabric_go_sdk/fixtures/artifacts/fudanfabric.channel.tx",

	// Chaincode parameters
	ChainCodeID:     "",
	ChaincodeGoPath: os.Getenv("GOPATH"),
	ChaincodePath:   "github.com/thorweiyan/ABEPasswordPlatform/chaincodeImpl/chaincode/",
	ChaincodeVersion: "2",
	OrgAdmin:        "Admin",
	OrgName:         "org1",
	ConfigFile:      os.Getenv("GOPATH") + "/src/github.com/thorweiyan/fabric_go_sdk/config.yaml",

	// User parameters
	UserName: "User1",
}
func OneButtonStartUp() {
	aa1pubkey,_ := hex.DecodeString("04e3dd49e00dce869da09afc266d707a3e59377d28aded9d8f264ab890790aa92735f7ca9df8507a1d0823092e29ab7d74336dd9938521c479")
	aa1prikey,_ := hex.DecodeString("3068020101041c359724a4ffdbe790e10fec6de810f8aee1c526dd15d383b9a113b6b6a00706052b81040021a13c033a0004e3dd49e00dce869da09afc266d707a3e59377d28aded9d8f264ab890790aa92735f7ca9df8507a1d0823092e29ab7d74336dd9938521c479")
	aa2pubkey,_ := hex.DecodeString("04a579d5b6764ae6b0f2f302d8717b4c8bb866ab6915be971798447d1918ec0ea739950ab7784a389e60a0f077ade960ae0e53d353f247c581")
	aa2prikey,_ := hex.DecodeString("3068020101041cd599de5c3c9d65903a263ee6d0d6df6fe7acb6616642aed261a4bac2a00706052b81040021a13c033a0004a579d5b6764ae6b0f2f302d8717b4c8bb866ab6915be971798447d1918ec0ea739950ab7784a389e60a0f077ade960ae0e53d353f247c581")
	aa3pubkey,_ := hex.DecodeString("04142100e66804198329ee8ac6e389f4d4448523f3cb13135f4fd4e0bfa816b00b7f5e53ae6d16c9a23dd8c0a7913934a6d19013a641a8cc8d")
	aa3prikey,_ := hex.DecodeString("3068020101041c56a7bc34c3901d167d8d563c922b3ac0f2c2d3e0a6fa05e51bef493fa00706052b81040021a13c033a0004142100e66804198329ee8ac6e389f4d4448523f3cb13135f4fd4e0bfa816b00b7f5e53ae6d16c9a23dd8c0a7913934a6d19013a641a8cc8d")

	SdkInitialize()
	SdkInstallAndInstantiateCC_SYS(2, 3)
	SdkInstallAndInstantiateCC_STR()
	SdkInstallAndInstantiateCC_AA("github.com/thorweiyan/ABEPasswordPlatform/chaincodeImpl/chaincode/AA_1/","AA_1cc",aa1pubkey,aa1pubkey,aa1prikey)
	SdkInstallAndInstantiateCC_AA("github.com/thorweiyan/ABEPasswordPlatform/chaincodeImpl/chaincode/AA_2/","AA_2cc",aa2pubkey,aa2pubkey,aa2prikey)
	SdkInstallAndInstantiateCC_AA("github.com/thorweiyan/ABEPasswordPlatform/chaincodeImpl/chaincode/AA_3/","AA_3cc",aa3pubkey,aa3pubkey,aa3prikey)
	SdkSYSputaaList()
	SdkUpdateAAList("AA_1cc",aa1prikey)
	SdkUpdateAAList("AA_2cc",aa2prikey)
	SdkUpdateAAList("AA_2cc",aa2prikey)

	SdkStartABE1("AA_1cc",aa1prikey)
	SdkStartABE1("AA_2cc",aa2prikey)
	SdkStartABE1("AA_3cc",aa3prikey)
	SdkStartABE2("AA_1cc",aa1prikey)
	SdkStartABE2("AA_2cc",aa2prikey)
	SdkStartABE2("AA_3cc",aa3prikey)
	SdkStartABE3("AA_1cc",aa1prikey)
	SdkStartABE3("AA_2cc",aa2prikey)
	SdkStartABE3("AA_3cc",aa3prikey)
}


//Just a example, need environment
func SdkInitialize() {
	// Initialization of the Fabric SDK from the previously set properties
	err := fSetup.Initialize()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
		return
	}
	// Close SDK
	defer fSetup.CloseSDK()
}

func SdkInstallAndInstantiateCC_SYS(t,n int) {
	fSetup.ChaincodePath = "github.com/thorweiyan/ABEPasswordPlatform/chaincodeImpl/chaincode/System/"
	fSetup.ChainCodeID = "SYScc"
	// Install and instantiate the chaincode
	err := fSetup.InstallAndInstantiateCC([]string{"init",strconv.Itoa(t), strconv.Itoa(n)})
	if err != nil {
		fmt.Printf("Unable to install and instantiate the chaincode: %v\n", err)
		return
	}
}

func SdkInstallAndInstantiateCC_STR() {
	fSetup.ChaincodePath = "github.com/thorweiyan/ABEPasswordPlatform/chaincodeImpl/chaincode/Storage/"
	fSetup.ChainCodeID = "STRcc"
	// Install and instantiate the chaincode
	err := fSetup.InstallAndInstantiateCC([]string{"init"})
	if err != nil {
		fmt.Printf("Unable to install and instantiate the chaincode: %v\n", err)
		return
	}
}

func SdkInstallAndInstantiateCC_AA(ccPath string, ccId string, ownerPubkey []byte, aaPubkey []byte, aaPrikey []byte) {
	fSetup.ChaincodePath = ccPath
	fSetup.ChainCodeID = ccId

	// Install and instantiate the chaincode
	err := fSetup.InstallAndInstantiateCC([]string{"init",string(ownerPubkey), "AA_1",string(aaPrikey),string(aaPubkey)})
	if err != nil {
		fmt.Printf("Unable to instantiate the chaincode: %v\n", err)
		return
	}
}

func signTransaction(priKey []byte, args []string) ([]string, error){
	sigMsg := ""
	for _,v := range args {
		sigMsg += v
	}

	r, s, err := wrapper.EcdsaSignNormal(priKey, sigMsg)
	if err != nil {
		return []string{}, fmt.Errorf(err.Error())
	}
	re := []string{string(r), string(s)}
	re = append(re, args...)
	return re, nil
}

func SdkStartABE1(ccId string, OwnerPriKey []byte) {
	fSetup.ChainCodeID = ccId
	params,_ := signTransaction(OwnerPriKey, []string{string("startABE1")})
	params = append([]string{"startABE1"}, params...)
	trcid, err := fSetup.Invoke(params)
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	fmt.Println(trcid)
	fmt.Printf("%x\n",trcid)
}

func SdkStartABE2(ccId string, OwnerPriKey []byte) {
	fSetup.ChainCodeID = ccId
	params,_ := signTransaction(OwnerPriKey, []string{string("startABE2")})
	params = append([]string{"startABE2"}, params...)
	trcid, err := fSetup.Invoke(params)
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	fmt.Println(trcid)
	fmt.Printf("%x\n",trcid)
}

func SdkStartABE3(ccId string, OwnerPriKey []byte) {
	fSetup.ChainCodeID = ccId
	params,_ := signTransaction(OwnerPriKey, []string{string("startABE3")})
	params = append([]string{"startABE3"}, params...)
	trcid, err := fSetup.Invoke(params)
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	fmt.Println(trcid)
	fmt.Printf("%x\n",trcid)
}

func SdkUpdateAAList(ccId string, OwnerPriKey []byte) {
	fSetup.ChainCodeID = ccId
	params,_ := signTransaction(OwnerPriKey, []string{string("updateAAList")})
	params = append([]string{"updateAAList"}, params...)
	trcid, err := fSetup.Invoke(params)
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	fmt.Println(trcid)
	fmt.Printf("%x\n",trcid)
}

func SdkAAccGetPubKey(ccId string) string{
	fSetup.ChainCodeID = ccId
	trcid, err := fSetup.Invoke([]string{"getPubKey"})
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	
	fmt.Printf("%x\n",trcid)
	return trcid
}

func SdkSYSputaaList() {
	fSetup.ChainCodeID = "SYScc"
	trcid, err := fSetup.Invoke([]string{"putaaList"})
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	fmt.Println(trcid)
	fmt.Printf("%x\n",trcid)
}


func SdkUserMethods(ccId string, OwnerPriKey []byte, data wrapper.UserData, method string) {
	fSetup.ChainCodeID = ccId
	pass,err := data.Serialize()
	fmt.Println("serialerr:",err)
	params,_ := signTransaction(OwnerPriKey, []string{string(pass)})
	params = append([]string{"userMethod",method+"Special"}, params...)
	trcid, err := fSetup.Invoke(params)
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	//fmt.Println(trcid)
	fmt.Printf("%v\n",trcid)
	fmt.Printf("%x\n",trcid)
}

func SdkUserMethodn(ccId string, OwnerPriKey []byte, data wrapper.UserData, method string) {
	fSetup.ChainCodeID = ccId

	pass,err := data.Serialize()
	fmt.Println("serialerr:",err)
	params,_ := signTransaction(OwnerPriKey, []string{string(pass)})
	params = append([]string{"userMethod",method}, params...)
	trcid, err := fSetup.Invoke(params)
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	//fmt.Println(trcid)
	fmt.Printf("%v\n",trcid)
	fmt.Printf("%x\n",trcid)
}


func SdkThirdParty(ccId string, OwnerPriKey []byte, userName string, passwordHash []byte) {
	fSetup.ChainCodeID = ccId
	passParams, err := signTransaction(OwnerPriKey,[]string{userName, string(passwordHash)})
	if err != nil {
		fmt.Errorf("sendToAA AASecret" + err.Error())
	}
	passParams = append([]string{"thirdVerify"}, passParams...)
	trcid, err := fSetup.Invoke(passParams)
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	fmt.Println(trcid)
	fmt.Printf("%x\n",trcid)
}


