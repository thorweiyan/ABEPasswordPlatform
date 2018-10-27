package wrapper

import (
	"testing"
	"fmt"
	"encoding/hex"
)

func TestSYS(t *testing.T) {
	pubkeybyte := SYSInit(2,3)
	//aa pubkey
	aa1Pubkey,_ := hex.DecodeString("04e3dd49e00dce869da09afc266d707a3e59377d28aded9d8f264ab890790aa92735f7ca9df8507a1d0823092e29ab7d74336dd9938521c479")
	aa2Pubkey,_ := hex.DecodeString("04a579d5b6764ae6b0f2f302d8717b4c8bb866ab6915be971798447d1918ec0ea739950ab7784a389e60a0f077ade960ae0e53d353f247c581")
	aa3Pubkey,_ := hex.DecodeString("04142100e66804198329ee8ac6e389f4d4448523f3cb13135f4fd4e0bfa816b00b7f5e53ae6d16c9a23dd8c0a7913934a6d19013a641a8cc8d")

	RPCADDRESS = "localhost:10000"
	AASetup1(pubkeybyte,aa1Pubkey)
	RPCADDRESS = "localhost:10001"
	AASetup1(pubkeybyte,aa2Pubkey)
	RPCADDRESS = "localhost:10002"
	AASetup1(pubkeybyte,aa3Pubkey)


	//aa communicate
	RPCADDRESS = "localhost:10000"
	aa1sij := AACommunicate([]string{string(aa2Pubkey),string(aa3Pubkey)})
	RPCADDRESS = "localhost:10001"
	aa2sij := AACommunicate([]string{string(aa1Pubkey),string(aa3Pubkey)})
	RPCADDRESS = "localhost:10002"
	aa3sij := AACommunicate([]string{string(aa1Pubkey),string(aa2Pubkey)})

	//aa setup2
	RPCADDRESS = "localhost:10000"
	AppendSij(aa2sij[0])
	fmt.Println(AppendSij(aa3sij[0]))
	aa1pk,aa1aid :=AASetup2()
	RPCADDRESS = "localhost:10001"
	AppendSij(aa1sij[0])
	fmt.Println(AppendSij(aa3sij[1]))
	aa2pk,aa2aid :=AASetup2()
	RPCADDRESS = "localhost:10002"
	AppendSij(aa1sij[1])
	fmt.Println(AppendSij(aa2sij[1]))
	aa3pk,aa3aid :=AASetup2()

	aa1needpki := [][]byte{aa1pk, aa2pk}
	aa1needaid := [][]byte{aa1aid, aa2aid}
	aa2needpki := [][]byte{aa2pk, aa3pk}
	aa2needaid := [][]byte{aa2aid, aa3aid}
	aa3needpki := [][]byte{aa3pk, aa1pk}
	aa3needaid := [][]byte{aa3aid, aa1aid}
	//aa setup3
	RPCADDRESS = "localhost:10000"
	AASetup3(aa1needpki,aa1needaid)
	RPCADDRESS = "localhost:10001"
	AASetup3(aa2needpki,aa2needaid)
	RPCADDRESS = "localhost:10002"
	AASetup3(aa3needpki,aa3needaid)
}

func TestKeyGen(t *testing.T) {
	aa1Pubkey,_ := hex.DecodeString("04e3dd49e00dce869da09afc266d707a3e59377d28aded9d8f264ab890790aa92735f7ca9df8507a1d0823092e29ab7d74336dd9938521c479")
	aa2Pubkey,_ := hex.DecodeString("04a579d5b6764ae6b0f2f302d8717b4c8bb866ab6915be971798447d1918ec0ea739950ab7784a389e60a0f077ade960ae0e53d353f247c581")
	//aa3Pubkey,_ := hex.DecodeString("04142100e66804198329ee8ac6e389f4d4448523f3cb13135f4fd4e0bfa816b00b7f5e53ae6d16c9a23dd8c0a7913934a6d19013a641a8cc8d")


	RPCADDRESS = "localhost:10000"
	AddAttr([]string{"czn","shuai","19970212","chou"})
	attrs, nowlen, err :=MarshalMap()
	fmt.Println(err)
	fmt.Println(nowlen)
	partsk1,err := PartUserSkGen([]string{"czn","shuai","19970212"})
	fmt.Println(err)

	RPCADDRESS = "localhost:10001"
	err =UnMarshalMap(attrs,nowlen)
	fmt.Println(err)
	partsk2,err := PartUserSkGen([]string{"czn","shuai","19970212"})
	fmt.Println(err)
	//
	//RPCADDRESS = "localhost:10002"
	//err =UnMarshalMap(attrs)
	//fmt.Println(err)
	//partsk3,err := PartUserSkGen([]string{"czn","shuai","19970212"})
	//fmt.Println(err)

	RPCADDRESS = "localhost:10000"
	usersk := KeyGen([][]byte{partsk1,partsk2},[][]byte{aa1Pubkey,aa2Pubkey})
	cy := Encrypt("this is test message", "(czn AND chou AND 19970212)")
	me, err := Decrypt(usersk, cy)
	fmt.Println(err)
	fmt.Println(string(me))
}

func TestUnMarshalMap(t *testing.T) {
	RPCADDRESS = "localhost:10000"
	attrs,nowlen, err :=MarshalMap()
	fmt.Println(err)
	fmt.Printf("%x\n", attrs)
	fmt.Println(UnMarshalMap(attrs,nowlen))
}
