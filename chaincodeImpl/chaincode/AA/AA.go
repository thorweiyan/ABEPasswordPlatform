package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Chaincode struct {
}



func (t *Chaincode) signTransaction(stub shim.ChaincodeStubInterface, args []string, AAList []string) pb.Response {

}


//***************************  Communicate with SYS  ***************************
func (t *Chaincode) registerToSYS(stub shim.ChaincodeStubInterface, args []string, AAList []string) pb.Response {

}


//***************************  Communicate with AA  ***************************
func (t *Chaincode) sendToAA(stub shim.ChaincodeStubInterface, args []string, AAList []string) pb.Response {

}

func (t *Chaincode) handleFromAA(stub shim.ChaincodeStubInterface, args []string, AAList []string) pb.Response {

}

//***************************  Communicate with STR  ***************************
//**** get ****
func (t *Chaincode) aaFromSTR(stub shim.ChaincodeStubInterface, args []string, AAList []string) pb.Response {

}

func (t *Chaincode) userDataFromSTR(stub shim.ChaincodeStubInterface, args []string, AAList []string) pb.Response {

}

func (t *Chaincode) userChangePasswordFromSTR(stub shim.ChaincodeStubInterface, args []string, AAList []string) pb.Response {

}

func (t *Chaincode) userTipFromSTR(stub shim.ChaincodeStubInterface, args []string, AAList []string) pb.Response {

}

func (t *Chaincode) getFromSTR(stub shim.ChaincodeStubInterface, args []string, AAList []string) pb.Response {

}

//**** put ****
func (t *Chaincode) userDataToSTR(stub shim.ChaincodeStubInterface, args []string, AAList []string) pb.Response {

}

func (t *Chaincode) userChangePasswordToSTR(stub shim.ChaincodeStubInterface, args []string, AAList []string) pb.Response {

}

func (t *Chaincode) userTipToSTR(stub shim.ChaincodeStubInterface, args []string, AAList []string) pb.Response {

}

func (t *Chaincode) putToSTR(stub shim.ChaincodeStubInterface, args []string, AAList []string) pb.Response {

}


//***************************  User method  ***************************
func (t *Chaincode) userSignUp(stub shim.ChaincodeStubInterface, args []string, AAList []string) pb.Response {

}

func (t *Chaincode) userChangePassword(stub shim.ChaincodeStubInterface, args []string, AAList []string) pb.Response {

}

func (t *Chaincode) userGetTip(stub shim.ChaincodeStubInterface, args []string, AAList []string) pb.Response {

}

//special aa method
func (t *Chaincode) userSignUpSpecial(stub shim.ChaincodeStubInterface, args []string, AAList []string) pb.Response {

}

func (t *Chaincode) userChangePasswordSpecial(stub shim.ChaincodeStubInterface, args []string, AAList []string) pb.Response {

}

func (t *Chaincode) userGetTipSpecial(stub shim.ChaincodeStubInterface, args []string, AAList []string) pb.Response {

}

//***************************  ABE usage  ***************************
func (t *Chaincode) abeInit1(stub shim.ChaincodeStubInterface, args []string, AAList []string) pb.Response {

}

func (t *Chaincode) abeInit2(stub shim.ChaincodeStubInterface, args []string, AAList []string) pb.Response {

}

func (t *Chaincode) keyGen(stub shim.ChaincodeStubInterface, args []string, AAList []string) pb.Response {

}

func (t *Chaincode) encrypt(stub shim.ChaincodeStubInterface, args []string, AAList []string) pb.Response {

}

func (t *Chaincode) decrypt(stub shim.ChaincodeStubInterface, args []string, AAList []string) pb.Response {

}











//***************************  Chaincode interface  ***************************
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