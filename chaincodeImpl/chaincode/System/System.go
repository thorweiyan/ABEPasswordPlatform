package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/thorweiyan/MulticenterABEForFabric"
	"github.com/thorweiyan/ABEPasswordPlatform/chaincodeImpl/wrapper"
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
func (t *Chaincode) abeInit(stub shim.ChaincodeStubInterface) error{
	sysAbe := new(MulticenterABEForFabric.MAFFscheme)
	sysAbe.SYSInit(t.T, t.N)
	err := stub.PutState("PubKeyParams", sysAbe.PublicKey.Serialize())
	if err != nil {
		return fmt.Errorf("Put MPK of ABE wrong!\n")
	}
	return nil
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

func (t *Chaincode) startAAsABECommunicate(stub shim.ChaincodeStubInterface) error {
	for i,aa := range t.AAList {
		//sign aa's pubkey
		args, err := wrapper.SignTransaction(stub, []string{aa})
		if err != nil {
			return err
		}
		payload := append([]string{"AA_"+strconv.Itoa(i)+"cc", "startABECommunicate"}, args...)
		tempResponse := wrapper.Call(stub, payload)
		if tempResponse.Status != 200 {
			return fmt.Errorf("StartAAsABECommmunicate error: " + tempResponse.Message)
		}
	}
	return nil
}

//***************************  Communicate with AA  ***************************
func (t *Chaincode) sendParamsToAA(stub shim.ChaincodeStubInterface, aaCCName string) error {
	abeMPK, err := stub.GetState("PubKeyParams")
	if err != nil {
		return fmt.Errorf("Get ABE's MPK error: " + err.Error())
	}

	args, err := wrapper.SignTransaction(stub, []string{string(abeMPK)})
	if err != nil {
		return fmt.Errorf("Sign ABE's MPK error: " + err.Error())
	}

	args = append([]string{aaCCName, "receiveMPK"}, args...)
	response := wrapper.Call(stub, args)
	if response.Status == 200 {
		return nil
	}
	return fmt.Errorf(response.Message)
}

//args: AA_ID(AA_1、AA_2...)
func (t *Chaincode) receiveRegisterFromAA(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if t.Initialized {
		return shim.Error("Already Initialized")
	}
	//获取AA的公钥,并检测其cc是否上线
	response := wrapper.Call(stub, []string{args[0]+"cc", "getPubKey"})
	if response.Status == 200 {
		t.Length += 1
		//去掉AA_
		id,err := strconv.Atoi(args[0][3:])
		if err != nil {
			return shim.Error("AA_ID's type error")
		}
		t.AAList[id] = args[0]

		//send MPK to AA
		err = t.sendParamsToAA(stub, args[0]+"cc")
		if err != nil {
			return shim.Error(err.Error())
		}

		if t.Length == t.N {
			// aa all online, start init

			// storage aalist
			err := t.aaListToSTR(stub)
			if err != nil {
				return shim.Error(err.Error())
			}

			// tell aa to communicate
			err = t.startAAsABECommunicate(stub)
			if err != nil {
				return shim.Error(err.Error())
			}

			// done
			t.Initialized = true
		}
		return shim.Success(nil)
	}
	//else
	return response
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
			return shim.Error("GetState PubKey error\n")
		}
		return shim.Success(PubKey)
	}
	return shim.Error("Invalid invoke function name. Expecting \"register\" \"getPubKey\"")
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

	t.AAList = make([]string, t.N)

	//chaincode's pair of keys
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