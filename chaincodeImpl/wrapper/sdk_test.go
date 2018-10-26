package wrapper

import (
"testing"
"os"
"fmt"
"github.com/thorweiyan/fabric_go_sdk"
)
var fSetup = fabric_go_sdk.FabricSetup{
	// Network parameters
	OrdererID: "orderer.fudan.edu.cn",
	OrgID: "org1.fudan.edu.cn",

	// Channel parameters
	ChannelID:     "fudanfabric",
	ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/thorweiyan/fabric_go_sdk/fixtures/artifacts/fudanfabric.channel.tx",

	// Chaincode parameters
	ChainCodeID:     "",
	ChaincodeGoPath: os.Getenv("GOPATH"),
	ChaincodePath:   "github.com/thorweiyan/ABEPasswordPlatform/chaincodeImpl/chaincode/",
	ChaincodeVersion: "1",
	OrgAdmin:        "Admin",
	OrgName:         "org1",
	ConfigFile:      os.Getenv("GOPATH") + "/src/github.com/thorweiyan/fabric_go_sdk/config.yaml",

	// User parameters
	UserName: "User1",
}

//Just a example, need environment
func TestInitialize(t *testing.T) {
	// Initialization of the Fabric SDK from the previously set properties
	err := fSetup.Initialize()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
		return
	}
	// Close SDK
	defer fSetup.CloseSDK()
}

func TestInstallAndInstantiateCC_SYS(t *testing.T) {
	fSetup.ChaincodePath += "System/"
	fSetup.ChainCodeID = "SYScc"
	fSetup.ChaincodeVersion = "4"
	// Install and instantiate the chaincode
	err := fSetup.InstallAndInstantiateCC([]string{"init", "2", "3"})
	if err != nil {
		fmt.Printf("Unable to install and instantiate the chaincode: %v\n", err)
		return
	}
}

func TestInstallAndInstantiateCC_STR(t *testing.T) {
	fSetup.ChaincodePath = "Storage/"
	// Install and instantiate the chaincode
	err := fSetup.InstallAndInstantiateCC([]string{"init"})
	if err != nil {
		fmt.Printf("Unable to install and instantiate the chaincode: %v\n", err)
		return
	}
}

func TestInstallAndInstantiateCC_AA1(t *testing.T) {
	fSetup.ChaincodePath += "AA/"
	fSetup.ChainCodeID = "AA_1cc"
	fSetup.ChaincodeVersion = "1"
	// Install and instantiate the chaincode
	err := fSetup.InstallAndInstantiateCC([]string{"init"})
	if err != nil {
		fmt.Printf("Unable to install and instantiate the chaincode: %v\n", err)
		return
	}
}

func TestInstallAndInstantiateCC_AA2(t *testing.T) {
	fSetup.ChaincodePath = "AA/"
	// Install and instantiate the chaincode
	err := fSetup.InstallAndInstantiateCC([]string{"init"})
	if err != nil {
		fmt.Printf("Unable to install and instantiate the chaincode: %v\n", err)
		return
	}
}

func TestInstallAndInstantiateCC_AA3(t *testing.T) {
	fSetup.ChaincodePath = "AA/"
	// Install and instantiate the chaincode
	err := fSetup.InstallAndInstantiateCC([]string{"init"})
	if err != nil {
		fmt.Printf("Unable to install and instantiate the chaincode: %v\n", err)
		return
	}
}

func TestFabricSetup_Invoke(t *testing.T) {
	trcid, err := fSetup.Invoke([]string{"creator3"})
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	fmt.Println(trcid)
}

func TestFabricSetup_Query(t *testing.T) {
	payload, err := fSetup.Query([]string{"invoke", "query", "hello"})
	if err != nil {
		fmt.Println("query error!", err)
	}
	fmt.Println(payload)
}
