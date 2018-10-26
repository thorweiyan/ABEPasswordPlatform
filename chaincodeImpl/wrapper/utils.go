package wrapper

import (
	"fmt"
	"strings"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//***************************  Chaincode Utils  ************************
func SplitStringbyn(a string) []string {
	return strings.SplitN(a, "\n\n", -1)
}

//args: r s args...
func IsOwner(stub shim.ChaincodeStubInterface, args []string) (bool,error) {
	//获取公钥
	pubKey, err := stub.GetState("OwnerPubKey")
	if err != nil {
		return false, fmt.Errorf("Don't have PubKey: " + err.Error())
	}

	sigMsg := ""
	for _,v := range args[2:] {
		sigMsg += v
	}

	isRight, err := EcdsaVerifyNormal(pubKey, sigMsg, []byte(args[0]), []byte(args[1]))

	if isRight {
		return true, nil
	}else {
		return false, fmt.Errorf("Owner Verify Error: " + err.Error())
	}
}

//args: AA_ID(AA_1) r s 参数...
func IsAA(stub shim.ChaincodeStubInterface, args []string) (bool,error) {
	//获取公钥
	pubKey, err := stub.GetState(args[0])
	if err != nil {
		return false, fmt.Errorf("Don't have this AA ID: " + err.Error())
	}

	sigMsg := ""
	for _,v := range args[3:] {
		sigMsg += v
	}

	isRight, err := EcdsaVerifyNormal(pubKey, sigMsg, []byte(args[1]), []byte(args[2]))

	if isRight {
		return true, nil
	}else {
		return false, fmt.Errorf("AA Verify Error: " + err.Error())
	}
}

//args: r s 参数...
func IsSYS(stub shim.ChaincodeStubInterface, args []string) (bool,error) {
	//获取公钥
	pubKey, err := stub.GetState("SYSChaincode")
	if err != nil {
		return false, fmt.Errorf("Can't get SYSChaincode's Pubkey: " + err.Error())
	}

	sigMsg := ""
	for _,v := range args[2:] {
		sigMsg += v
	}

	isRight, err := EcdsaVerifyNormal(pubKey, sigMsg, []byte(args[0]), []byte(args[1]))

	if isRight {
		return true, nil
	}else {
		return false, fmt.Errorf("SYS Verify Error: " + err.Error())
	}
}

func SignTransaction(stub shim.ChaincodeStubInterface, args []string) ([]string, error){
	//获取公钥
	priKey, err := stub.GetState("PriKey")
	if err != nil {
		return []string{}, fmt.Errorf("Don't have PriKey: " + err.Error())
	}

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

//{"Args":["chaincode","method"...]}'
func Call(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("call: ", args)
	sub_args := args[1:]
	return stub.InvokeChaincode(args[0], ToChaincodergs(sub_args...), stub.GetChannelID())
}

func ToChaincodergs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}
