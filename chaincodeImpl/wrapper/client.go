package wrapper

import (
	"log"
	"net/rpc"
)

const RPCADDRESS = "localhost:9999"

func dial() *rpc.Client {
	client, err := rpc.DialHTTP("tcp", RPCADDRESS)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	return client
}

type sysinit struct {
	t, n int
}
//Client functions
func SYSInit(t, n int) []byte {
	client := dial()

	args := &sysinit{t:t,n:n}
	var reply []byte
	err := client.Call("MAFF.SYSInit", args, &reply)
	if err != nil {
		log.Fatal("MAFF error:", err)
	}

	return reply
}

func AASetup1(aapubkey []byte)(int,int) {
	client := dial()

	var reply sysinit
	err := client.Call("MAFF.AASetup1", aapubkey, &reply)
	if err != nil {
		log.Fatal("MAFF error:", err)
	}

	return reply.t,reply.n
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

type setup2 struct {
	pki, aid []byte
}
func AASetup2()([]byte,[]byte) {
	client := dial()

	var reply setup2
	err := client.Call("MAFF.AASetup2", "", &reply)
	if err != nil {
		log.Fatal("MAFF error:", err)
	}

	return reply.pki, reply.aid
}

type setup3 struct {
	pki,aid [][]byte
}
func AASetup3(pki,aid [][]byte) {
	client := dial()

	args := setup3{pki:pki, aid:aid}
	var reply []byte
	err := client.Call("MAFF.AASetup3", args, &reply)
	if err != nil {
		log.Fatal("MAFF error:", err)
	}
}


func MarshalMap() ([]byte,error) {
	client := dial()

	var reply []byte
	err := client.Call("MAFF.MarshalMap", "", &reply)
	if err != nil {
		log.Fatal("MAFF error:", err)
	}
	return reply, err
}

func UnMarshalMap(attrs []byte) (error) {
	client := dial()

	var reply []byte
	err := client.Call("MAFF.UnMarshalMap", attrs, &reply)
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

type keygen struct {
	partSk, aid [][]byte
}
func KeyGen(partSk [][]byte, aid [][]byte) []byte {
	client := dial()

	args := keygen{partSk:partSk, aid:aid}
	var reply []byte
	err := client.Call("MAFF.KeyGen", args, &reply)
	if err != nil {
		log.Fatal("MAFF error:", err)
	}
	return reply
}

type encrypt struct {
	message, policy string
}
func Encrypt(message string, policy string) []byte {
	client := dial()

	args := encrypt{message:message, policy:policy}
	var reply []byte
	err := client.Call("MAFF.Encrypt", args, &reply)
	if err != nil {
		log.Fatal("MAFF error:", err)
	}
	return reply
}

type decrypt struct {
	secretKey, cryptText []byte
}
func Decrypt(secretKey []byte, cryptText []byte) ([]byte, error) {
	client := dial()

	args := decrypt{secretKey:secretKey, cryptText:cryptText}
	var reply []byte
	err := client.Call("MAFF.Decrypt", args, &reply)
	if err != nil {
		log.Fatal("MAFF error:", err)
	}
	return reply, err
}
