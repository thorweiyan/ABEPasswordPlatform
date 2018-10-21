package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"crypto/ecdsa"
	"crypto/elliptic"
	"strings"
	"math/big"
)

type Chaincode struct {
}

func splitStringbyn(a string) []string {
	return strings.SplitN(a, "\n\n", -1)
}

//args: str(AA的序号) r s 参数...
func (t *Chaincode) isAA(stub shim.ChaincodeStubInterface, args []string) (bool,error) {
	//获取公钥
	pubKeyByte, err := stub.GetState("AA_" + args[0])
	if err != nil {
		return false, fmt.Errorf("Don't have this AA ID: " + err.Error())
	}
	pubKeyString := splitStringbyn(string(pubKeyByte))
	x := big.NewInt(0)
	y := big.NewInt(0)
	r := big.NewInt(0)
	s := big.NewInt(0)
	x.SetBytes([]byte(pubKeyString[0]))
	y.SetBytes([]byte(pubKeyString[1]))
	r.SetBytes([]byte(args[1]))
	s.SetBytes([]byte(args[2]))

	pubKey := ecdsa.PublicKey{elliptic.P224(), x, y}

	sigMsg := ""
	for _,v := range args[3:] {
		sigMsg += v
	}

	isRight := ecdsa.Verify(&pubKey, []byte(sigMsg), r, s)

	if isRight {
		return true, nil
	}else {
		return false, fmt.Errorf("AA Verify Error!\n")
	}
}

//args: r s 参数...
func (t *Chaincode) isSYS(stub shim.ChaincodeStubInterface, args []string) (bool,error) {
	//获取公钥
	pubKeyByte, err := stub.GetState("SysChaincode")
	if err != nil {
		return false, fmt.Errorf("Can't get sysChaincode's Pubkey: " + err.Error())
	}
	pubKeyString := splitStringbyn(string(pubKeyByte))
	x := big.NewInt(0)
	y := big.NewInt(0)
	r := big.NewInt(0)
	s := big.NewInt(0)
	x.SetBytes([]byte(pubKeyString[0]))
	y.SetBytes([]byte(pubKeyString[1]))
	r.SetBytes([]byte(args[0]))
	s.SetBytes([]byte(args[1]))

	pubKey := ecdsa.PublicKey{elliptic.P224(), x, y}

	sigMsg := ""
	for _,v := range args[2:] {
		sigMsg += v
	}

	isRight := ecdsa.Verify(&pubKey, []byte(sigMsg), r, s)

	if isRight {
		return true, nil
	}else {
		return false, fmt.Errorf("SYS Verify Error!\n")
	}
}

func (t *Chaincode) put(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	switch args[0] {
	//存放AA的列表
	case "AAList":
		return t.putAAList(stub, args[1:])
		//存放用户账户密码
	case "UserData":
		if len(args) != 3 {
			return shim.Error("Incorrect number of arguments. Expecting 3")
		}
		return t.putUserData(stub, args[1:])
		//存放改密信息
	case "ChangePasswordData":
		if len(args) != 3 {
			return shim.Error("Incorrect number of arguments. Expecting 3")
		}
		return t.putChangePasswordData(stub, args[1:])
		//存放用户自设提示信息
	case "UserTip":
		if len(args) != 3 {
			return shim.Error("Incorrect number of arguments. Expecting 3")
		}
		return t.putUserTip(stub, args[1:])
	default:
		return shim.Error("Can't match any one")
	}

}

//AA初始化过没
func (t *Chaincode) isAAInitialized(stub shim.ChaincodeStubInterface) (bool, error) {
	isInit, err := stub.GetState("SYS_isInitAA")
	if err != nil {
		return false, fmt.Errorf("Failed to get isInitAA state: " + err.Error())
	}
	if string(isInit) != "True" {
		return false, fmt.Errorf("AA is Initialized " + err.Error())
	}
	return true, nil
}

func (t *Chaincode) putAAList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//AA初始化完了没
	isInit, err := t.isAAInitialized(stub)
	if isInit{
		return shim.Error(err.Error())
	}
	//是不是SYS调用的
	rightCreator, err := t.isSYS(stub, args)
	if (err != nil) || !rightCreator {
		return shim.Error("Putting AAList: " + err.Error())
	}
	//去掉rs
	args = args[2:]
	//执行存储命令
	for i := range args {
		err = stub.PutState("AA_"+strconv.Itoa(i), []byte(args[i]))
		if err != nil {
			return shim.Error("Putting AAList: " + err.Error())
		}
	}
	//设置AA已经初始化
	err = stub.PutState("SYS_isInitAA", []byte("True"))
	if err != nil {
		return shim.Error("Putting AAList: " + err.Error())
	}
	err = stub.PutState("AAListLength", []byte(strconv.Itoa(len(args))))
	if err != nil {
		return shim.Error("Putting AAList: " + err.Error())
	}


	return shim.Success(nil)
}

func (t *Chaincode) putUserData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//AA初始化
	isInit, err := t.isAAInitialized(stub)
	if !isInit || err != nil {
		return shim.Error(err.Error())
	}
	//是不是AA调用的
	rightCreator, err := t.isAA(stub, args)
	if (err != nil) || !rightCreator {
		return shim.Error("Putting UserData: " + err.Error())
	}
	//去掉rs
	args = args[3:]
	//直接覆盖用户数据
	err = stub.PutState("UserData_"+args[1], []byte(args[2]))
	if err != nil {
		return shim.Error("Putting UserData: " + err.Error())
	}
	return shim.Success(nil)
}

func (t *Chaincode) putChangePasswordData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//AA初始化
	isInit, err := t.isAAInitialized(stub)
	if !isInit || err != nil {
		return shim.Error(err.Error())
	}
	//是不是AA调用的
	rightCreator, err := t.isAA(stub, args)
	if (err != nil) || !rightCreator {
		return shim.Error("Putting ChangePasswordData: " + err.Error())
	}
	//去掉rs
	args = args[3:]
	//直接覆盖用户数据
	err = stub.PutState("ChangePasswordData_"+args[1], []byte(args[2]))
	if err != nil {
		return shim.Error("Putting ChangePasswordData: " + err.Error())
	}
	return shim.Success(nil)
}

func (t *Chaincode) putUserTip(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//AA初始化
	isInit, err := t.isAAInitialized(stub)
	if !isInit || err != nil {
		return shim.Error(err.Error())
	}
	//是不是AA调用的
	rightCreator, err := t.isAA(stub, args)
	if (err != nil) || !rightCreator {
		return shim.Error("Putting UserTip: " + err.Error())
	}
	//去掉rs
	args = args[3:]
	//直接覆盖用户数据
	err = stub.PutState("UserTip_"+args[1], []byte(args[2]))
	if err != nil {
		return shim.Error("Putting UserTip: " + err.Error())
	}
	return shim.Success(nil)
}

//args: func aa_id r s ...
func (t *Chaincode) get(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}
	//是不是AA调用的
	rightCreator, err := t.isAA(stub, args[1:])
	if (err != nil) || !rightCreator {
		return shim.Error("When Get: " + err.Error())
	}

	//AA初始化完了没
	isInit, err := t.isAAInitialized(stub)
	if !isInit || err != nil {
		return shim.Error(err.Error())
	}
	switch args[0] {
	case "AAList":
		rs, err := t.getAAList(stub, args[4])
		if err != nil {
			return shim.Error(err.Error())
		}else {
			return shim.Success(rs)
		}
	case "UserData":
		result, err := stub.GetState("UserData_" + args[4])
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(result)
	case "ChangePasswordData":
		result, err := stub.GetState("ChangePasswordData_" + args[4])
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(result)
	case "UserTip":
		result, err := stub.GetState("UserTip_" + args[4])
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(result)
	default:
		return shim.Error("Can't match any one")
	}
}

func (t *Chaincode) getAAList(stub shim.ChaincodeStubInterface, aaid string) ([]byte, error) {
	aaListLength, err := stub.GetState("AAListLength")
	if err != nil {
		return []byte{}, fmt.Errorf(err.Error())
	}
	var rs []byte
	length, err := strconv.Atoi(string(aaListLength))
	if err != nil {
		return []byte{}, fmt.Errorf(err.Error())
	}
	for i := 0;i< length;i++{
		tempre, err := stub.GetState("AA_" + strconv.Itoa(i))
		if err != nil {
			return []byte{}, fmt.Errorf(err.Error())
		}
		if string(tempre) == aaid {
			continue
		}else {
			rs = append(rs, tempre...)
			rs = append(rs, []byte("\n\n")...)
		}
	}
	return rs, nil
}

func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Storage Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "put" {
		return t.put(stub, args)
	}else if function == "get" {
		return t.get(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"get\" \"put\"")
}

func (t *Chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Storage Init")
	_, args := stub.GetFunctionAndParameters()

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("SysChaincode", []byte(args[0]))
	if err != nil {
		return shim.Error("Putting SysChaincode: " + err.Error())
	}

	err = stub.PutState("SYS_isInitAA", []byte("False"))
	if err != nil {
		return shim.Error("Putting isInitAA: " + err.Error())
	}

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("Error starting STRChaincode: %s", err)
	}
}