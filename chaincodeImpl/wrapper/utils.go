package wrapper

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"strings"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"math/big"
	"crypto/x509"
	"encoding/gob"
	"bytes"
	"io"
	"crypto/sha256"
	"golang.org/x/crypto/pbkdf2"
)


type UserData struct {
	UserName             string		//"UN:xxxxxx"
	UserPasswordHash     []byte		//"xxxxxxxxx"
	ChangePasswordPolicy string		//"CPP:xxxxx"
	GetTipPolicy         string		//"GTP:xxxxx"
	GetTipMessage        string		//"GTM:xxxxx"
	UserAttributes       []string	//"xxxxxxxxx"
	SpecialAAId			 string		//"AA_1"
	PartSk               [][]byte   //nil
	Aid                  [][]byte	//nil
}

func (u *UserData) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(u)
	if err != nil {
		return []byte{}, err
	}

	return result.Bytes(), nil
}

func DeserializeUserData(d []byte) (*UserData, error) {
	ud := new(UserData)

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&ud)
	if err != nil {
		return nil, err
	}
	return ud, nil
}

func Pbkdf2(password []byte) ([]byte, error) {
	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("rand salt error:" + err.Error())
	}

	return pbkdf2.Key(password, salt, 4096, 32, sha256.New), nil
}

func EcdsaSetUp() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error){
	prk, err := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("GenerateKey error: " + err.Error())
	}
	puk := prk.PublicKey

	return prk, &puk, nil
}

func EcdsaSign(prk *ecdsa.PrivateKey, sigMsg string ) ([]byte, []byte, error){
	r, s, err := ecdsa.Sign(rand.Reader, prk, []byte(sigMsg))
	if err != nil {
		return []byte(""), []byte(""), fmt.Errorf("Sign error: " + err.Error())
	}else {
		return r.Bytes(), s.Bytes(), nil
	}
}

func EcdsaVerify(puk *ecdsa.PublicKey, sigMsg string, r *big.Int, s *big.Int) (bool, error){
	isRight := ecdsa.Verify(puk, []byte(sigMsg), r, s)
	if isRight {
		return true, nil
	}else {
		return false, fmt.Errorf("Verify Error!\n")
	}
}

func EcdsaSetUpNormal() ([]byte, []byte, error){
	prk, err := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("GenerateKey error: " + err.Error())
	}
	puk := prk.PublicKey

	reprk, err := marshalEcdsaPrivateKey(prk)
	if err != nil {
		return []byte(""), []byte(""), fmt.Errorf(err.Error())
	}
	repuk := marshalEcdsaPublicKey(&puk)
	return reprk, repuk, nil
}

func EcdsaSignNormal(prk []byte, sigMsg string ) ([]byte, []byte, error){
	realprk, err := parseEcdsaPrivateKey(prk)
	if err != nil {
		return []byte(""), []byte(""), fmt.Errorf(err.Error())
	}

	r, s, err := ecdsa.Sign(rand.Reader, realprk, []byte(sigMsg))
	if err != nil {
		return []byte(""), []byte(""), fmt.Errorf("Sign error: " + err.Error())
	}else {
		return r.Bytes(), s.Bytes(), nil
	}
}

func EcdsaVerifyNormal(puk []byte, sigMsg string, r []byte, s []byte) (bool, error){
	realpuk := parseEcdsaPublicKey(puk)

	realr := big.NewInt(0)
	reals := big.NewInt(0)
	realr.SetBytes(r)
	reals.SetBytes(s)

	isRight := ecdsa.Verify(realpuk, []byte(sigMsg), realr, reals)
	if isRight {
		return true, nil
	}else {
		return false, fmt.Errorf("Verify Error!\n")
	}
}

func marshalEcdsaPrivateKey(prk *ecdsa.PrivateKey) ([]byte, error){
	return x509.MarshalECPrivateKey(prk)
}

func marshalEcdsaPublicKey(puk *ecdsa.PublicKey) []byte{
	return elliptic.Marshal(elliptic.P224(), puk.X, puk.Y)
}

func parseEcdsaPrivateKey(prk []byte) (*ecdsa.PrivateKey, error){
	return x509.ParseECPrivateKey(prk)
}

func parseEcdsaPublicKey(puk []byte) *ecdsa.PublicKey{
	re := new(ecdsa.PublicKey)
	re.Curve = elliptic.P224()
	re.X = big.NewInt(0)
	re.Y = big.NewInt(0)
	re.X, re.Y =  elliptic.Unmarshal(elliptic.P224(), puk)
	return re
}

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
