package model

import (
	"encoding/json"
	"time"
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	ID       int    `json:",omitempty"`
	JWT      string `json:"jwt,omitempty"`
}

type DataTypes int

const (
	Cred DataTypes = iota //0
	BankCard
	TextFile
	BinaryFile
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

func (ad *AuthData) ToBinary() ([]byte, error) {
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

func (bcd *BankCardData) ToBinary() ([]byte, error) {
	return json.Marshal(bcd)
}

func (bcd *BankCardData) FromBinary(input []byte) error {
	err := json.Unmarshal(input, bcd)
	if err != nil {
		return err
	}
	return nil
}

// func main() {
// 	aData := AuthData{Login: "user", Password: "user"}
// 	mData := Metadata{UploadDate: time.Now(), UploadType: PAI, Description: "Sberbank"}

// 	b, err := aData.ToBinary()
// 	if err != nil {
// 		return
// 	}
// 	sd := SendData{MData: mData, Data: b}
// 	send, err := json.Marshal(sd)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	data2 := SendData{}
// 	err = json.Unmarshal(send, &data2)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	log.Println(data2)
// 	auth := AuthData{}
// 	err = auth.FromBinary(data2.Data)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	log.Println(auth)
// }
