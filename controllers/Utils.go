package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/thorweiyan/ABEPasswordPlatform/chaincodeImpl/wrapper"
	"github.com/thorweiyan/ABEPasswordPlatform/models"
	"math/big"
	"strconv"
)

func DoSdk(userdata wrapper.UserData, method string) (result string) {
	fmt.Printf("xxxxxxxxxxxxxxxx----------------%v\n",userdata)
	fmt.Printf("xxxxxxxxxxxxxxxx----------------%v\n",userdata.UserAttributes)
	fmt.Printf("xxxxxxxxxxxxxxxx----------------%v\n",len(userdata.UserAttributes))
	//return "sdad"
	//调用special AA合约
	special := Rand2(big.NewInt(0)).Int64()
	sId := "AA_" + strconv.Itoa(int(special))
	sccId := sId + "cc"
	fmt.Println("special ccId: ", sccId)

	sOwnerPriKey,_ := hex.DecodeString(AAkey[special-1].prikey)
	fmt.Println("sOwnerPriKey: ", sOwnerPriKey)

	fmt.Println("method: ", method+"Special")

	models.SdkUserMethods(sccId, sOwnerPriKey, userdata, method)


	//调用normal AA合约
	data := wrapper.UserData{
		UserName:             userdata.UserName,
		UserAttributes:       userdata.UserAttributes,
		SpecialAAId: sId,
	}
	fmt.Println(data)
	normal := Rand2(big.NewInt(special)).Int64()
	nccId := "AA_" + strconv.Itoa(int(normal)) + "cc"
	fmt.Println("normal ccId: ", nccId)

	nOwnerPriKey,_ := hex.DecodeString(AAkey[normal-1].prikey)
	fmt.Println("nOwnerPriKey: ", nOwnerPriKey)

	models.SdkUserMethodn(nccId, nOwnerPriKey, data, method)

	return "SUCCESS!"
}

//传进来一个数值，生成除了这个数值外的任意一个随机数
func Rand2(ex *big.Int) (index *big.Int) {
	max := big.NewInt(3)
	exp := ex.Sub(ex, big.NewInt(1))
	i, _ := rand.Int(rand.Reader, max)

	for {
		if i.String() == exp.String() {
			i, _ = rand.Int(rand.Reader, max)

		}else {
			break
		}
	}

	i.Add(i, big.NewInt(1))
	return i
}
