package client

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/manifoldco/promptui"
)

func (client *KeeperClient) RunUserShell(ctx context.Context) {

	prompt := promptui.Select{
		Label: "Select command",
		Items: []string{"Register new user", "User login", "Save data in keeper", "List data in keeper", "Read data from keeper", "Delete data from keeper", "Exit"},
	}

	for {

		select {
		case <-ctx.Done():
			return
		default:
			index, _, err := prompt.Run()
			if err != nil {
				slog.Warn("run select prompt", "error", err.Error())
				return
			}
			switch index {
			case 0:
				err := client.AddUser()
				if err != nil {
					fmt.Println("user add error")
					slog.Error("add user", "error", err.Error())
					continue
				}
				slog.Info("user successful add")
				fmt.Println("user successful add")
			case 1:
				err := client.AuthUser()
				if err != nil {
					fmt.Println("user auth error")
					slog.Error("user auth", "error", err.Error())
					continue
				}
				slog.Info("user successful auth")
			case 2:
				if client.JWT == "" {
					slog.Info("Please login in keeper")
					continue
				}
				err := client.SaveData()
				if err != nil {
					fmt.Println("save data error")
					slog.Error("save data", "error", err.Error())
					continue
				}
			case 3:
				if client.JWT == "" {
					slog.Info("Please login in keeper")
					continue
				}
				err := client.ListDataOnServer()
				if err != nil {
					fmt.Println("list data error")
					slog.Error("list data", "error", err.Error())
					continue
				}
			case 4:
				if client.JWT == "" {
					slog.Info("Please login in keeper")
					continue
				}
				err := client.ReadDataFromServer()
				if err != nil {
					fmt.Println("error read data from server")
					slog.Error("read data from server", "error", err.Error())
					continue
				}
			case 5:
				if client.JWT == "" {
					slog.Info("Please login in keeper")
					continue
				}
				err := client.DeleteDataFromServer()
				if err != nil {
					fmt.Println("error delete data from server")
					slog.Error("delete data from server", "error", err.Error())
					continue
				}
			case 6:
				return
			}
		}

	}
}

func (client *KeeperClient) ListDataOnServer() error {
	data, err := client.GetAllData()
	client.UserData = data
	if err != nil {
		return fmt.Errorf("err get data for user: %w", err)
	}
	for _, md := range *client.UserData {
		fmt.Printf("Data type %d data desription : %s  data hash: %s\n", md.UploadType, md.Description, md.Hash)
	}
	return nil

}
