package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/thorweiyan/ABEPasswordPlatform/chaincodeImpl/wrapper"
	"github.com/thorweiyan/MulticenterABEForFabric"
	"strconv"
	"encoding/json"
)

type Chaincode struct {
	MyId        string
	Initialized bool
	MAFF        *MulticenterABEForFabric.MAFFscheme
	AAList      []string //pubkey
	N           int      //所有aa的数量
	T           int      //阈值
	//ABE中间变量
	PKi [][]byte
	Aid [][]byte
	//User发送属性，special aa临时储存,username--userdata
	TempUserParams map[string]wrapper.UserData
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
	}else if flag[:4] == "User"{
		passparams, err := wrapper.SignTransaction(stub, []string{flag, params[0], params[1], string(t.Aid[0])})
		if err != nil {
			return fmt.Errorf(err.Error())
		}
		passparams = append([]string{params[2], "handleFromAA", t.MyId}, passparams...)
		response := wrapper.Call(stub, passparams)
		if response.Status != 200 {
			return fmt.Errorf(response.Message)
		}
	}else {
		return fmt.Errorf("Don't match all methods\n")
	}
	return nil
}

//args: AA_ID r s ("AASecret" "AASij")/("AAPKi" "PKi" "Aid")/("UserSignUp"/"UserChangePassword"/"UserGetTip" "UserName" "PartSk" "Aid")
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
	case "UserSignUp":
		if len(args) != 7 {
			return shim.Error("Incorrect number of arguments. Expecting 7")
		}
		if temp, ok := t.TempUserParams[args[4]]; !ok {
			return shim.Error("Don't have this user")
		}else {
			temp.PartSk = append(temp.PartSk, []byte(args[5]))
			temp.Aid = append(temp.Aid, []byte(args[6]))
			t.TempUserParams[args[4]] = temp
			if len(temp.Aid) == t.T {
				//sk gen
				userSk := t.keyGen(args[4])
				err = t.signUp(stub, temp, userSk)
				if err != nil {
					return shim.Error("UserSignUp: " + err.Error())
				}
				delete(t.TempUserParams, args[4])
			}
			return shim.Success(nil)
		}
	case "UserChangePassword":
		if len(args) != 7 {
			return shim.Error("Incorrect number of arguments. Expecting 7")
		}
		if temp, ok := t.TempUserParams[args[4]]; !ok {
			return shim.Error("Don't have this user")
		}else {
			temp.PartSk = append(temp.PartSk, []byte(args[5]))
			temp.Aid = append(temp.Aid, []byte(args[6]))
			t.TempUserParams[args[4]] = temp
			if len(temp.Aid) == t.T {
				//sk gen
				userSk := t.keyGen(args[4])
				//修改密码
				err = t.changePassword(stub, temp, userSk)
				if err != nil {
					return shim.Error("UserChangePassword: " + err.Error())
				}
				delete(t.TempUserParams, args[4])
			}
			return shim.Success(nil)
		}
	case "UserGetTip":
		if len(args) != 7 {
			return shim.Error("Incorrect number of arguments. Expecting 7")
		}
		if temp, ok := t.TempUserParams[args[4]]; !ok {
			return shim.Error("Don't have this user")
		}else {
			temp.PartSk = append(temp.PartSk, []byte(args[5]))
			temp.Aid = append(temp.Aid, []byte(args[6]))
			t.TempUserParams[args[4]] = temp
			if len(temp.Aid) == t.T {
				//sk gen
				userSk := t.keyGen(args[4])
				//加密数据上链
				message, err := t.getTip(stub, temp, userSk)
				if err != nil {
					return shim.Error("UserSignUp: " + err.Error())
				}
				delete(t.TempUserParams, args[4])
				return shim.Success(message)
			}
			return shim.Success(nil)
		}
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

func (t *Chaincode) userDataFromSTR(stub shim.ChaincodeStubInterface, userName string) ([]byte, error) {
	passparams, err := wrapper.SignTransaction(stub, []string{userName})
	if err != nil {
		return nil, fmt.Errorf("userDataFromSTR: " + err.Error())
	}
	passparams = append([]string{"STRcc", "get", "UserData", t.MyId}, passparams...)
	response := wrapper.Call(stub, passparams)
	if response.Status != 200 {
		return nil, fmt.Errorf("userDataFromSTR: " + response.Message)
	}

	return response.Payload, nil
}

func (t *Chaincode) userChangePasswordDataFromSTR(stub shim.ChaincodeStubInterface, userName string) ([]byte, error) {
	passparams, err := wrapper.SignTransaction(stub, []string{userName})
	if err != nil {
		return nil, fmt.Errorf("userChangePasswordDataFromSTR: " + err.Error())
	}
	passparams = append([]string{"STRcc", "get", "ChangePasswordData", t.MyId}, passparams...)
	response := wrapper.Call(stub, passparams)
	if response.Status != 200 {
		return nil, fmt.Errorf("userChangePasswordDataFromSTR: " + response.Message)
	}

	return response.Payload, nil
}

func (t *Chaincode) userTipFromSTR(stub shim.ChaincodeStubInterface, userName string) ([]byte, error) {
	passparams, err := wrapper.SignTransaction(stub, []string{userName})
	if err != nil {
		return nil, fmt.Errorf("userTipFromSTR: " + err.Error())
	}
	passparams = append([]string{"STRcc", "get", "UserTip",t.MyId}, passparams...)
	response := wrapper.Call(stub, passparams)
	if response.Status != 200 {
		return nil, fmt.Errorf("userTipFromSTR: " + response.Message)
	}

	return response.Payload, nil
}

func (t *Chaincode) attrFromSTR(stub shim.ChaincodeStubInterface) error {
	passparams, err := wrapper.SignTransaction(stub, []string{"ABEAttr"})
	if err != nil {
		return fmt.Errorf("Deserialize User's Data error: " + err.Error())
	}
	passparams = append([]string{"STRcc", "get", "ABEAttr",t.MyId}, passparams...)
	response := wrapper.Call(stub, passparams)
	if response.Status != 200 {
		return fmt.Errorf(response.Message)
	}
	err = json.Unmarshal(response.Payload, t.MAFF.Omega.Rhos_map)
	if err != nil {
		return fmt.Errorf("Unmarshal ABE's map error: " + err.Error())
	}
	return nil
}

//**** put ****
func (t *Chaincode) userDataToSTR(stub shim.ChaincodeStubInterface, password []byte, userName string) error {
	//加盐hash数据,存储用户账户密码
	passwordSaltHash, err := wrapper.Pbkdf2(password)
	if err != nil {
		return fmt.Errorf("UserDataToSTR: " + err.Error())
	}
	passparams, err := wrapper.SignTransaction(stub, []string{userName, string(passwordSaltHash)})
	if err != nil {
		return fmt.Errorf("UserDataToSTR: " + err.Error())
	}
	passparams = append([]string{"STRcc", "put", "UserData", t.MyId}, passparams...)
	response := wrapper.Call(stub, passparams)
	if response.Status != 200 {
		return fmt.Errorf("UserDataToSTR: " + response.Message)
	}

	return nil
}

func (t *Chaincode) userChangePasswordDataToSTR(stub shim.ChaincodeStubInterface, userData wrapper.UserData, userSk []byte) error {
	cypherText := t.encrypt(userData.UserName, userData.ChangePasswordPolicy)
	passparams, err := wrapper.SignTransaction(stub, []string{userData.UserName, string(cypherText)})
	if err != nil {
		return fmt.Errorf("userChangePasswordDataToSTR: " + err.Error())
	}
	passparams = append([]string{"STRcc", "put", "ChangePasswordData", t.MyId}, passparams...)
	response := wrapper.Call(stub, passparams)
	if response.Status != 200 {
		return fmt.Errorf("userChangePasswordDataToSTR: " + response.Message)
	}
	return nil
}

func (t *Chaincode) userTipToSTR(stub shim.ChaincodeStubInterface, userData wrapper.UserData, userSk []byte) error {
	UserTipData := t.encrypt(userData.GetTipMessage, userData.GetTipPolicy)
	passparams, err := wrapper.SignTransaction(stub, []string{userData.UserName, string(UserTipData)})
	if err != nil {
		return fmt.Errorf("UserTipToSTR: " + err.Error())
	}
	passparams = append([]string{"STRcc", "put", "UserTip", t.MyId}, passparams...)
	response := wrapper.Call(stub, passparams)
	if response.Status != 200 {
		return fmt.Errorf("UserTipToSTR: " + response.Message)
	}
	return nil
}

func (t *Chaincode) attrToSTR(stub shim.ChaincodeStubInterface) error {
	attrs, err := json.Marshal(t.MAFF.Omega.Rhos_map)
	if err != nil {
		return fmt.Errorf("Marshal ABE's map error: " + err.Error())
	}

	passparams, err := wrapper.SignTransaction(stub, []string{string(attrs)})
	if err != nil {
		return fmt.Errorf("attrToSTR: " + err.Error())
	}
	passparams = append([]string{"STRcc", "put", "ABEAttr", t.MyId}, passparams...)
	response := wrapper.Call(stub, passparams)
	if response.Status != 200 {
		return fmt.Errorf(response.Message)
	}
	return nil
}

//---------------------------  初始化之后     ---------------------------
//***************************  User method  ***************************
//args: method r s serializedData
func (t *Chaincode) userMethod(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	if t.Initialized {
		return shim.Error("Already initialized")
	}
	rightOwner, err := wrapper.IsOwner(stub, args[1:])
	if !rightOwner {
		return shim.Error(err.Error())
	}

	switch args[0] {
	case "userSignUp":
		return t.userSignUp(stub, args[3:])
	case "userSignUpSpecial":
		return t.userSignUpSpecial(stub, args[3:])
	case "userChangePassword":
		return t.userChangePassword(stub, args[3:])
	case "userChangePasswordSpecial":
		return t.userOthersSpecial(stub, args[3:])
	case "userGetTip":
		return t.userGetTip(stub, args[3:])
	case "userGetTipSpecial":
		return t.userOthersSpecial(stub, args[3:])
	default:
		return shim.Error("Invalid invoke function name. Expecting \"userSignUp\" \"userSignUpSpecial\" \"userChangePassword\" " +
			"\"userChangePasswordSpecial\" \"userGetTip\" \"userGetTipSpecial\" ")
	}
}

//args: r s serializedData
func (t *Chaincode) userSignUp(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	tempUserParams, err := wrapper.DeserializeUserData([]byte(args[0]))
	if err != nil {
		return shim.Error("Deserialize User's Data error: " + err.Error())
	}

	if tempUserParams.SpecialAAId == "" || tempUserParams.Aid != nil || tempUserParams.PartSk != nil || tempUserParams.UserName == "" || len(tempUserParams.UserAttributes) != 6{
		return shim.Error("Passing params something wrong")
	}

	//生成对应user's sk并发给special AA
	//先从STR处取得现在的ATTR
	err = t.attrFromSTR(stub)
	if err != nil {
		return shim.Error("userSignUp error: " + err.Error())
	}
	partSk, err := t.partUserSkGen(tempUserParams.UserAttributes)
	if err != nil {
		return shim.Error("Generate Part of User's sk error: " + err.Error())
	}

	passparams, err := wrapper.SignTransaction(stub, []string{string(partSk)})
	if err != nil {
		return shim.Error("userSignUp: " + err.Error())
	}
	passparams = append([]string{tempUserParams.SpecialAAId+"cc", "userMethod", "UserSignUp", t.MyId}, passparams...)
	response := wrapper.Call(stub, passparams)
	return response
}

func (t *Chaincode) userChangePassword(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	tempUserParams, err := wrapper.DeserializeUserData([]byte(args[0]))
	if err != nil {
		return shim.Error("Deserialize User's Data error: " + err.Error())
	}

	if tempUserParams.SpecialAAId == "" || tempUserParams.Aid != nil || tempUserParams.PartSk != nil || tempUserParams.UserName == "" || len(tempUserParams.UserAttributes) == 0 {
		return shim.Error("Passing params something wrong")
	}

	//生成对应user's sk并发给special AA
	//先从STR处取得现在的ATTR
	err = t.attrFromSTR(stub)
	if err != nil {
		return shim.Error("userChangePassword error: " + err.Error())
	}
	partSk, err := t.partUserSkGen(tempUserParams.UserAttributes)
	if err != nil {
		return shim.Error("Generate Part of User's sk error: " + err.Error())
	}

	passparams, err := wrapper.SignTransaction(stub, []string{string(partSk)})
	if err != nil {
		return shim.Error("userChangePassword: " + err.Error())
	}
	passparams = append([]string{tempUserParams.SpecialAAId+"cc", "userMethod", "UserChangePassword", t.MyId}, passparams...)
	response := wrapper.Call(stub, passparams)
	return response
}

func (t *Chaincode) userGetTip(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	tempUserParams, err := wrapper.DeserializeUserData([]byte(args[0]))
	if err != nil {
		return shim.Error("Deserialize User's Data error: " + err.Error())
	}

	if tempUserParams.SpecialAAId == "" || tempUserParams.Aid != nil || tempUserParams.PartSk != nil || tempUserParams.UserName == "" || len(tempUserParams.UserAttributes) == 0 {
		return shim.Error("Passing params something wrong")
	}

	//生成对应user's sk并发给special AA
	//先从STR处取得现在的ATTR
	err = t.attrFromSTR(stub)
	if err != nil {
		return shim.Error("userGetTip error: " + err.Error())
	}
	partSk, err := t.partUserSkGen(tempUserParams.UserAttributes)
	if err != nil {
		return shim.Error("Generate Part of User's sk error: " + err.Error())
	}

	passparams, err := wrapper.SignTransaction(stub, []string{string(partSk)})
	if err != nil {
		return shim.Error("userGetTip: " + err.Error())
	}
	passparams = append([]string{tempUserParams.SpecialAAId+"cc", "userMethod", "UserGetTip", t.MyId}, passparams...)
	response := wrapper.Call(stub, passparams)
	return response
}

//special aa method
func (t *Chaincode) userSignUpSpecial(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	tempUserParams, err := wrapper.DeserializeUserData([]byte(args[0]))
	if err != nil {
		return shim.Error("Deserialize User's Data error: " + err.Error())
	}

	if tempUserParams.SpecialAAId == "" || tempUserParams.Aid != nil || tempUserParams.PartSk != nil || tempUserParams.UserName == "" || len(tempUserParams.UserAttributes) != 6{
		return shim.Error("Passing params something wrong")
	}

	//生成对应user's sk
	//先从STR处取得现在的ATTR
	err = t.attrFromSTR(stub)
	if err != nil {
		return shim.Error("userSignUpSpecial error: " + err.Error())
	}
	//判断用户名是否已经存在
	if _,ok := t.MAFF.Omega.Rhos_map[tempUserParams.UserName]; ok {
		return shim.Error("UserName already exists")
	}
	//加上新的ATTR，并存储**************************其他没有这两步
	t.MAFF.AddAttr(tempUserParams.UserAttributes)
	err = t.attrToSTR(stub)
	if err != nil {
		return shim.Error("userSignUpSpecial error: " + err.Error())
	}

	partSk, err := t.partUserSkGen(tempUserParams.UserAttributes)
	if err != nil {
		return shim.Error("Generate Part of User's sk error: " + err.Error())
	}

	tempUserParams.PartSk = append(tempUserParams.PartSk, partSk)
	tempUserParams.Aid = append(tempUserParams.Aid, t.Aid[0])
	t.TempUserParams[tempUserParams.UserName] = *tempUserParams
	return shim.Success(nil)
}

//changePassword and getTip
func (t *Chaincode) userOthersSpecial(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	tempUserParams, err := wrapper.DeserializeUserData([]byte(args[0]))
	if err != nil {
		return shim.Error("Deserialize User's Data error: " + err.Error())
	}

	if tempUserParams.SpecialAAId == "" || tempUserParams.Aid != nil || tempUserParams.PartSk != nil || tempUserParams.UserName == "" || len(tempUserParams.UserAttributes) == 0 {
		return shim.Error("Passing params something wrong")
	}

	//生成对应user's sk
	//先从STR处取得现在的ATTR
	err = t.attrFromSTR(stub)
	if err != nil {
		return shim.Error("userOthersSpecial error: " + err.Error())
	}
	partSk, err := t.partUserSkGen(tempUserParams.UserAttributes)
	if err != nil {
		return shim.Error("Generate Part of User's sk error: " + err.Error())
	}
	tempUserParams.PartSk = append(tempUserParams.PartSk, partSk)
	tempUserParams.Aid = append(tempUserParams.Aid, t.Aid[0])
	t.TempUserParams[tempUserParams.UserName] = *tempUserParams
	return shim.Success(nil)
}

func (t *Chaincode) signUp(stub shim.ChaincodeStubInterface, userData wrapper.UserData, userSk []byte) error {
	//账户密码数据上链
	err := t.userDataToSTR(stub, userData.UserPasswordHash, userData.UserName)
	if err != nil {
		return fmt.Errorf("signUp: " + err.Error())
	}

	//修改密码凭证上链
	err = t.userChangePasswordDataToSTR(stub, userData, userSk)
	if err != nil {
		return fmt.Errorf("signUp: " + err.Error())
	}

	//存储用户提示信息如果有的话,没有直接返回
	if userData.GetTipMessage == "" || userData.GetTipPolicy == "" {
		return nil
	}
	err = t.userTipToSTR(stub, userData, userSk)
	if err != nil {
		return fmt.Errorf("signUp: " + err.Error())
	}
	return nil
}

func (t *Chaincode) changePassword(stub shim.ChaincodeStubInterface, userData wrapper.UserData, userSk []byte) error {
	//检测属性是否能够解密，先获取密文
	cyperText, err := t.userChangePasswordDataFromSTR(stub, userData.UserName)
	if err != nil {
		return fmt.Errorf("changePassword: " + err.Error())
	}
	message,err := t.decrypt(userSk, cyperText)
	if err != nil {
		return fmt.Errorf("changePassword: " + err.Error())
	}

	if string(message) == userData.UserName {
		err = t.userDataToSTR(stub, userData.UserPasswordHash, userData.UserName)
		if err != nil {
			return fmt.Errorf("changePassword: " + err.Error())
		}
		return nil
	}
	return fmt.Errorf("Policy not matched!\n")
}

func (t *Chaincode) getTip(stub shim.ChaincodeStubInterface, userData wrapper.UserData, userSk []byte) ([]byte, error) {
	//检测属性是否能够解密，先获取密文
	cyperText, err := t.userTipFromSTR(stub, userData.UserName)
	if err != nil {
		return nil, fmt.Errorf("getTip: " + err.Error())
	}
	message, err := t.decrypt(userSk, cyperText)
	if err != nil {
		return nil, fmt.Errorf("changePassword: " + err.Error())
	}
	return message, nil
}

//***************************  Third party method  ***************************
//args: r s userName PasswordHash
func (t *Chaincode) thirdVerify(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) !=4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	rightOwner, err := wrapper.IsOwner(stub, args)
	if !rightOwner {
		return shim.Error(err.Error())
	}
	//获取用户密钥hash
	hash, err := t.userDataFromSTR(stub, args[2])
	if err != nil {
		return shim.Error("thirdVerify:" + err.Error())
	}
	if string(hash) != args[3] {
		return shim.Error("Password unmatched!")
	}
	return shim.Success(nil)
}


//***************************  ABE usage  ***************************
func (t *Chaincode) partUserSkGen(attrs []string) ([]byte, error) {
	return t.MAFF.SKGEN_AA(attrs)
}

func (t *Chaincode) keyGen(userName string) []byte {
	return t.MAFF.SKGEN_USER(t.TempUserParams[userName].PartSk, t.TempUserParams[userName].Aid)
}

func (t *Chaincode) encrypt(message string, policy string) []byte {
	return t.MAFF.ENCRYPT([]byte(message), policy)
}

func (t *Chaincode) decrypt(secretKey []byte, cryptText []byte) ([]byte, error) {
	return t.MAFF.DECRYPT(secretKey, cryptText)
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
	case "userMethod":
		return t.userMethod(stub, args)
	case "thirdVerify":
		return t.thirdVerify(stub, args)
	default:
		return shim.Error("Invalid invoke function name. Expecting \"getPubKey\" \"registerToSYS\" \"receiveMPK\" " +
			"\"startABECommunicate\" \"handleFromAA\" \"userMethod\" \"thirdVerify\"")
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