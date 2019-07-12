package wrapper

import (
"testing"
"os"
"fmt"
"github.com/thorweiyan/fabric_go_sdk"
	"encoding/hex"
	"strings"
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
	ChaincodeVersion: "1",
	OrgAdmin:        "Admin",
	OrgName:         "org1",
	ConfigFile:      os.Getenv("GOPATH") + "/src/github.com/thorweiyan/fabric_go_sdk/config.yaml",

	// User parameters
	UserName: "User1",
}
func TestAll(t *testing.T) {
	//初始化channel
	TestInitialize(t)
	//安装并实例化合约
	TestInstallAndInstantiateCC_SYS(t)
	TestInstallAndInstantiateCC_STR(t)
	TestInstallAndInstantiateCC_AA1(t)
	TestInstallAndInstantiateCC_AA2(t)
	TestInstallAndInstantiateCC_AA3(t)
	//一系列配置操作
	TestSYSputaaList(t)
	TestUpdate1(t)
	TestUpdate2(t)
	TestUpdate3(t)
	TestStartABE11(t)
	TestStartABE12(t)
	TestStartABE13(t)
	TestStartABE21(t)
	TestStartABE22(t)
	TestStartABE23(t)
	TestStartABE31(t)
	TestStartABE32(t)
	TestStartABE33(t)

	//TestUserMethodsignups1(t)
	//TestUserMethodsignupn2(t)
	//TestUserMethodsignups1(t)
	//TestUserMethodsignupn2(t)
	//
	//
	//TestUserMethodChanges1(t)
	//TestUserMethodChangen2(t)
	//
	//TestUserMethodGetTips1(t)
	//TestUserMethodGetTipn2(t)
	//
	//TestThirdParty(t)
}

func Test2(t *testing.T) {
	TestUserMethodChanges1(t)
	TestUserMethodChangen2(t)

	TestUserMethodGetTips1(t)
	TestUserMethodGetTipn2(t)

	TestThirdParty(t)
}

//Just a example, need environment
func TestInitialize(t *testing.T) {
	// Initialization of the Fabric SDK from the previously set properties
	err := fSetup.Initialize()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
		return
	}
	// Close SDK
	defer fSetup.CloseSDK()
}

func TestInstallAndInstantiateCC_SYS(t *testing.T) {
	fSetup.ChaincodePath = "github.com/thorweiyan/ABEPasswordPlatform/chaincodeImpl/chaincode/System/"
	fSetup.ChainCodeID = "SYScc"
	fSetup.ChaincodeVersion = "2"
	// Install and instantiate the chaincode
	err := fSetup.InstallAndInstantiateCC([]string{"init","2", "3"})
	if err != nil {
		fmt.Printf("Unable to install and instantiate the chaincode: %v\n", err)
		return
	}
}

func TestInstallAndInstantiateCC_STR(t *testing.T) {
	fSetup.ChaincodePath = "github.com/thorweiyan/ABEPasswordPlatform/chaincodeImpl/chaincode/Storage/"
	fSetup.ChainCodeID = "STRcc"
	fSetup.ChaincodeVersion = "2"
	// Install and instantiate the chaincode
	err := fSetup.InstallAndInstantiateCC([]string{"init"})
	if err != nil {
		fmt.Printf("Unable to install and instantiate the chaincode: %v\n", err)
		return
	}
}

func TestInstallAndInstantiateCC_AA1(t *testing.T) {
	fSetup.ChaincodePath = "github.com/thorweiyan/ABEPasswordPlatform/chaincodeImpl/chaincode/AA_1/"
	fSetup.ChainCodeID = "AA_1cc"
	fSetup.ChaincodeVersion = "1"

	// Install and instantiate the chaincode
	aa1pubkey,_ := hex.DecodeString("04e3dd49e00dce869da09afc266d707a3e59377d28aded9d8f264ab890790aa92735f7ca9df8507a1d0823092e29ab7d74336dd9938521c479")
	aa1prikey,_ := hex.DecodeString("3068020101041c359724a4ffdbe790e10fec6de810f8aee1c526dd15d383b9a113b6b6a00706052b81040021a13c033a0004e3dd49e00dce869da09afc266d707a3e59377d28aded9d8f264ab890790aa92735f7ca9df8507a1d0823092e29ab7d74336dd9938521c479")
	err := fSetup.InstallAndInstantiateCC([]string{"init",string(aa1pubkey), "AA_1",string(aa1prikey),string(aa1pubkey)})
	if err != nil {
		fmt.Printf("Unable to instantiate the chaincode: %v\n", err)
		return
	}
}

func TestInstallAndInstantiateCC_AA2(t *testing.T) {
	fSetup.ChaincodePath = "github.com/thorweiyan/ABEPasswordPlatform/chaincodeImpl/chaincode/AA_2/"
	fSetup.ChainCodeID = "AA_2cc"
	fSetup.ChaincodeVersion = "1"

	// Install and instantiate the chaincode
	aa1pubkey,_ := hex.DecodeString("04a579d5b6764ae6b0f2f302d8717b4c8bb866ab6915be971798447d1918ec0ea739950ab7784a389e60a0f077ade960ae0e53d353f247c581")
	aa1prikey,_ := hex.DecodeString("3068020101041cd599de5c3c9d65903a263ee6d0d6df6fe7acb6616642aed261a4bac2a00706052b81040021a13c033a0004a579d5b6764ae6b0f2f302d8717b4c8bb866ab6915be971798447d1918ec0ea739950ab7784a389e60a0f077ade960ae0e53d353f247c581")
	err := fSetup.InstallAndInstantiateCC([]string{"init",string(aa1pubkey), "AA_2",string(aa1prikey),string(aa1pubkey)})
	if err != nil {
		fmt.Printf("Unable to instantiate the chaincode: %v\n", err)
		return
	}
}

func TestInstallAndInstantiateCC_AA3(t *testing.T) {
	fSetup.ChaincodePath = "github.com/thorweiyan/ABEPasswordPlatform/chaincodeImpl/chaincode/AA_3/"
	fSetup.ChainCodeID = "AA_3cc"
	fSetup.ChaincodeVersion = "1"

	// Install and instantiate the chaincode
	aa1pubkey,_ := hex.DecodeString("04142100e66804198329ee8ac6e389f4d4448523f3cb13135f4fd4e0bfa816b00b7f5e53ae6d16c9a23dd8c0a7913934a6d19013a641a8cc8d")
	aa1prikey,_ := hex.DecodeString("3068020101041c56a7bc34c3901d167d8d563c922b3ac0f2c2d3e0a6fa05e51bef493fa00706052b81040021a13c033a0004142100e66804198329ee8ac6e389f4d4448523f3cb13135f4fd4e0bfa816b00b7f5e53ae6d16c9a23dd8c0a7913934a6d19013a641a8cc8d")
	err := fSetup.InstallAndInstantiateCC([]string{"init",string(aa1pubkey), "AA_3",string(aa1prikey),string(aa1pubkey)})
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

	r, s, err := EcdsaSignNormal(priKey, sigMsg)
	if err != nil {
		return []string{}, fmt.Errorf(err.Error())
	}
	re := []string{string(r), string(s)}
	re = append(re, args...)
	return re, nil
}
func TestStartABE11(t *testing.T) {
	fSetup.ChaincodePath += "AA_1/"
	fSetup.ChainCodeID = "AA_1cc"
	fSetup.ChaincodeVersion = "1"
	aa1prikey,_ := hex.DecodeString("3068020101041c359724a4ffdbe790e10fec6de810f8aee1c526dd15d383b9a113b6b6a00706052b81040021a13c033a0004e3dd49e00dce869da09afc266d707a3e59377d28aded9d8f264ab890790aa92735f7ca9df8507a1d0823092e29ab7d74336dd9938521c479")
	params,_ := signTransaction(aa1prikey, []string{string("startABE1")})
	params = append([]string{"startABE1"}, params...)
	trcid, err := fSetup.Invoke(params)
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	fmt.Println(trcid)
	fmt.Printf("%x\n",trcid)
}
func TestStartABE12(t *testing.T) {
	fSetup.ChaincodePath += "AA_1/"
	fSetup.ChainCodeID = "AA_2cc"
	fSetup.ChaincodeVersion = "1"
	aa1prikey,_ := hex.DecodeString("3068020101041cd599de5c3c9d65903a263ee6d0d6df6fe7acb6616642aed261a4bac2a00706052b81040021a13c033a0004a579d5b6764ae6b0f2f302d8717b4c8bb866ab6915be971798447d1918ec0ea739950ab7784a389e60a0f077ade960ae0e53d353f247c581")
	params,_ := signTransaction(aa1prikey, []string{string("startABE1")})
	params = append([]string{"startABE1"}, params...)
	trcid, err := fSetup.Invoke(params)
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	fmt.Println(trcid)
	fmt.Printf("%x\n",trcid)
}
func TestStartABE13(t *testing.T) {
	fSetup.ChaincodePath += "AA_3/"
	fSetup.ChainCodeID = "AA_3cc"
	fSetup.ChaincodeVersion = "1"
	aa1prikey,_ := hex.DecodeString("3068020101041c56a7bc34c3901d167d8d563c922b3ac0f2c2d3e0a6fa05e51bef493fa00706052b81040021a13c033a0004142100e66804198329ee8ac6e389f4d4448523f3cb13135f4fd4e0bfa816b00b7f5e53ae6d16c9a23dd8c0a7913934a6d19013a641a8cc8d")
	params,_ := signTransaction(aa1prikey, []string{string("startABE1")})
	params = append([]string{"startABE1"}, params...)
	trcid, err := fSetup.Invoke(params)
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	fmt.Println(trcid)
	fmt.Printf("%x\n",trcid)
}
func TestStartABE21(t *testing.T) {
	fSetup.ChaincodePath += "AA_1/"
	fSetup.ChainCodeID = "AA_1cc"
	fSetup.ChaincodeVersion = "1"
	aa1prikey,_ := hex.DecodeString("3068020101041c359724a4ffdbe790e10fec6de810f8aee1c526dd15d383b9a113b6b6a00706052b81040021a13c033a0004e3dd49e00dce869da09afc266d707a3e59377d28aded9d8f264ab890790aa92735f7ca9df8507a1d0823092e29ab7d74336dd9938521c479")
	params,_ := signTransaction(aa1prikey, []string{string("startABE2")})
	params = append([]string{"startABE2"}, params...)
	trcid, err := fSetup.Invoke(params)
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	fmt.Println(trcid)
	fmt.Printf("%x\n",trcid)
}
func TestStartABE22(t *testing.T) {
	fSetup.ChaincodePath += "AA_1/"
	fSetup.ChainCodeID = "AA_2cc"
	fSetup.ChaincodeVersion = "1"
	aa1prikey,_ := hex.DecodeString("3068020101041cd599de5c3c9d65903a263ee6d0d6df6fe7acb6616642aed261a4bac2a00706052b81040021a13c033a0004a579d5b6764ae6b0f2f302d8717b4c8bb866ab6915be971798447d1918ec0ea739950ab7784a389e60a0f077ade960ae0e53d353f247c581")
	params,_ := signTransaction(aa1prikey, []string{string("startABE2")})
	params = append([]string{"startABE2"}, params...)
	trcid, err := fSetup.Invoke(params)
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	fmt.Println(trcid)
	fmt.Printf("%x\n",trcid)
}
func TestStartABE23(t *testing.T) {
	fSetup.ChaincodePath += "AA_3/"
	fSetup.ChainCodeID = "AA_3cc"
	fSetup.ChaincodeVersion = "1"
	aa1prikey,_ := hex.DecodeString("3068020101041c56a7bc34c3901d167d8d563c922b3ac0f2c2d3e0a6fa05e51bef493fa00706052b81040021a13c033a0004142100e66804198329ee8ac6e389f4d4448523f3cb13135f4fd4e0bfa816b00b7f5e53ae6d16c9a23dd8c0a7913934a6d19013a641a8cc8d")
	params,_ := signTransaction(aa1prikey, []string{string("startABE2")})
	params = append([]string{"startABE2"}, params...)
	trcid, err := fSetup.Invoke(params)
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	fmt.Println(trcid)
	fmt.Printf("%x\n",trcid)
}
func TestStartABE31(t *testing.T) {
	fSetup.ChaincodePath += "AA_1/"
	fSetup.ChainCodeID = "AA_1cc"
	fSetup.ChaincodeVersion = "1"
	aa1prikey,_ := hex.DecodeString("3068020101041c359724a4ffdbe790e10fec6de810f8aee1c526dd15d383b9a113b6b6a00706052b81040021a13c033a0004e3dd49e00dce869da09afc266d707a3e59377d28aded9d8f264ab890790aa92735f7ca9df8507a1d0823092e29ab7d74336dd9938521c479")
	params,_ := signTransaction(aa1prikey, []string{string("startABE3")})
	params = append([]string{"startABE3"}, params...)
	trcid, err := fSetup.Invoke(params)
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	fmt.Println(trcid)
	fmt.Printf("%x\n",trcid)
}
func TestStartABE32(t *testing.T) {
	fSetup.ChaincodePath += "AA_1/"
	fSetup.ChainCodeID = "AA_2cc"
	fSetup.ChaincodeVersion = "1"
	aa1prikey,_ := hex.DecodeString("3068020101041cd599de5c3c9d65903a263ee6d0d6df6fe7acb6616642aed261a4bac2a00706052b81040021a13c033a0004a579d5b6764ae6b0f2f302d8717b4c8bb866ab6915be971798447d1918ec0ea739950ab7784a389e60a0f077ade960ae0e53d353f247c581")
	params,_ := signTransaction(aa1prikey, []string{string("startABE3")})
	params = append([]string{"startABE3"}, params...)
	trcid, err := fSetup.Invoke(params)
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	fmt.Println(trcid)
	fmt.Printf("%x\n",trcid)
}
func TestStartABE33(t *testing.T) {
	fSetup.ChaincodePath += "AA_3/"
	fSetup.ChainCodeID = "AA_3cc"
	fSetup.ChaincodeVersion = "1"
	aa1prikey,_ := hex.DecodeString("3068020101041c56a7bc34c3901d167d8d563c922b3ac0f2c2d3e0a6fa05e51bef493fa00706052b81040021a13c033a0004142100e66804198329ee8ac6e389f4d4448523f3cb13135f4fd4e0bfa816b00b7f5e53ae6d16c9a23dd8c0a7913934a6d19013a641a8cc8d")
	params,_ := signTransaction(aa1prikey, []string{string("startABE3")})
	params = append([]string{"startABE3"}, params...)
	trcid, err := fSetup.Invoke(params)
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	fmt.Println(trcid)
	fmt.Printf("%x\n",trcid)
}
func TestUpdate1(t *testing.T) {
	fSetup.ChaincodePath += "AA_1/"
	fSetup.ChainCodeID = "AA_1cc"
	fSetup.ChaincodeVersion = "1"
	aa1prikey,_ := hex.DecodeString("3068020101041c359724a4ffdbe790e10fec6de810f8aee1c526dd15d383b9a113b6b6a00706052b81040021a13c033a0004e3dd49e00dce869da09afc266d707a3e59377d28aded9d8f264ab890790aa92735f7ca9df8507a1d0823092e29ab7d74336dd9938521c479")
	params,_ := signTransaction(aa1prikey, []string{string("updateAAList")})
	params = append([]string{"updateAAList"}, params...)
	trcid, err := fSetup.Invoke(params)
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	fmt.Println(trcid)
	fmt.Printf("%x\n",trcid)
}
func TestUpdate2(t *testing.T) {
	fSetup.ChaincodePath += "AA_1/"
	fSetup.ChainCodeID = "AA_2cc"
	fSetup.ChaincodeVersion = "1"
	aa1prikey,_ := hex.DecodeString("3068020101041cd599de5c3c9d65903a263ee6d0d6df6fe7acb6616642aed261a4bac2a00706052b81040021a13c033a0004a579d5b6764ae6b0f2f302d8717b4c8bb866ab6915be971798447d1918ec0ea739950ab7784a389e60a0f077ade960ae0e53d353f247c581")
	params,_ := signTransaction(aa1prikey, []string{string("updateAAList")})
	params = append([]string{"updateAAList"}, params...)
	trcid, err := fSetup.Invoke(params)
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	fmt.Println(trcid)
	fmt.Printf("%x\n",trcid)
}
func TestUpdate3(t *testing.T) {
	fSetup.ChaincodePath += "AA_3/"
	fSetup.ChainCodeID = "AA_3cc"
	fSetup.ChaincodeVersion = "1"
	aa1prikey,_ := hex.DecodeString("3068020101041c56a7bc34c3901d167d8d563c922b3ac0f2c2d3e0a6fa05e51bef493fa00706052b81040021a13c033a0004142100e66804198329ee8ac6e389f4d4448523f3cb13135f4fd4e0bfa816b00b7f5e53ae6d16c9a23dd8c0a7913934a6d19013a641a8cc8d")
	params,_ := signTransaction(aa1prikey, []string{string("updateAAList")})
	params = append([]string{"updateAAList"}, params...)
	trcid, err := fSetup.Invoke(params)
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	fmt.Println(trcid)
	fmt.Printf("%x\n",trcid)
}

func TestSignTransaction(t *testing.T) {
	fSetup.ChainCodeID = "AA_2cc"
	fSetup.ChaincodeVersion = "1"
	aa1prikey,_ := hex.DecodeString("3068020101041c56a7bc34c3901d167d8d563c922b3ac0f2c2d3e0a6fa05e51bef493fa00706052b81040021a13c033a0004142100e66804198329ee8ac6e389f4d4448523f3cb13135f4fd4e0bfa816b00b7f5e53ae6d16c9a23dd8c0a7913934a6d19013a641a8cc8d")
	passParams, err := signTransaction(aa1prikey,[]string{"AASecret", "4324234324"})
	if err != nil {
		fmt.Errorf("sendToAA AASecret" + err.Error())
	}
	passParams = append([]string{"handleFromAA", "AA_3"}, passParams...)
	trcid, err := fSetup.Invoke(passParams)
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	fmt.Println(trcid)
	fmt.Printf("%x\n",trcid)
}

func TestAAccGetPubKey(t *testing.T) {
	fSetup.ChaincodePath += "AA_3/"
	fSetup.ChainCodeID = "AA_3cc"
	fSetup.ChaincodeVersion = "1"
	trcid, err := fSetup.Invoke([]string{"getPubKey"})
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	//fmt.Println(trcid)
	fmt.Printf("%v\n",trcid)
	fmt.Printf("%x\n",trcid)
}

func TestSYSputaaList(t *testing.T) {
	fSetup.ChaincodePath += "System/"
	fSetup.ChainCodeID = "SYScc"
	fSetup.ChaincodeVersion = "2"
	trcid, err := fSetup.Invoke([]string{"putaaList"})
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	fmt.Println(trcid)
	fmt.Printf("%x\n",trcid)
}


func TestUserMethodsignups1(t *testing.T) {
	fSetup.ChainCodeID = "AA_1cc"
	fSetup.ChaincodeVersion = "1"
	aa1prikey,_ := hex.DecodeString("3068020101041c359724a4ffdbe790e10fec6de810f8aee1c526dd15d383b9a113b6b6a00706052b81040021a13c033a0004e3dd49e00dce869da09afc266d707a3e59377d28aded9d8f264ab890790aa92735f7ca9df8507a1d0823092e29ab7d74336dd9938521c479")
	data := UserData{UserName:"UN:czn",UserPasswordHash:[]byte("123456712342344"),GetTipMessage:"ti shi message",GetTipPolicy:"(UN:czn AND BD:19970212)"}
	data.UserAttributes= []string{"UN:czn","BD:19970212","SFZ:3xxxxxxxxxxxxxxxxxx1","YX:czn@fudan.edu.cn","ZS:dog","ZS:cat"}
	data.ChangePasswordPolicy = "(UN:czn AND BD:19970212 AND SFZ:3xxxxxxxxxxxxxxxxxx1 AND YX:czn@fudan.edu.cn AND ZS:dog)"
	data.GetTipPolicy = "(UN:czn AND BD:19970212 AND SFZ:3xxxxxxxxxxxxxxxxxx1)"
	data.SpecialAAId = "AA_1"
	pass,err := data.Serialize()
	fmt.Println("serialerr:",err)
	params,_ := signTransaction(aa1prikey, []string{string(pass)})
	params = append([]string{"userMethod","userSignUpSpecial"}, params...)
	trcid, err := fSetup.Invoke(params)
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	//fmt.Println(trcid)
	fmt.Printf("%v\n",trcid)
	fmt.Printf("%x\n",trcid)
}

func TestUserMethodsignupn2(t *testing.T) {
	fSetup.ChainCodeID = "AA_2cc"
	fSetup.ChaincodeVersion = "1"
	aa1prikey,_ := hex.DecodeString("3068020101041cd599de5c3c9d65903a263ee6d0d6df6fe7acb6616642aed261a4bac2a00706052b81040021a13c033a0004a579d5b6764ae6b0f2f302d8717b4c8bb866ab6915be971798447d1918ec0ea739950ab7784a389e60a0f077ade960ae0e53d353f247c581")
	data := UserData{UserName:"UN:czn"}
	data.UserAttributes= []string{"UN:czn","BD:19970212","SFZ:3xxxxxxxxxxxxxxxxxx1","YX:czn@fudan.edu.cn","ZS:dog","ZS:cat"}
	//
	data.SpecialAAId = "AA_1"
	pass,err := data.Serialize()
	fmt.Println("serialerr:",err)
	params,_ := signTransaction(aa1prikey, []string{string(pass)})
	params = append([]string{"userMethod","userSignUp"}, params...)
	trcid, err := fSetup.Invoke(params)
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	//fmt.Println(trcid)
	fmt.Printf("%v\n",trcid)
	fmt.Printf("%x\n",trcid)
}

func TestUserMethodChanges1(t *testing.T) {
	fSetup.ChainCodeID = "AA_2cc"
	fSetup.ChaincodeVersion = "1"
	aa1prikey,_ := hex.DecodeString("3068020101041cd599de5c3c9d65903a263ee6d0d6df6fe7acb6616642aed261a4bac2a00706052b81040021a13c033a0004a579d5b6764ae6b0f2f302d8717b4c8bb866ab6915be971798447d1918ec0ea739950ab7784a389e60a0f077ade960ae0e53d353f247c581")
	data := UserData{UserName:"UN:roy",UserPasswordHash:[]byte("cznzuishuai")}
	data.UserAttributes= strings.Split("UN:roy,SFZ:678987000236787654,SJ:17317301908,ZS:shuai,ZS:hei",",")
	pass,err := data.Serialize()
	fmt.Println("serialerr:",err)
	params,_ := signTransaction(aa1prikey, []string{string(pass)})
	params = append([]string{"userMethod","userChangePasswordSpecial"}, params...)
	trcid, err := fSetup.Invoke(params)
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	//fmt.Println(trcid)
	fmt.Printf("%v\n",trcid)
	fmt.Printf("%x\n",trcid)
}

func TestUserMethodChangen2(t *testing.T) {
	fSetup.ChainCodeID = "AA_3cc"
	fSetup.ChaincodeVersion = "1"
	aa1prikey,_ := hex.DecodeString("3068020101041c56a7bc34c3901d167d8d563c922b3ac0f2c2d3e0a6fa05e51bef493fa00706052b81040021a13c033a0004142100e66804198329ee8ac6e389f4d4448523f3cb13135f4fd4e0bfa816b00b7f5e53ae6d16c9a23dd8c0a7913934a6d19013a641a8cc8d")
	data := UserData{UserName:"UN:roy"}
	data.UserAttributes= strings.Split("UN:roy,SFZ:678987000236787654,SJ:17317301908,ZS:shuai,ZS:hei",",")
	//
	data.SpecialAAId = "AA_2"
	pass,err := data.Serialize()
	fmt.Println("serialerr:",err)
	params,_ := signTransaction(aa1prikey, []string{string(pass)})
	params = append([]string{"userMethod","userChangePassword"}, params...)
	trcid, err := fSetup.Invoke(params)
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	//fmt.Println(trcid)
	fmt.Printf("%v\n",trcid)
	fmt.Printf("%x\n",trcid)
}
func TestUserMethodGetTips1(t *testing.T) {
	fSetup.ChainCodeID = "AA_1cc"
	fSetup.ChaincodeVersion = "1"
	aa1prikey,_ := hex.DecodeString("3068020101041c359724a4ffdbe790e10fec6de810f8aee1c526dd15d383b9a113b6b6a00706052b81040021a13c033a0004e3dd49e00dce869da09afc266d707a3e59377d28aded9d8f264ab890790aa92735f7ca9df8507a1d0823092e29ab7d74336dd9938521c479")
	data := UserData{UserName:"UN:czn"}
	data.UserAttributes= []string{"UN:czn","BD:19970212","SFZ:3xxxxxxxxxxxxxxxxxx1"}
	pass,err := data.Serialize()
	fmt.Println("serialerr:",err)
	params,_ := signTransaction(aa1prikey, []string{string(pass)})
	params = append([]string{"userMethod","userGetTipSpecial"}, params...)
	trcid, err := fSetup.Invoke(params)
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	//fmt.Println(trcid)
	fmt.Printf("%v\n",trcid)
	fmt.Printf("%x\n",trcid)
}

func TestUserMethodGetTipn2(t *testing.T) {
	fSetup.ChainCodeID = "AA_2cc"
	fSetup.ChaincodeVersion = "1"
	aa1prikey,_ := hex.DecodeString("3068020101041cd599de5c3c9d65903a263ee6d0d6df6fe7acb6616642aed261a4bac2a00706052b81040021a13c033a0004a579d5b6764ae6b0f2f302d8717b4c8bb866ab6915be971798447d1918ec0ea739950ab7784a389e60a0f077ade960ae0e53d353f247c581")
	data := UserData{UserName:"UN:czn"}
	data.UserAttributes= []string{"UN:czn","BD:19970212","SFZ:3xxxxxxxxxxxxxxxxxx1"}
	//
	data.SpecialAAId = "AA_1"
	pass,err := data.Serialize()
	fmt.Println("serialerr:",err)
	params,_ := signTransaction(aa1prikey, []string{string(pass)})
	params = append([]string{"userMethod","userGetTip"}, params...)
	trcid, err := fSetup.Invoke(params)
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	//fmt.Println(trcid)
	fmt.Printf("%v\n",trcid)
	fmt.Printf("%x\n",trcid)
}





func TestFabricSetup_Invoke(t *testing.T) {
	fSetup.ChaincodePath += "System/"
	fSetup.ChainCodeID = "SYScc"
	fSetup.ChaincodeVersion = "2"
	trcid, err := fSetup.Invoke([]string{"getPubKey"})
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	fmt.Println(trcid)
	fmt.Printf("%x\n",trcid)
}

func TestFabricSetup_Query(t *testing.T) {
	payload, err := fSetup.Query([]string{"invoke", "query", "hello"})
	if err != nil {
		fmt.Println("query error!", err)
	}
	fmt.Println(payload)
}

func TestThirdParty(t *testing.T) {
	fSetup.ChainCodeID = "AA_2cc"
	fSetup.ChaincodeVersion = "1"
	aa1prikey,_ := hex.DecodeString("3068020101041cd599de5c3c9d65903a263ee6d0d6df6fe7acb6616642aed261a4bac2a00706052b81040021a13c033a0004a579d5b6764ae6b0f2f302d8717b4c8bb866ab6915be971798447d1918ec0ea739950ab7784a389e60a0f077ade960ae0e53d353f247c581")
	passParams, err := signTransaction(aa1prikey,[]string{"UN:czn", string([]byte("cznzuishuai"))})
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

func TestRecover(t *testing.T) {
	fSetup.ChainCodeID = "AA_2cc"
	fSetup.ChaincodeVersion = "1"
	aa1prikey,_ := hex.DecodeString("3068020101041cd599de5c3c9d65903a263ee6d0d6df6fe7acb6616642aed261a4bac2a00706052b81040021a13c033a0004a579d5b6764ae6b0f2f302d8717b4c8bb866ab6915be971798447d1918ec0ea739950ab7784a389e60a0f077ade960ae0e53d353f247c581")
	passParams, err := signTransaction(aa1prikey,[]string{"recoverParams","shabi"})
	if err != nil {
		fmt.Errorf("sendToAA AASecret" + err.Error())
	}
	passParams = append([]string{"recoverParams"}, passParams...)
	trcid, err := fSetup.Invoke(passParams)
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	fmt.Println(trcid)
	fmt.Printf("%x\n",trcid)
}

