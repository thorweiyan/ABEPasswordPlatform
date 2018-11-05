package wrapper

import (
	"log"
	"net/rpc"
)

var RPCADDRESS string = "localhost:10000"

func dial() *rpc.Client {
	client, err := rpc.DialHTTP("tcp", RPCADDRESS)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	return client
}

type Sysinit struct {
	T, N int
}
//Client functions
func SYSInit(t, n int) []byte {
	client := dial()

	args := &Sysinit{T:t, N:n}
	var reply []byte
	err := client.Call("MAFF.SYSInit", args, &reply)
	if err != nil {
		log.Fatal("MAFF error:", err)
	}

	return reply
}

type Setup1 struct {
	Pubkey []byte
	AAaid []byte
}

func AASetup1(pubkey, aapubkey []byte)(int,int) {
	client := dial()

	var reply Sysinit
	setup1 := Setup1{Pubkey:pubkey,AAaid:aapubkey}
	err := client.Call("MAFF.AASetup1", setup1, &reply)
	if err != nil {
		log.Fatal("MAFF error:", err)
	}

	return reply.T,reply.N
}

func AACommunicate(aaList []string) []string{
	client := dial()

	var reply []string
	err := client.Call("MAFF.AACommunicate", aaList, &reply)
	if err != nil {
		log.Fatal("MAFF error:", err)
	}

	return reply
}

func AppendSij(sij string) bool{
	client := dial()

	var reply bool
	err := client.Call("MAFF.AppendSij", sij, &reply)
	if err != nil {
		log.Fatal("MAFF error:", err)
	}

	return reply
}

type Setup2 struct {
	Pki, Aid []byte
}
func AASetup2()([]byte,[]byte) {
	client := dial()

	var reply Setup2
	err := client.Call("MAFF.AASetup2", "", &reply)
	if err != nil {
		log.Fatal("MAFF error:", err)
	}

	return reply.Pki, reply.Aid
}

type Setup3 struct {
	Pki, Aid [][]byte
}
func AASetup3(pki,aid [][]byte) {
	client := dial()

	args := Setup3{Pki:pki, Aid:aid}
	var reply []byte
	err := client.Call("MAFF.AASetup3", args, &reply)
	if err != nil {
		log.Fatal("MAFF error:", err)
	}
}

type Mmap struct {
	Map []byte
	NowLen int
}

func MarshalMap() ([]byte,int,error) {
	client := dial()

	var reply *Mmap
	err := client.Call("MAFF.MarshalMap", "", &reply)
	if err != nil {
		log.Fatal("MAFF error:", err)
	}
	return reply.Map, reply.NowLen, err
}

func UnMarshalMap(attrs []byte, nowlen int) (error) {
	client := dial()

	var reply []byte
	args := Mmap{Map:attrs, NowLen:nowlen}
	err := client.Call("MAFF.UnMarshalMap", args, &reply)
	if err != nil {
		log.Fatal("MAFF error:", err)
	}
	return err
}

func IsUserExists(userName string)(bool) {
	client := dial()

	var reply bool
	err := client.Call("MAFF.IsUserExists", userName, &reply)
	if err != nil {
		log.Fatal("MAFF error:", err)
	}
	return reply
}

func AddAttr(userAttributes []string) error{
	client := dial()

	var reply []byte
	err := client.Call("MAFF.AddAttr", userAttributes, &reply)
	if err != nil {
		log.Fatal("MAFF error:", err)
	}
	return err
}

func PartUserSkGen(attrs []string) ([]byte, error) {
	client := dial()

	var reply []byte
	err := client.Call("MAFF.PartUserSkGen", attrs, &reply)
	if err != nil {
		log.Fatal("MAFF error:", err)
	}
	return reply,err
}

type Keygen struct {
	PartSk, Aid [][]byte
}
func KeyGen(partSk [][]byte, aid [][]byte) []byte {
	client := dial()

	args := Keygen{PartSk:partSk, Aid:aid}
	var reply []byte
	err := client.Call("MAFF.KeyGen", args, &reply)
	if err != nil {
		log.Fatal("MAFF error:", err)
	}
	return reply
}

type ENcrypt struct {
	Message, Policy string
}
func Encrypt(message string, policy string) []byte {
	client := dial()

	args := ENcrypt{Message:message, Policy:policy}
	var reply []byte
	err := client.Call("MAFF.Encrypt", args, &reply)
	if err != nil {
		log.Fatal("MAFF error:", err)
	}
	return reply
}

type DEcrypt struct {
	SecretKey, CryptText []byte
}
func Decrypt(secretKey []byte, cryptText []byte) ([]byte, error) {
	client := dial()

	args := DEcrypt{SecretKey:secretKey, CryptText:cryptText}
	var reply []byte
	err := client.Call("MAFF.Decrypt", args, &reply)
	if err != nil {
		log.Fatal("MAFF error:", err)
	}
	return reply, err
}
