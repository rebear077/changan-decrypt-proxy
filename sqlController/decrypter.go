package sql

import (
	"io/ioutil"
	"os"

	encrypter "github.com/FISCO-BCOS/go-sdk/encryption"
	"github.com/sirupsen/logrus"
)

// 解密用的
type Decrypter struct {
	encrypte *encrypter.Encrypter
	symKey   []byte
	pubKey   []byte
	priKey   []byte
}

func NewDecrypter() *Decrypter {
	en := encrypter.NewEncrypter()
	symkey, err := getSymKey("configs/symPri.txt")
	if err != nil {
		logrus.Fatalln(err)
	}
	pubkey, err := getRSAPublicKey("configs/public.pem")
	if err != nil {
		logrus.Fatalln(err)
	}
	prikey, err := getRSAPrivateKey("configs/private.pem")
	if err != nil {
		logrus.Fatalln(err)
	}

	return &Decrypter{
		encrypte: en,
		symKey:   symkey,
		pubKey:   pubkey,
		priKey:   prikey,
	}
}

func (d *Decrypter) ValidateHash(hash []byte, plain []byte) bool {
	resHash := d.encrypte.Signature(plain)
	if string(resHash) == string(hash) {
		return true
	} else {
		return false
	}
}
func getSymKey(path string) ([]byte, error) {
	filesymPrivate, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	stat, err := filesymPrivate.Stat()
	if err != nil {
		return nil, err
	}
	symkey := make([]byte, stat.Size())
	filesymPrivate.Read(symkey)
	filesymPrivate.Close()
	return symkey, nil
}
func getRSAPublicKey(path string) ([]byte, error) {
	pubKey, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return pubKey, nil
}
func getRSAPrivateKey(path string) ([]byte, error) {
	privateKey, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return privateKey, err
}

func (d *Decrypter) DecryptSymkey(ensymkey []byte) ([]byte, error) {
	symkey, err := d.encrypte.AsymDecrypt(ensymkey, d.priKey)
	return symkey, err
}

func (d *Decrypter) DecryptData(endata string, symkey []byte) ([]byte, error) {
	data, err := d.encrypte.SymDecrypt([]byte(endata), symkey)
	return data, err
}
