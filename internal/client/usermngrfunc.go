package client

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/skdiver33/gophkeeper/model"
)

func (client *KeeperClient) AddUser() error {
	user := model.User{}
	err := survey.Ask(authQuest, &user)
	if err != nil {
		return fmt.Errorf("error get new user login: %w", err)
	}
	var passwd1, passwd2 string
	prompt := &survey.Password{
		Message: "Please type password",
	}
	err = survey.AskOne(prompt, &passwd1)
	if err != nil {
		return fmt.Errorf("error get new user password: %w", err)
	}
	prompt.Message = "Please retype your password"
	err = survey.AskOne(prompt, &passwd2)
	if err != nil {
		return fmt.Errorf("error get new user password: %w", err)
	}
	if strings.Compare(passwd1, passwd2) != 0 {
		return fmt.Errorf("passwd not compares")
	}
	user.Password = passwd1
	err = client.UserAuth(&user, "/api/user/register")
	if err != nil {
		return fmt.Errorf("user add error: %w", err)
	}
	prompt.Message = "User successful added and authorized.\n Create and remember your secret key.\nEnter your secret key."
	key := ""
	err = survey.AskOne(prompt, &key)
	if err != nil {
		return fmt.Errorf("error get user secret key: %w", err)
	}
	h := sha256.New()
	h.Write([]byte(key))
	client.CryptKey = h.Sum(nil)
	return nil
}

func (client *KeeperClient) AuthUser() error {
	user := model.User{}
	err := survey.Ask(authQuest, &user.Login)
	if err != nil {
		return fmt.Errorf("error get user login: %w", err)
	}
	prompt := &survey.Password{
		Message: "Please type your password",
	}
	err = survey.AskOne(prompt, &user.Password)
	if err != nil {
		return fmt.Errorf("error get user password: %w", err)
	}
	err = client.UserAuth(&user, "/api/user/login")
	if err != nil {
		return fmt.Errorf("user login error: %w", err)
	}
	prompt.Message = "User successful authorized.\n Enter your secret key."
	key := ""
	err = survey.AskOne(prompt, &key)
	if err != nil {
		return fmt.Errorf("error get user secret key: %w", err)
	}
	h := sha256.New()
	h.Write([]byte(key))
	client.CryptKey = h.Sum(nil)
	return nil
}
