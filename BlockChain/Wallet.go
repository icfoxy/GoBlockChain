package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Addr       string
}

type Signature struct {
	R *big.Int
	S *big.Int
}

type WalletTansfer struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
	Addr       string `json:"addr"`
}

func NewWallet() *Wallet {
	w := new(Wallet)
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	w.PrivateKey = privateKey
	w.PublicKey = &privateKey.PublicKey
	w.Addr = PublicKeyToAddr(*w.PublicKey)
	return w
}

func (w *Wallet) GetPrivateKey() *ecdsa.PrivateKey {
	return w.PrivateKey
}

func (w *Wallet) GetPrivateKeyStr() string {
	return fmt.Sprintf("%x", w.PrivateKey.D.Bytes())
}

func (w *Wallet) GetPublicKey() *ecdsa.PublicKey {
	return w.PublicKey
}

func (w *Wallet) GetPublicKeyStr() string {
	return fmt.Sprintf("%x%x", w.PublicKey.X.Bytes(), w.PublicKey.Y.Bytes())
}

func (w *Wallet) GetAddr() string {
	return w.Addr
}

func PublicKeyToAddr(publicKey ecdsa.PublicKey) string {
	//step1 get publicKey
	//step2 SHA-256
	h2 := sha256.New()
	h2.Write(publicKey.X.Bytes())
	h2.Write(publicKey.Y.Bytes())
	digest2 := h2.Sum(nil)
	//step3 RIPEMD-160 (20 Btye)
	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)
	//step4 add version to front
	digest4 := make([]byte, 21)
	digest4[0] = 0x00 //main net
	copy(digest4[1:], digest3[:])
	//step5 SHA-256 on digest4
	h5 := sha256.New()
	h5.Write(digest4)
	digest5 := h5.Sum(nil)
	//step6 SHA-256 on digest5
	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)
	//step7 form checkSum
	checkSum := digest6[:4]
	//setp8 form final address bytes
	addrBytes := make([]byte, 25)
	copy(addrBytes[:21], digest4[:])
	copy(addrBytes[21:], checkSum[:])
	//step9 convert into string
	addr := base58.Encode(addrBytes)
	return addr
}

func (w *Wallet) SignTransaction(transaction *Transaction) *Signature {
	data, _ := json.Marshal(transaction)
	hash := sha256.Sum256([]byte(data))
	r, s, _ := ecdsa.Sign(rand.Reader, w.GetPrivateKey(), hash[:])
	return &Signature{
		R: r,
		S: s,
	}
}

func SignTransaction(transaction *Transaction, privateKey *ecdsa.PrivateKey) *Signature {
	data, _ := json.Marshal(transaction)
	hash := sha256.Sum256([]byte(data))
	r, s, _ := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	return &Signature{
		R: r,
		S: s,
	}
}

func (w *Wallet) NewSignedTransaction(
	senderAddr, receiverAddr string, value int, info string) (*Transaction, *Signature) {
	transaction := NewTransaction(senderAddr, receiverAddr, value, info)
	sign := w.SignTransaction(transaction)
	return transaction, sign
}

func (w *Wallet) ToTransfer() WalletTansfer {
	return WalletTansfer{
		PrivateKey: w.GetPrivateKeyStr(),
		PublicKey:  w.GetPublicKeyStr(),
		Addr:       w.Addr,
	}
}

func (s Signature) ToString() string {
	return fmt.Sprintf("%x%x", s.R, s.S)
}

func StringToPrivateKey(privateKeyStr string) (*ecdsa.PrivateKey, error) {
	privateKeyBytes, err := hex.DecodeString(privateKeyStr)
	if err != nil {
		return nil, err
	}
	d := new(big.Int).SetBytes(privateKeyBytes)
	return &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     nil,
			Y:     nil,
		},
		D: d,
	}, nil
}

func StringToPublicKey(publicKeyStr string) (*ecdsa.PublicKey, error) {
	publicKeyBytes, err := hex.DecodeString(publicKeyStr)
	if err != nil {
		return nil, err
	}
	x := new(big.Int).SetBytes(publicKeyBytes[:len(publicKeyBytes)/2])
	y := new(big.Int).SetBytes(publicKeyBytes[len(publicKeyBytes)/2:])
	return &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}, nil
}

func StringToSignature(signatureStr string) (*Signature, error) {
	signatureBytes, err := hex.DecodeString(signatureStr)
	if err != nil {
		return nil, err
	}
	r := new(big.Int).SetBytes(signatureBytes[:len(signatureBytes)/2])
	s := new(big.Int).SetBytes(signatureBytes[len(signatureBytes)/2:])
	return &Signature{
		R: r,
		S: s,
	}, nil
}
