package protocol

import (
	"crypto/sha256"
	"time"

	"github.com/skdiver33/gophkeeper/model"
)

type ProtocolPackage struct {
	MData model.Metadata `json:"metadata"`
	Data  []byte         `json:"data"`
}

func CreateProtoPackage(data []byte, dt model.DataTypes, descript string) (*ProtocolPackage, error) {
	pkg := ProtocolPackage{Data: data}
	md := model.Metadata{UploadDate: time.Now(), UploadType: dt, Description: descript}
	h := sha256.New()
	_, err := h.Write(data)
	if err != nil {
		return nil, err
	}
	md.Hash = string(h.Sum([]byte(md.Description)))
	pkg.MData = md
	return &pkg, nil
}
