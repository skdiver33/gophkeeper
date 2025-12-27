package protocol

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/skdiver33/gophkeeper/model"
)

type ProtocolPackage struct {
	MData model.Metadata `json:"metadata"`
	Data  []byte         `json:"data"`
}

func CreateProtoPackage(data []byte, dt model.DataTypes, descript string) (*ProtocolPackage, error) {
	pkg := ProtocolPackage{Data: data}
	md := model.Metadata{UploadDate: time.Now().UTC(), UploadType: dt, Description: descript}
	h := sha256.New()
	d := append(data, []byte(descript)...)
	_, err := h.Write(d)
	if err != nil {
		return nil, err
	}
	md.Hash = base64.RawStdEncoding.EncodeToString(h.Sum(nil))
	pkg.MData = md
	return &pkg, nil
}

func (pkg *ProtocolPackage) CryptPkgData(key []byte) error {

	aesblock, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return err
	}
	hash := []byte(pkg.MData.Hash)
	nonce := hash[len(hash)-aesgcm.NonceSize():]
	pkg.Data = aesgcm.Seal(nil, nonce, pkg.Data, nil)
	return nil
}

func (pkg *ProtocolPackage) DecryptPkgData(key []byte) error {
	aesblock, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return err
	}
	hash := []byte(pkg.MData.Hash)
	nonce := hash[len(hash)-aesgcm.NonceSize():]
	pkg.Data, err = aesgcm.Open(nil, nonce, pkg.Data, nil) // расшифровываем
	if err != nil {
		return err
	}
	return nil
}
