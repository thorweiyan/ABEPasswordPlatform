package wrapper

//Client functions
func SYSInit(t, n int) []byte {
	return nil
}

func AASetup1(aapubkey []byte)(int,int) {
	return 0,0
}

func AACommunicate(aaList []string) []string{
	return nil
}

func AppendSij(sij string) bool{
	return true
}

func AASetup2()([]byte,[]byte) {
	return nil,nil
}

func AASetup3(pki,aid [][]byte) {

}

func MarshalMap() ([]byte,error) {
	return nil, nil
}

func UnMarshalMap(attrs []byte) (error) {

	return nil
}

func IsUserExists(userName string)(bool) {

	return false
}

func AddAttr(userAttributes []string) error{
	return nil
}

func PartUserSkGen(attrs []string) ([]byte, error) {
	return nil,nil
}

func KeyGen(partSk [][]byte, aid [][]byte) []byte {
	return nil
}

func Encrypt(message string, policy string) []byte {
	return nil
}

func Decrypt(secretKey []byte, cryptText []byte) ([]byte, error) {
	return nil,nil
}
