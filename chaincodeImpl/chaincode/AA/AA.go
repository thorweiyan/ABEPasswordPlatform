package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/thorweiyan/ABEPasswordPlatform/chaincodeImpl/wrapper"
	"github.com/thorweiyan/MulticenterABEForFabric"
	"strconv"
)

type Chaincode struct {
	MyId string
	Initialized bool
	MAFF *MulticenterABEForFabric.MAFFscheme
	AAList []string //pubkey
	N int //所有aa的数量
	T int //阈值
	//ABE中间变量
	PKi [][]byte
	Aid [][]byte
	//User发送属性，special aa临时储存
	AttrSk map[]
}

//---------------------------  初始化阶段  ---------------------------------------
//***************************  Communicate with SYS  ***************************
//args: r s(for "registerToSYS")
func (t *Chaincode) registerToSYS(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if t.Initialized {
		return shim.Error("Already initialized")
	}
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	//is owner?
	payload := append(args, "registerToSYS")
	rightOwner, err := wrapper.IsOwner(stub, payload)
	if !rightOwner {
		return shim.Error(err.Error())
	}

	// call SYS receiveRegisterFromAA
	return wrapper.Call(stub, []string{"SYScc", "register", t.MyId})
}

//args: r s AA_ID(AA_1)
func (t *Chaincode) startABECommunicate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if t.Initialized {
		return shim.Error("Already initialized")
	}
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	rightSYS, err := wrapper.IsSYS(stub, args)
	if !rightSYS {
		return shim.Error(err.Error())
	}

	//从STR获取AAList
	aaList, err := t.aaFromSTR(stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	t.AAList = aaList[:]

	//存储其他AA pubkey
	for i := 1; i <= len(aaList)+1; i++ {
		if "AA_"+strconv.Itoa(i) == t.MyId {
			i++
		}
		err = stub.PutState("AA_"+strconv.Itoa(i), []byte(aaList[0]))
		if err != nil {
			return shim.Error("Put AA error: " + err.Error())
		}
		aaList = aaList[1:]
	}

	//生成其他AA所需的Sij秘密并发送
	aaSijs := t.MAFF.AACommunicate(t.AAList)
	err = t.sendToAA(stub, aaSijs, "AASecret")
	if err != nil {
		return shim.Error("Communicate Secret error: " + err.Error())
	}
	return shim.Success(nil)
}

//args: r s MPK
func (t *Chaincode) receiveMPK(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if t.Initialized {
		return shim.Error("Already initialized")
	}
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	//先获取SYS的公钥
	response := wrapper.Call(stub, []string{"SYScc", "getPubKey"})
	if response.Status != 200 {
		return response
	}

	//存储
	err := stub.PutState("SYSChaincode", response.Payload)
	if err != nil {
		return shim.Error("Put SYSChaincode error: " + err.Error())
	}

	//判断是否来自SYS
	rightSYS,err := wrapper.IsSYS(stub, args)
	if !rightSYS {
		return shim.Error(err.Error())
	}

	//存储并使用MPK初始化
	err = stub.PutState("PubKeyParams", []byte(args[2]))
	if err != nil {
		return shim.Error("Put PubKeyParams error: " + err.Error())
	}
	aapubkey, err := stub.GetState("PubKey")
	if err != nil {
		return shim.Error("Get AAPubKey error: " + err.Error())
	}
	t.MAFF.AASETUP1([]uint8(args[3]), aapubkey)
	t.N = t.MAFF.PublicKey.GetN()
	t.T = t.MAFF.PublicKey.GetT()

	return shim.Success(nil)
}

//***************************  Communicate with AA  ***************************
//发送给其他AA, flag: "AASecret"/"AAPKi"
func (t *Chaincode) sendToAA(stub shim.ChaincodeStubInterface, params []string, flag string) error {
	if flag == "AASecret" {
		j := 0 //顺序对应的Sij
		for i := 1; i <= t.N; i++ {
			if "AA_"+strconv.Itoa(i) == t.MyId {
				i++
			}
			//call aa
			passParams, err := wrapper.SignTransaction(stub, []string{flag, params[j]})
			if err != nil {
				return fmt.Errorf(err.Error())
			}
			passParams = append([]string{"AA_" + strconv.Itoa(i) + "cc", "handleFromAA", t.MyId}, passParams...)

			response := wrapper.Call(stub, passParams)
			if response.Status != 200 {
				return fmt.Errorf(response.Message)
			}
			j++
		}
	}else if flag == "AAPKi" {
		param := params[0]
		start,err := strconv.Atoi(t.MyId[3:])
		if err != nil {
			return fmt.Errorf("Turn ID to int error: " + err.Error())
		}
		//给自己id后的t-1个发PKi
		for i:= start+1; i <= start + t.T-1; i++ {
			j := i % (t.N+1)
			//不变说明在id后面,变小说明在前面
			var id string
			if j==i {
				id = "AA_"+strconv.Itoa(j)
			}else {
				id = "AA_"+strconv.Itoa(j+1)
			}

			passParams, err := wrapper.SignTransaction(stub, []string{flag, param, string(t.Aid[0])})
			if err != nil {
				return fmt.Errorf(err.Error())
			}
			passParams = append([]string{id + "cc", "handleFromAA", t.MyId}, passParams...)

			response := wrapper.Call(stub, passParams)
			if response.Status != 200 {
				return fmt.Errorf(response.Message)
			}
		}
	}else {
		return fmt.Errorf("Don't match all methods\n")
	}
	return nil
}

//args: AA_ID r s ("AASecret" "AASij")/("AAPKi" "PKi" "Aid")
func (t *Chaincode) handleFromAA(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if t.Initialized {
		return shim.Error("Already initialized")
	}
	rightAA, err := wrapper.IsAA(stub, args)
	if !rightAA {
		return shim.Error(err.Error())
	}

	switch args[3] {
	case "AASecret":
		if len(args) != 5 {
			return shim.Error("Incorrect number of arguments. Expecting 5")
		}
		t.MAFF.AppendSij(args[1])
		//集齐其他n个aa的秘密Sij,生成AA的SK和PK，将PK分享出去
		if len(t.MAFF.Sij) == t.N {
			t.MAFF.AASETUP2()
			t.PKi = append(t.PKi, t.MAFF.Pki.Bytes())
			t.Aid = append(t.Aid, t.MAFF.Aid.Bytes())
			err = t.sendToAA(stub, []string{string(t.MAFF.Pki.Bytes())}, "AAPKi")
			if err != nil {
				return shim.Error("Communicate PKi error: " + err.Error())
			}
		}
		return shim.Success(nil)
	case "AAPKi":
		if len(args) != 6 {
			return shim.Error("Incorrect number of arguments. Expecting 6")
		}
		t.PKi = append(t.PKi, []byte(args[4]))
		t.Aid = append(t.Aid, []byte(args[5]))
		//集齐其他t-1个aa的Pki，生成e(g,g)^alpha
		if len(t.PKi) == t.T {
			t.MAFF.AASETUP3(t.PKi, t.Aid)
			t.Initialized = true
		}
		return shim.Success(nil)
	default:
		return shim.Error("Invalid invoke function name. Expecting \"AASecret\" \"AAPKi\" ")
	}
}

//***************************  Communicate with STR  ***************************
//**** get ****
func (t *Chaincode) aaFromSTR(stub shim.ChaincodeStubInterface) ([]string, error) {
	//从STR获取AAList
	passParams, err := wrapper.SignTransaction(stub, []string{t.MyId})
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	passParams = append([]string{"STRcc", "get", "AAList", t.MyId}, passParams...)
	response := wrapper.Call(stub, passParams)
	if response.Status != 200 {
		return nil, fmt.Errorf(response.Message)
	}
	aaList := wrapper.SplitStringbyn(string(response.Payload))
	return aaList, nil
}

func (t *Chaincode) userDataFromSTR(stub shim.ChaincodeStubInterface, args []string) pb.Response {

}

func (t *Chaincode) userChangePasswordDataFromSTR(stub shim.ChaincodeStubInterface, args []string) pb.Response {

}

func (t *Chaincode) userTipFromSTR(stub shim.ChaincodeStubInterface, args []string) pb.Response {

}

func (t *Chaincode) getFromSTR(stub shim.ChaincodeStubInterface, args []string) pb.Response {

}

//**** put ****
func (t *Chaincode) userDataToSTR(stub shim.ChaincodeStubInterface, args []string) pb.Response {

}

func (t *Chaincode) userChangePasswordDataToSTR(stub shim.ChaincodeStubInterface, args []string) pb.Response {

}

func (t *Chaincode) userTipToSTR(stub shim.ChaincodeStubInterface, args []string) pb.Response {

}

func (t *Chaincode) putToSTR(stub shim.ChaincodeStubInterface, args []string) pb.Response {

}


//***************************  User method  ***************************
func (t *Chaincode) userSignUp(stub shim.ChaincodeStubInterface, args []string) pb.Response {

}

func (t *Chaincode) userChangePassword(stub shim.ChaincodeStubInterface, args []string) pb.Response {

}

func (t *Chaincode) userGetTip(stub shim.ChaincodeStubInterface, args []string) pb.Response {

}

//special aa method
func (t *Chaincode) userSignUpSpecial(stub shim.ChaincodeStubInterface, args []string) pb.Response {

}

func (t *Chaincode) userChangePasswordSpecial(stub shim.ChaincodeStubInterface, args []string) pb.Response {

}

func (t *Chaincode) userGetTipSpecial(stub shim.ChaincodeStubInterface, args []string) pb.Response {

}

//***************************  Third party method  ***************************
func (t *Chaincode) thirdVerify(stub shim.ChaincodeStubInterface, args []string) pb.Response {

}


//***************************  ABE usage  ***************************
func (t *Chaincode) keyGen(stub shim.ChaincodeStubInterface, args []string) pb.Response {

}

func (t *Chaincode) encrypt(stub shim.ChaincodeStubInterface, args []string) pb.Response {

}

func (t *Chaincode) decrypt(stub shim.ChaincodeStubInterface, args []string) pb.Response {

}











//***************************  Chaincode interface  ***************************
func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("System Invoke")
	function, args := stub.GetFunctionAndParameters()

	switch function {
	case "getPubKey":
		PubKey, err := stub.GetState("PubKey")
		if err != nil {
			return shim.Error("GetState PubKey error\n")
		}
		return shim.Success(PubKey)
	case "registerToSYS":
		return t.registerToSYS(stub, args)
	case "receiveMPK":
		return t.receiveMPK(stub, args)
	case "startABECommunicate":
		return t.startABECommunicate(stub, args)
	case "handleFromAA":
		return t.handleFromAA(stub, args)
	case "userSignUp":
		return t.userSignUp(stub, args)
	case "userSignUpSpecial":
		return t.userSignUpSpecial(stub, args)
	case "userChangePassword":
		return t.userChangePassword(stub, args)
	case "userChangePasswordSpecial":
		return t.userChangePasswordSpecial(stub, args)
	case "userGetTip":
		return t.userGetTip(stub, args)
	case "userGetTipSpecial":
		return t.userGetTipSpecial(stub, args)
	case "thirdVerify":
		return t.thirdVerify(stub, args)
	default:
		return shim.Error("Invalid invoke function name. Expecting \"getPubKey\" \"registerToSYS\" \"receiveMPK\" " +
			"\"startABECommunicate\" \"handleFromAA\" \"userSignUp\" \"userSignUpSpecial\" \"userChangePassword\" " +
				"\"userChangePasswordSpecial\" \"userGetTip\" \"userGetTipSpecial\" \"thirdVerify\"")
	}
}

//args: owner's_pubkey  AA_ID("AA_1")
func (t *Chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("AA Init")
	_, args := stub.GetFunctionAndParameters()

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	//storage owner's pubkey
	err := stub.PutState("OwnerPubKey", []byte(args[0]))
	if err != nil {
		return shim.Error("Put Ownerkey Error!")
	}

	//storage my ID
	t.MyId = args[1]

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