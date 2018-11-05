package wrapper

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"math/big"
	"crypto/x509"
	"crypto/rand"
)

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
