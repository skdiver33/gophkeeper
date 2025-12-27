package model

import (
	"encoding/json"
	"os"
	"time"
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	ID       int    `json:",omitempty"`
	JWT      string `json:"-"`
}

type DataTypes int

const (
	AuthDataType DataTypes = iota //0
	BankCardType
	FileType
)

type Metadata struct {
	UploadDate  time.Time `json:"uploaddate,format:unix"`
	UploadType  DataTypes `json:"uploadtype"`
	Description string    `json:"description"`
	Hash        string    `json:"hash"`
}

type AuthData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (ad AuthData) ToBinary() ([]byte, error) {
	return json.Marshal(ad)
}

func (ad *AuthData) FromBinary(input []byte) error {
	err := json.Unmarshal(input, ad)
	if err != nil {
		return err
	}
	return nil
}

type BankCardData struct {
	CardNumber string `json:"cardnumber"`
	ExpireDate string `json:"expiredate"`
	CSVCode    int    `json:"csvcode"`
	CardHolder string `json:"cardholder"`
}

func (bcd BankCardData) ToBinary() ([]byte, error) {
	return json.Marshal(bcd)
}

func (bcd *BankCardData) FromBinary(input []byte) error {
	err := json.Unmarshal(input, bcd)
	if err != nil {
		return err
	}
	return nil
}

type FileData struct {
	Filename string `json:"filename"`
}

func (file FileData) ToBinary() ([]byte, error) {
	data, err := os.ReadFile(file.Filename)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (file *FileData) FromBinary(input []byte) error {

	return nil
}
