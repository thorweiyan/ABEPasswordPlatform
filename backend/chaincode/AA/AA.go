package AA

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"fmt"
	"encoding/base64"
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