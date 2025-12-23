package client

import (
	"fmt"
	"log/slog"

	"github.com/manifoldco/promptui"
)

func (client *KeeperClient) DeleteDataFromServer() error {
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
	err = client.DeleteData((*client.UserData)[index])
	if err != nil {
		return fmt.Errorf("error get user`s data from server %w", err)
	}

	return nil
}
