package client

import (
	"fmt"
	"log/slog"

	"github.com/AlecAivazis/survey/v2"
	"github.com/manifoldco/promptui"
	"github.com/skdiver33/gophkeeper/model"
	"github.com/skdiver33/gophkeeper/protocol"
)

type BynaryConverter interface {
	ToBinary() ([]byte, error)
}

func GetProtocolPackage[K BynaryConverter](i int) (*protocol.ProtocolPackage, error) {

	var data K
	err := survey.Ask(dataQuest[i], &data)
	if err != nil {
		return nil, err
	}
	desc, err := GetDataDescription()
	if err != nil {
		return nil, err
	}
	sendData, err := data.ToBinary()
	if err != nil {
		return nil, err
	}
	pkg, err := protocol.CreateProtoPackage(sendData, model.DataTypes(i), *desc)
	if err != nil {
		return nil, err
	}
	return pkg, nil
}

func (client *KeeperClient) SaveData() error {
	prompt := promptui.Select{
		Label: "Select data type",
		Items: []string{"Authentication data", "Bank card data", "File"},
	}
	index, _, err := prompt.Run()
	if err != nil {
		return err
	}
	var pkg *protocol.ProtocolPackage
	switch index {
	case 0:
		pkg, err = GetProtocolPackage[model.AuthData](index)
	case 1:
		pkg, err = GetProtocolPackage[model.BankCardData](index)
	case 2:
		pkg, err = GetProtocolPackage[model.FileData](index)
	}
	if err != nil {
		return err
	}
	err = pkg.CryptPkgData(client.CryptKey)
	if err != nil {
		return err
	}
	err = client.SendData(pkg)
	if err != nil {
		return err
	}
	fmt.Println("data successfull save")

	return nil
}

func GetDataDescription() (*string, error) {
	var desc string
	err := survey.Ask(descQuest, &desc)
	if err != nil {
		slog.Error("error get data description")
		return nil, err
	}
	return &desc, nil
}
