package main

import (
	"encoding/base64"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"fmt"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Chaincode struct {
}

func toChaincodergs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}

func (t *Chaincode) isAA(stub shim.ChaincodeStubInterface) (bool,error) {
	id, err := cid.GetID(stub)
	if err != nil {
		return false, fmt.Errorf("getid error: " + err.Error())
	}

	idReadable, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return false, fmt.Errorf("base64 decode error: " + err.Error())
	}
	strId := string(idReadable)

	boolValue, err := stub.GetState(strId)
	if err != nil {
		return false, fmt.Errorf("Failed to get AAChaincode state\n")
	}
	if string(boolValue) != "True" {
		return false, fmt.Errorf("It's not from AA\n")
	}

	return true, nil
}

//{"Args":["call","chaincode","method"...]}'
func (t *Chaincode) call(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("call: ", args)
	sub_args := args[1:]
	return stub.InvokeChaincode(args[0], toChaincodergs(sub_args...), stub.GetChannelID())
}

func (t *Chaincode) sendToAA(stub shim.ChaincodeStubInterface, args []string, AAList []string) pb.Response {

}

func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("System Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "put" {
		return t.put(stub, args)
	}else if function == "get" {
		return t.get(stub, args)
	}
	
	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}

func (t *Chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("System Init")
	_, args := stub.GetFunctionAndParameters()

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("Error starting SYSChaincode: %s", err)
	}
}