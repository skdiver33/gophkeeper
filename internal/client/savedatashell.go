package client

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/manifoldco/promptui"
	"github.com/skdiver33/gophkeeper/model"
	"github.com/skdiver33/gophkeeper/protocol"
)

func (client *KeeperClient) SaveData() error {
	prompt := promptui.Select{
		Label: "Select data type",
		Items: []string{"Authentication data", "Bank card data", "File"},
	}
	index, _, err := prompt.Run()
	if err != nil {
		return err
	}
	switch index {
	case 0:
		ad := model.AuthData{}
		err := survey.Ask(authDataQuest, &ad)
		if err != nil {
			return err
		}
		desc, err := client.GetDataDescription()
		if err != nil {
			return err
		}
		sendData, err := ad.ToBinary()
		if err != nil {
			return err
		}
		pkg, err := protocol.CreateProtoPackage(sendData, model.AuthDataType, *desc)
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
		fmt.Println("user authentification data successfull save")
	case 1:
		bankData := model.BankCardData{}
		err := survey.Ask(bankCardQuest, &bankData)
		if err != nil {
			return err
		}
		desc, err := client.GetDataDescription()
		if err != nil {
			return err
		}
		fmt.Println(bankData)
		sendData, err := bankData.ToBinary()
		if err != nil {
			return err
		}
		pkg, err := protocol.CreateProtoPackage(sendData, model.BankCardType, *desc)
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
		fmt.Println("bank card data successfull save")
	case 2:
		fd := model.FileData{}
		err := survey.Ask(fileQuest, &fd)
		if err != nil {
			return err
		}
		desc, err := client.GetDataDescription()
		if err != nil {
			return err
		}
		data, err := os.ReadFile(fd.Filename)
		if err != nil {
			return err
		}
		pkg, err := protocol.CreateProtoPackage(data, model.FileType, *desc)
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
		fmt.Println("user file successfull save")

	}
	return nil
}

func (client *KeeperClient) GetDataDescription() (*string, error) {

	var desc string
	err := survey.Ask(descQuest, &desc)
	if err != nil {
		slog.Error("error get data description")
		return nil, err
	}
	return &desc, nil
}
