package client

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/skdiver33/gophkeeper/model"
)

func (client *KeeperClient) ReadDataFromServer() error {
	if client.UserData == nil {
		slog.Warn("No users data found. Call list data on server.")
		return nil
	}
	items := make([]string, 0)
	for _, i := range *client.UserData {
		items = append(items, i.Description)
	}
	userDataPromt := promptui.Select{
		Label: "Select data",
		Items: items,
	}
	index, _, err := userDataPromt.Run()
	if err != nil {
		return fmt.Errorf("run select data promt failed %w", err)
	}
	pkg, err := client.GetData((*client.UserData)[index])
	if err != nil {
		return fmt.Errorf("error get user`s data from server %w", err)
	}
	err = pkg.DecryptPkgData(client.CryptKey)
	if err != nil {
		return fmt.Errorf("error decrypt receive data %w", err)
	}
	switch pkg.MData.UploadType {
	case model.AuthDataType:
		ad := model.AuthData{}
		ad.FromBinary(pkg.Data)
		fmt.Println(ad)
	case model.BankCardType:
		bc := model.BankCardData{}
		bc.FromBinary(pkg.Data)
		fmt.Println(bc)
	case model.FileType:
		f, err := os.CreateTemp(".", "download-file.*.txt")
		if err != nil {
			slog.Error("create download file", "error", err)
		}
		defer f.Close()
		if _, err := f.Write(pkg.Data); err != nil {
			f.Close()
			return fmt.Errorf("error write tmpfile: %w", err)
		}
		fmt.Printf("download file %s\n", f.Name())
	}

	return nil
}
