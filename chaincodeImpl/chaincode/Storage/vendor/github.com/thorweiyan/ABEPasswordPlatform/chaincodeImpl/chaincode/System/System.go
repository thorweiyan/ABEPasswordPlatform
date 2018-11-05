package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"github.com/thorweiyan/ABEPasswordPlatform/chaincodeImpl/wrapper"
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
	pubKey := wrapper.SYSInit(t.T, t.N)
	err := stub.PutState("PubKeyParams", pubKey)
	if err != nil {
		return fmt.Errorf("abeInit:Put MPK of ABE wrong!\n")
	}
	return nil
}


//***************************  Communicate with STR  ***************************
func (t *Chaincode) aaListToSTR(stub shim.ChaincodeStubInterface) pb.Response {
	//if t.Length != t.N {
	//	return shim.Error("AA isn't all online")
	//}
	for i:=0;i<t.N;i++{
		response := wrapper.Call(stub, []string{"AA_"+strconv.Itoa(i+1)+"cc","getPubKey"})
		if response.Status != 200 {
			return response
		}
		t.AAList[i] = string(response.Payload)
	}

	passParams, err := wrapper.SignTransaction(stub, t.AAList)
	if err != nil {
		return shim.Error("aaListToSTR"+err.Error())
	}
	passParams = append([]string{"STRcc", "put", "AAList", "SYS"}, passParams...)
	response := wrapper.Call(stub, passParams)
	if response.Status != 200 {
		shim.Error("attrToSTR:" + response.Message)
	}
	//put ABEATTR
	attrs, nowattr, err := wrapper.MarshalMap()
	if err != nil {
		return shim.Error("attrToSTR:" + err.Error())
	}

	passparams, err := wrapper.SignTransaction(stub, []string{string(attrs), string(strconv.Itoa(nowattr))})
	if err != nil {
		return shim.Error("attrToSTR: " + err.Error())
	}
	passparams = append([]string{"STRcc", "put", "ABEAttr", "SYS"}, passparams...)
	response = wrapper.Call(stub, passparams)
	if response.Status != 200 {
		return shim.Error("attrToSTR:" + response.Message)
	}
	return shim.Success(nil)
}

//func (t *Chaincode) startAAsABECommunicate(stub shim.ChaincodeStubInterface) error {
//	for i,aa := range t.AAList {
//		//sign aa's pubkey
//		args, err := wrapper.SignTransaction(stub, []string{aa})
//		if err != nil {
//			return err
//		}
//		payload := append([]string{"AA_"+strconv.Itoa(i)+"cc", "startABECommunicate"}, args...)
//		tempResponse := wrapper.Call(stub, payload)
//		if tempResponse.Status != 200 {
//			return fmt.Errorf("StartAAsABECommmunicate error: " + tempResponse.Message)
//		}
//	}
//	return nil
//}

//***************************  Communicate with AA  ***************************

//args: AA_ID(AA_1、AA_2...) pubkey
func (t *Chaincode) receiveRegisterFromAA(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if t.Initialized {
		return shim.Error("receiveRegisterFromAA:Already Initialized")
	}
	//return MPK to AA
	abeMPK, err := stub.GetState("PubKeyParams")
	if err != nil {
		return shim.Error("sendParamsToAA:Get ABE's MPK error: " + err.Error())
	}

	////去掉AA_,get id
	//id, err := strconv.Atoi(args[0][3:])
	//if err != nil {
	//	return shim.Error("receiveRegisterFromAA:AA_ID's type error")
	//}

	//if t.AAList[id-1] != "" {
	//	return shim.Success(abeMPK)
	//}

	//t.AAList[id-1] = args[1]
	//t.Length ++

	//返回mpk
	return shim.Success(abeMPK)
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
			return shim.Error("Invoke:GetState PubKey error\n")
		}
		return shim.Success(PubKey)
	}else if function == "putaaList" {
		return t.aaListToSTR(stub)
	}
	return shim.Error("Invalid invoke function name. Expecting \"register\" \"getPubKey\"")
}

//args: t,n
func (t *Chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("System Init")
	_, args := stub.GetFunctionAndParameters()

	if len(args) != 2 {
		re := ""
		for _,v := range args{
			re+=v
		}
		return shim.Error(re+ "Incorrect number of arguments. Expecting 2")
	}

	var err error
	t.Length = 0
	t.T, err = strconv.Atoi(args[0])
	if err!= nil {
		return shim.Error("Init:t is not a string(int)\n")
	}
	t.N, err = strconv.Atoi(args[1])
	if err!= nil {
		return shim.Error("Init:n is not a string(int)\n")
	}
	t.Initialized = false

	t.AAList = make([]string, t.N)

	//chaincode's pair of keys
	CCPrikey, CCPubkey, err := wrapper.EcdsaSetUpNormal()
	if err!= nil {
		return shim.Error("Init::" + err.Error())
	}

	err = stub.PutState("PriKey", CCPrikey)
	if err!= nil {
		return shim.Error("Init:PutState Prikey error\n")
	}
	err = stub.PutState("PubKey", CCPubkey)
	if err!= nil {
		return shim.Error("Init:PutState Pubkey error\n")
	}

	err = t.abeInit(stub)
	if err!= nil {
		return shim.Error("Init:abe init error:" + err.Error())
	}
	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("Error starting SYSChaincode: %s", err)
	}
}