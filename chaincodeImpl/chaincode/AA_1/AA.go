package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/thorweiyan/ABEPasswordPlatform/chaincodeImpl/wrapper"
	"strconv"
)

type Chaincode struct {
	MyId        string
	Initialized bool
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

func (t *Chaincode) registerToSYS(stub shim.ChaincodeStubInterface, pubkey string) pb.Response {
	// call SYS receiveRegisterFromAA
	return wrapper.Call(stub, []string{"SYScc", "register", t.MyId, pubkey})
}

func (t *Chaincode) updateAAList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if t.Initialized {
		return shim.Error("Already initialized")
	}
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	rightOwner, err := wrapper.IsOwner(stub, args)
	if !rightOwner {
		return shim.Error("updateAAList" + err.Error())
	}

	//从STR获取AAList
	aaList, err := t.aaFromSTR(stub)
	if err != nil {
		return shim.Error("updateAAList" + err.Error())
	}
	aaList = aaList[:len(aaList)-1]

	t.AAList = aaList[:]
	for _,i :=range t.AAList{
		fmt.Printf("asdf:%x\n",i)
	}

	//存储其他AA pubkey
	for i := 1; i <= len(t.AAList)+1; i++ {
		if "AA_"+strconv.Itoa(i) == t.MyId {
			continue
		}
		err = stub.PutState("AA_"+strconv.Itoa(i), []byte(aaList[0]))
		if err != nil {
			return shim.Error("startABE1:Put AA error: " + err.Error())
		}
		aaList = aaList[1:]
	}
	return shim.Success(nil)
}

//args: r s "startABE1"
func (t *Chaincode) startABE1(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if t.Initialized {
		return shim.Error("Already initialized")
	}
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	rightOwner, err := wrapper.IsOwner(stub, args)
	if !rightOwner {
		return shim.Error("startABE1" + err.Error())
	}

	//生成其他AA所需的Sij秘密并发送
	aaSijs := wrapper.AACommunicate(t.AAList)
	err = t.sendToAA(stub, aaSijs, "AASecret")
	if err != nil {
		return shim.Error("startABE1:Communicate Secret error: " + err.Error())
	}
	return shim.Success(nil)
}

//args: r s "startABE2"
func (t *Chaincode) startABE2(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if t.Initialized {
		return shim.Error("Already initialized")
	}
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	rightOwner, err := wrapper.IsOwner(stub, args)
	if !rightOwner {
		return shim.Error("startABE1" + err.Error())
	}

	//AAsetup2
	tempPki, tempAid := wrapper.AASetup2()

	//发送pk,aid
	t.PKi = append(t.PKi, tempPki)
	t.Aid = append(t.Aid, tempAid)
	err = t.sendToAA(stub, []string{string(tempPki)}, "AAPKi")
	if err != nil {
		return shim.Error("handleFromAA AASecret:Communicate PKi error: " + err.Error())
	}
	return shim.Success(nil)
}

//args: r s "startABE3"
func (t *Chaincode) startABE3(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if t.Initialized {
		return shim.Error("Already initialized")
	}
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	rightOwner, err := wrapper.IsOwner(stub, args)
	if !rightOwner {
		return shim.Error("startABE1" + err.Error())
	}

	//AAsetup3
	//集齐其他t-1个aa的Pki，生成e(g,g)^alpha
	wrapper.AASetup3(t.PKi, t.Aid)
	t.Initialized = true
	err = stub.PutState("Initialized", []byte("true"))
	if err != nil {
		return shim.Error("Put Initialized error: " + err.Error())
	}
	return shim.Success(nil)
}


//***************************  Communicate with AA  ***************************
//发送给其他AA, flag: "AASecret"/"AAPKi"
func (t *Chaincode) sendToAA(stub shim.ChaincodeStubInterface, params []string, flag string) error {
	fmt.Printf("sendtoaa:%s,%x\n",flag, params)
	if flag == "AASecret" {
		j := 0 //顺序对应的Sij
		for i := 1; i <= t.N; i++ {
			if "AA_"+strconv.Itoa(i) == t.MyId {
				continue
			}
			//call aa
			passParams, err := wrapper.SignTransaction(stub, []string{flag, params[j]})
			if err != nil {
				return fmt.Errorf("sendToAA AASecret" + err.Error())
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
			return fmt.Errorf("sendToAA AAPKi: Turn ID to int error: " + err.Error())
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
			pubKey,err := stub.GetState("PubKey")
			if err != nil {
				return fmt.Errorf("sendToAA AAPKi:Get PubKey state:"+err.Error())
			}
			passParams, err := wrapper.SignTransaction(stub, []string{flag, param, string(pubKey)})
			if err != nil {
				return fmt.Errorf("sendToAA AAPKi:"+err.Error())
			}
			passParams = append([]string{id + "cc", "handleFromAA", t.MyId}, passParams...)

			response := wrapper.Call(stub, passParams)
			if response.Status != 200 {
				return fmt.Errorf("sendToAA AAPKi:"+response.Message)
			}
		}
	}else if flag[:4] == "User"{
		pubKey,err := stub.GetState("PubKey")
		if err != nil {
			return fmt.Errorf("sendToAA AAPKi:Get PubKey state:"+err.Error())
		}
		passparams, err := wrapper.SignTransaction(stub, []string{flag, params[0], params[1], string(pubKey)})
		if err != nil {
			return fmt.Errorf(err.Error())
		}
		passparams = append([]string{params[2], "handleFromAA", t.MyId}, passparams...)
		response := wrapper.Call(stub, passparams)
		if response.Status != 200 {
			return fmt.Errorf("sendToAA User:"+response.Message)
		}
	}else {
		return fmt.Errorf("sendToAA User:Don't match all methods\n")
	}
	return nil
}

//args: AA_ID r s ("AASecret" "AASij")/("AAPKi" "PKi" "Aid")/("UserSignUp"/"UserChangePassword"/"UserGetTip" "UserName" "PartSk" "Aid")
func (t *Chaincode) handleFromAA(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	rightAA, err := wrapper.IsAA(stub, args)
	if !rightAA {
		return shim.Error(err.Error())
	}

	switch args[3] {
	case "AASecret":
		if len(args) != 5 {
			return shim.Error("handleFromAA AASecret:Incorrect number of arguments. Expecting 5")
		}
		wrapper.AppendSij(args[4])
		return shim.Success(nil)
	case "AAPKi":
		if len(args) != 6 {
			return shim.Error("handleFromAA AAPki:Incorrect number of arguments. Expecting 6")
		}
		t.PKi = append(t.PKi, []byte(args[4]))
		t.Aid = append(t.Aid, []byte(args[5]))
		return shim.Success(nil)
	case "UserSignUp":
		if len(args) != 7 {
			return shim.Error("handleFromAA UserSignUp:Incorrect number of arguments. Expecting 7")
		}
		if temp, ok := t.TempUserParams[args[4]]; !ok {
			return shim.Error("handleFromAA UserSignUp:Don't have this user")
		}else {
			temp.PartSk = append(temp.PartSk, []byte(args[5]))
			temp.Aid = append(temp.Aid, []byte(args[6]))
			t.TempUserParams[args[4]] = temp
			if len(temp.Aid) == t.T {
				//sk gen
				userSk := wrapper.KeyGen(temp.PartSk, temp.Aid)
				err = t.signUp(stub, temp, userSk)
				if err != nil {
					return shim.Error("handleFromAA UserSignUp: " + err.Error())
				}
				delete(t.TempUserParams, args[4])
			}
			return shim.Success(nil)
		}
	case "UserChangePassword":
		if len(args) != 7 {
			return shim.Error("handleFromAA UserChangePassword:Incorrect number of arguments. Expecting 7")
		}
		if temp, ok := t.TempUserParams[args[4]]; !ok {
			return shim.Error("handleFromAA UserChangePassword:Don't have this user")
		}else {
			temp.PartSk = append(temp.PartSk, []byte(args[5]))
			temp.Aid = append(temp.Aid, []byte(args[6]))
			t.TempUserParams[args[4]] = temp
			if len(temp.Aid) == t.T {
				//sk gen
				userSk := wrapper.KeyGen(temp.PartSk, temp.Aid)
				//修改密码
				err = t.changePassword(stub, temp, userSk)
				if err != nil {
					return shim.Error("handleFromAA UserChangePassword: " + err.Error())
				}
				delete(t.TempUserParams, args[4])
			}
			return shim.Success(nil)
		}
	case "UserGetTip":
		if len(args) != 7 {
			return shim.Error("handleFromAA UserGetTip:Incorrect number of arguments. Expecting 7")
		}
		if temp, ok := t.TempUserParams[args[4]]; !ok {
			return shim.Error("Don't have this user")
		}else {
			temp.PartSk = append(temp.PartSk, []byte(args[5]))
			temp.Aid = append(temp.Aid, []byte(args[6]))
			t.TempUserParams[args[4]] = temp
			if len(temp.Aid) == t.T {
				//sk gen
				userSk := wrapper.KeyGen(temp.PartSk, temp.Aid)
				//加密数据上链
				message, err := t.getTip(stub, temp, userSk)
				if err != nil {
					return shim.Error("handleFromAA UserSignUp: " + err.Error())
				}
				delete(t.TempUserParams, args[4])
				return shim.Success(message)
			}
			return shim.Success(nil)
		}
	default:
		return shim.Error("Invalid invoke function name. Expecting \"AASecret\" \"AAPKi\" \"UserSignUp\" \"UserChangePassword\" \"UserGetTip\" ")
	}
}

//***************************  Communicate with STR  ***************************
//**** get ****
func (t *Chaincode) aaFromSTR(stub shim.ChaincodeStubInterface) ([]string, error) {
	//从STR获取AAList
	passParams, err := wrapper.SignTransaction(stub, []string{t.MyId})
	if err != nil {
		return nil, fmt.Errorf("aaFromSTR"+err.Error())
	}
	passParams = append([]string{"STRcc", "get", "AAList", t.MyId}, passParams...)
	response := wrapper.Call(stub, passParams)
	if response.Status != 200 {
		return nil, fmt.Errorf("aaFromSTR"+response.Message)
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
		return fmt.Errorf("attrFromSTR + Deserialize User's Data error: " + err.Error())
	}
	passparams = append([]string{"STRcc", "get", "ABEAttr",t.MyId}, passparams...)
	response := wrapper.Call(stub, passparams)
	if response.Status != 200 {
		return fmt.Errorf("attrFromSTR:" + response.Message)
	}


	temp := wrapper.SplitStringbyn(string(response.Payload))
	nowAttr,err := strconv.Atoi(temp[1])
	if err != nil {
		return fmt.Errorf("attrFromSTR:Get NowAttr error!\n")
	}
	abeAttrs := []byte(temp[0])

	err = wrapper.UnMarshalMap(abeAttrs, nowAttr)
	if err != nil {
		return fmt.Errorf("attrFromSTR: Unmarshal ABE's map error: " + err.Error())
	}
	return nil
}

//**** put ****
func (t *Chaincode) userDataToSTR(stub shim.ChaincodeStubInterface, password []byte, userName string) error {
	//加盐hash数据,存储用户账户密码
	passwordSaltHash, err := wrapper.Pbkdf2New(password)
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
	cypherText := wrapper.Encrypt(userData.UserName, userData.ChangePasswordPolicy)
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
	UserTipData := wrapper.Encrypt(userData.GetTipMessage, userData.GetTipPolicy)
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
	attrs, nowattr, err := wrapper.MarshalMap()
	if err != nil {
		return fmt.Errorf("attrToSTR:" + err.Error())
	}

	passparams, err := wrapper.SignTransaction(stub, []string{string(attrs), string(strconv.Itoa(nowattr))})
	if err != nil {
		return fmt.Errorf("attrToSTR: " + err.Error())
	}
	passparams = append([]string{"STRcc", "put", "ABEAttr", t.MyId}, passparams...)
	response := wrapper.Call(stub, passparams)
	if response.Status != 200 {
		return fmt.Errorf("attrToSTR:" + response.Message)
	}
	return nil
}

//---------------------------  初始化之后     ---------------------------
//***************************  User method  ***************************
//args: method r s serializedData
func (t *Chaincode) userMethod(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("userMethod:Incorrect number of arguments. Expecting 4")
	}
	if !t.Initialized {
		return shim.Error("userMethod:Not initialized")
	}
	rightOwner, err := wrapper.IsOwner(stub, args[1:])
	if !rightOwner {
		return shim.Error("userMethod:"+err.Error())
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
		return shim.Error("userSignUp:Deserialize User's Data error: " + err.Error())
	}

	if tempUserParams.SpecialAAId == "" || tempUserParams.Aid != nil || tempUserParams.PartSk != nil || tempUserParams.UserName == "" || len(tempUserParams.UserAttributes) != 6{
		return shim.Error("userSignUp:Passing params something wrong")
	}

	//生成对应user's sk并发给special AA
	//先从STR处取得现在的ATTR
	err = t.attrFromSTR(stub)
	if err != nil {
		return shim.Error("userSignUp error: " + err.Error())
	}
	partSk, err := wrapper.PartUserSkGen(tempUserParams.UserAttributes)
	if err != nil {
		return shim.Error("userSignUp:Generate Part of User's sk error: " + err.Error())
	}
	pubKey,err := stub.GetState("PubKey")
	if err != nil {
		return shim.Error("sendToAA AAPKi:Get PubKey state:"+err.Error())
	}
	passparams, err := wrapper.SignTransaction(stub, []string{"UserSignUp",tempUserParams.UserName,string(partSk),string(pubKey)})
	if err != nil {
		return shim.Error("userSignUp: " + err.Error())
	}
	passparams = append([]string{tempUserParams.SpecialAAId+"cc", "handleFromAA", t.MyId}, passparams...)
	response := wrapper.Call(stub, passparams)
	return response
}

func (t *Chaincode) userChangePassword(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	tempUserParams, err := wrapper.DeserializeUserData([]byte(args[0]))
	if err != nil {
		return shim.Error("userChangePassword:Deserialize User's Data error: " + err.Error())
	}

	if tempUserParams.SpecialAAId == "" || tempUserParams.Aid != nil || tempUserParams.PartSk != nil || tempUserParams.UserName == "" || len(tempUserParams.UserAttributes) == 0 {
		return shim.Error("userChangePassword:Passing params something wrong")
	}

	//生成对应user's sk并发给special AA
	//先从STR处取得现在的ATTR
	err = t.attrFromSTR(stub)
	if err != nil {
		return shim.Error("userChangePassword error: " + err.Error())
	}
	partSk, err := wrapper.PartUserSkGen(tempUserParams.UserAttributes)
	if err != nil {
		return shim.Error("GuserChangePassword:enerate Part of User's sk error: " + err.Error())
	}
	pubKey,err := stub.GetState("PubKey")
	if err != nil {
		return shim.Error("sendToAA AAPKi:Get PubKey state:"+err.Error())
	}
	passparams, err := wrapper.SignTransaction(stub, []string{"UserChangePassword", tempUserParams.UserName,string(partSk),string(pubKey)})
	if err != nil {
		return shim.Error("userChangePassword: " + err.Error())
	}
	passparams = append([]string{tempUserParams.SpecialAAId+"cc", "handleFromAA", t.MyId}, passparams...)
	response := wrapper.Call(stub, passparams)
	return response
}

func (t *Chaincode) userGetTip(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	tempUserParams, err := wrapper.DeserializeUserData([]byte(args[0]))
	if err != nil {
		return shim.Error("userGetTip:Deserialize User's Data error: " + err.Error())
	}

	if tempUserParams.SpecialAAId == "" || tempUserParams.Aid != nil || tempUserParams.PartSk != nil || tempUserParams.UserName == "" || len(tempUserParams.UserAttributes) == 0 {
		return shim.Error("userGetTip:Passing params something wrong")
	}

	//生成对应user's sk并发给special AA
	//先从STR处取得现在的ATTR
	err = t.attrFromSTR(stub)
	if err != nil {
		return shim.Error("userGetTip error: " + err.Error())
	}
	partSk, err := wrapper.PartUserSkGen(tempUserParams.UserAttributes)
	if err != nil {
		return shim.Error("userGetTip:Generate Part of User's sk error: " + err.Error())
	}

	pubKey,err := stub.GetState("PubKey")
	if err != nil {
		return shim.Error("sendToAA AAPKi:Get PubKey state:"+err.Error())
	}
	passparams, err := wrapper.SignTransaction(stub, []string{"UserGetTip", tempUserParams.UserName,string(partSk),string(pubKey)})
	if err != nil {
		return shim.Error("userGetTip: " + err.Error())
	}
	passparams = append([]string{tempUserParams.SpecialAAId+"cc", "handleFromAA", t.MyId}, passparams...)
	response := wrapper.Call(stub, passparams)
	return response
}

//special aa method
func (t *Chaincode) userSignUpSpecial(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	tempUserParams, err := wrapper.DeserializeUserData([]byte(args[0]))
	if err != nil {
		return shim.Error("userSignUpSpecial:Deserialize User's Data error: " + err.Error())
	}

	if tempUserParams.Aid != nil || tempUserParams.PartSk != nil || tempUserParams.UserName == "" || len(tempUserParams.UserAttributes) != 6{
		return shim.Error("userSignUpSpecial:Passing params something wrong")
	}

	//生成对应user's sk
	//先从STR处取得现在的ATTR
	err = t.attrFromSTR(stub)
	if err != nil {
		return shim.Error("userSignUpSpecial error: " + err.Error())
	}
	//判断用户名是否已经存在
	if wrapper.IsUserExists(tempUserParams.UserName) {
		return shim.Error("userSignUpSpecial:UserName already exists")
	}
	//加上新的ATTR，并存储**************************其他没有这两步
	err = wrapper.AddAttr(tempUserParams.UserAttributes)
	if err != nil {
		return shim.Error("userSignUpSpecial error: " + err.Error())
	}

	err = t.attrToSTR(stub)
	if err != nil {
		return shim.Error("userSignUpSpecial error: " + err.Error())
	}

	partSk, err := wrapper.PartUserSkGen(tempUserParams.UserAttributes)
	if err != nil {
		return shim.Error("userSignUpSpecial:Generate Part of User's sk error: " + err.Error())
	}

	tempUserParams.PartSk = append(tempUserParams.PartSk, partSk)
	pubKey,err := stub.GetState("PubKey")
	if err != nil {
		return shim.Error("userSignupSpecial:Get PubKey state:"+err.Error())
	}
	tempUserParams.Aid = append(tempUserParams.Aid, pubKey)
	t.TempUserParams[tempUserParams.UserName] = *tempUserParams
	return shim.Success(nil)
}

//changePassword and getTip
func (t *Chaincode) userOthersSpecial(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	tempUserParams, err := wrapper.DeserializeUserData([]byte(args[0]))
	if err != nil {
		return shim.Error("userOthersSpecial:Deserialize User's Data error: " + err.Error())
	}

	if tempUserParams.Aid != nil || tempUserParams.PartSk != nil || tempUserParams.UserName == "" || len(tempUserParams.UserAttributes) == 0 {
		return shim.Error("userOthersSpecial:Passing params something wrong")
	}

	//生成对应user's sk
	//先从STR处取得现在的ATTR
	err = t.attrFromSTR(stub)
	if err != nil {
		return shim.Error("userOthersSpecial error: " + err.Error())
	}
	partSk, err := wrapper.PartUserSkGen(tempUserParams.UserAttributes)
	if err != nil {
		return shim.Error("userOthersSpecial:Generate Part of User's sk error: " + err.Error())
	}
	tempUserParams.PartSk = append(tempUserParams.PartSk, partSk)
	pubKey,err := stub.GetState("PubKey")
	if err != nil {
		return shim.Error("userOtherSpecial:Get PubKey state:"+err.Error())
	}
	tempUserParams.Aid = append(tempUserParams.Aid, pubKey)
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
	message,err := wrapper.Decrypt(userSk, cyperText)
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
	message, err := wrapper.Decrypt(userSk, cyperText)
	if err != nil {
		return nil, fmt.Errorf("changePassword: " + err.Error())
	}
	return message, nil
}

//***************************  Third party method  ***************************
//args: r s userName PasswordHash
func (t *Chaincode) thirdVerify(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//TODO 加入innitialize
	if len(args) !=4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	rightOwner, err := wrapper.IsOwner(stub, args)
	if !rightOwner {
		return shim.Error(err.Error())
	}
	//获取用户密钥hash
	hash, err := t.userDataFromSTR(stub, args[2])
	fmt.Printf("hash:%x\n",hash)
	if err != nil {
		return shim.Error("thirdVerify:" + err.Error())
	}

	if !wrapper.Pbkdf2Verify([]byte(args[3]), hash){
		return shim.Error("Password unmatched!")
	}
	return shim.Success(nil)
}
//恢复所有参数
//args: r s recoverParams
func (t *Chaincode) recoverParams(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) !=4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	rightOwner, err := wrapper.IsOwner(stub, args)
	if !rightOwner {
		return shim.Error(err.Error())
	}
	//获取initialize等
	temp, err := stub.GetState("Initialized")
	if err != nil {
		return shim.Error("Get Initialized error: " + err.Error())
	}
	if string(temp) == "true"{
		t.Initialized = true
	}else {
		t.Initialized = false
	}

	temp, err = stub.GetState("MyId")
	if err != nil {
		return shim.Error("Get MyId error: " + err.Error())
	}
	t.MyId = string(temp)

	temp, err = stub.GetState("N")
	if err != nil {
		return shim.Error("Get N error: " + err.Error())
	}
	t.N,err = strconv.Atoi(string(temp))
	if err != nil {
		return shim.Error("Get N error: " + err.Error())
	}

	temp, err = stub.GetState("T")
	if err != nil {
		return shim.Error("Get T error: " + err.Error())
	}
	t.T,err = strconv.Atoi(string(temp))
	if err != nil {
		return shim.Error("Get T error: " + err.Error())
	}

	return shim.Success(nil)
}


//***************************  Chaincode interface  ***************************
func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("System Invoke")
	function, args := stub.GetFunctionAndParameters()

	switch function {
	case "getPubKey":
		PubKey, err := stub.GetState("PubKey")
		if err != nil {
			return shim.Error("Invoke:GetState PubKey error\n")
		}
		return shim.Success(PubKey)
	case "updateAAList":
		return t.updateAAList(stub, args)
	case "startABE1":
		return t.startABE1(stub, args)
	case "startABE2":
		return t.startABE2(stub, args)
	case "startABE3":
		return t.startABE3(stub, args)
	case "handleFromAA":
		return t.handleFromAA(stub, args)
	case "userMethod":
		return t.userMethod(stub, args)
	case "thirdVerify":
		return t.thirdVerify(stub, args)
	case "recoverParams":
		return t.recoverParams(stub, args)
	default:
		return shim.Error("Invalid invoke function name. Expecting \"getPubKey\" " +
			"\"startABE1\" \"startABE2\" \"startABE3\" \"handleFromAA\" \"userMethod\" \"thirdVerify\"")
	}
}

//args: owner's_pubkey  AA_ID("AA_1") AA's_prikey AA's_pubkey
func (t *Chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("AA Init")
	_, args := stub.GetFunctionAndParameters()

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	//storage owner's pubkey
	err := stub.PutState("OwnerPubKey", []byte(args[0]))
	if err != nil {
		return shim.Error("Put Ownerkey Error!")
	}

	//storage my ID
	t.MyId = args[1]
	err = stub.PutState("MyId", []byte(args[1]))
	if err != nil {
		return shim.Error("Put MyId error: " + err.Error())
	}

	//chaincode's pair of keys
	CCPrikey, CCPubkey := []byte(args[2]),[]byte(args[3])

	err = stub.PutState("PriKey", CCPrikey)
	if err!= nil {
		return shim.Error("PutState Prikey error\n")
	}
	err = stub.PutState("PubKey", CCPubkey)
	if err!= nil {
		return shim.Error("PutState Pubkey error\n")
	}
	fmt.Printf("%x\n",CCPrikey)
	fmt.Printf("%x\n",CCPubkey)
	//开始注册
	response := t.registerToSYS(stub, string(CCPubkey))
	if response.Status != 200 {
		return response
	}
	mpk := response.Payload
	//存储并使用MPK初始化
	err = stub.PutState("PubKeyParams", []byte(mpk))
	if err != nil {
		return shim.Error("receiveMPK:Put PubKeyParams error: " + err.Error())
	}
	t.T,t.N = wrapper.AASetup1([]byte(mpk), CCPubkey)
	err = stub.PutState("N", []byte(strconv.Itoa(t.N)))
	if err != nil {
		return shim.Error("Put N error: " + err.Error())
	}
	err = stub.PutState("T", []byte(strconv.Itoa(t.T)))
	if err != nil {
		return shim.Error("Put T error: " + err.Error())
	}
	t.TempUserParams = make(map[string]wrapper.UserData,1000)

	t.Initialized = false
	err = stub.PutState("Initialized", []byte("false"))
	if err != nil {
		return shim.Error("Put Initialized error: " + err.Error())
	}
	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("Error starting SYSChaincode: %s", err)
	}
}