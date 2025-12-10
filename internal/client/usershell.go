package client

import (
	"fmt"
	"log"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/manifoldco/promptui"

	"github.com/skdiver33/gophkeeper/model"
)

var lq = []*survey.Question{
	{
		Name:      "login",
		Prompt:    &survey.Input{Message: "Enter your login"},
		Validate:  survey.Required,
		Transform: survey.ToLower,
	},
}

func (client *KeeperClient) RunUserShell() {

	prompt := promptui.Select{
		Label: "Select command",
		Items: []string{"Register new user", "Authentifiacate user", "Save data in keeper", "List data in keeper", "Delete data from keeper", "Read data from keeper", "Exit"},
	}

	index, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch index {
	case 0:
		user := model.User{}
		err := survey.Ask(lq, &user)
		if err != nil {
			log.Println("error get new user login")
		}
		var passwd1, passwd2 string
		prompt := &survey.Password{
			Message: "Please type your password",
		}
		err = survey.AskOne(prompt, &passwd1)
		if err != nil {
			log.Println("error get new user password")
		}
		prompt.Message = "Please retype your password"
		err = survey.AskOne(prompt, &passwd2)
		if err != nil {
			log.Println("error get new user password")
		}
		if strings.Compare(passwd1, passwd2) != 0 {
			log.Println("passwd not compares")
			break
		}
		user.Password = passwd1
		client.RegisterNewUser(&user)
	case 1:

		err = survey.Ask(lq, &client.ClientUser.Login)
		if err != nil {
			log.Println("error get user login")
		}
		prompt := &survey.Password{
			Message: "Please type your password",
		}
		err = survey.AskOne(prompt, &client.ClientUser.Password)
		if err != nil {
			log.Println("error get user login")
		}
	case 2:
	case 3:
	case 4:
	case 6:
	case 7:

	default:
		return
	}
}
