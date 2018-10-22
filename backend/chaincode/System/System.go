package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	pb "github.com/hyperledger/fabric/protos/peer"

	"github.com/thorweiyan/ABEPasswordPlatform/backend/wrapper"
	"strconv"
)

type Chaincode struct {
	Length int
	T int //阈值
	N int //AA总数
	Initialized bool //是否初始化完成
	AAList []string
}



//***************************  ABE Init  ***************************
func (t *Chaincode) abeInit(stub shim.ChaincodeStubInterface, args []string) pb.Response {

}


//***************************  Communicate with STR  ***************************
func (t *Chaincode) aaListToSTR(stub shim.ChaincodeStubInterface) error {
	passParams, err := wrapper.SignTransaction(stub, t.AAList)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	passParams = append([]string{"STRcc", "put", "AAList"}, passParams...)
	response := wrapper.Call(stub, passParams)
	if response.Status == 200 {
		return nil
	}else {
		return fmt.Errorf(response.Message)
	}
}

func (t *Chaincode) startAAsABECommunicate(stub shim.ChaincodeStubInterface, args []string) pb.Response {

}

//***************************  Communicate with AA  ***************************
func (t *Chaincode) sendParamsToAA(stub shim.ChaincodeStubInterface) error {

}

//args: AA_ID(AA_1、AA_2...)
func (t *Chaincode) receiveRegisterFromAA(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if t.Initialized {
		return shim.Error("System is already up")
	}
	response := wrapper.Call(stub, []string{args[0]+"cc", "getPubKey"})
	if response.Status == 200 {
		t.Length += 1
		t.AAList = append(t.AAList, string(response.Payload))
		err := t.sendParamsToAA(stub)
		if err != nil {
			return shim.Error(err.Error())
		}

		if t.Length == t.N {
			// aa all online
			err := t.aaListToSTR(stub)
			if err != nil {
				return shim.Error(err.Error())
			}


		}
	}else if response.Status == 200 {
		return response
	}

}








//***************************  Chaincode interface  ***************************
func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("System Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "register" {
		return t.receiveRegisterFromAA(stub, args)
	}else if function == "getPubKey" {
		PubKey, err := stub.GetState("PubKey")
		if err!= nil {
			return shim.Error("GetState Pubkey error\n")
		}
		return shim.Success(PubKey)
	}
	return shim.Error("Invalid invoke function name. Expecting \"register\" ")
}

//args: t,n
func (t *Chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("System Init")
	_, args := stub.GetFunctionAndParameters()

	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Expecting 0")
	}

	var err error
	t.Length = 0
	t.T, err = strconv.Atoi(args[0])
	if err!= nil {
		return shim.Error("t is not a string(int)\n")
	}
	t.N, err = strconv.Atoi(args[1])
	if err!= nil {
		return shim.Error("n is not a string(int)\n")
	}
	t.Initialized = false

	CCPrikey, CCPubkey, err := wrapper.EcdsaSetUpNormal()
	if err!= nil {
		return shim.Error("t is not a string(int)\n")
	}

	err = stub.PutState("PriKey", CCPrikey)
	if err!= nil {
		return shim.Error("PutState Prikey error\n")
	}
	err = stub.PutState("PubKey", CCPubkey)
	if err!= nil {
		return shim.Error("PutState Pubkey error\n")
	}

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("Error starting SYSChaincode: %s", err)
	}
}