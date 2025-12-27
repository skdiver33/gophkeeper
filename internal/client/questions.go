package client

import "github.com/AlecAivazis/survey/v2"

var authQuest = []*survey.Question{
	{
		Name:      "login",
		Prompt:    &survey.Input{Message: "Enter your login"},
		Validate:  survey.Required,
		Transform: survey.ToLower,
	},
}

var authDataQuest = []*survey.Question{
	{
		Name:      "login",
		Prompt:    &survey.Input{Message: "Enter login"},
		Validate:  survey.Required,
		Transform: survey.ToLower,
	},
	{
		Name:      "password",
		Prompt:    &survey.Input{Message: "Enter password"},
		Validate:  survey.Required,
		Transform: survey.ToLower,
	},
}

var bankCardQuest = []*survey.Question{
	{
		Name:      "CardNumber",
		Prompt:    &survey.Input{Message: "Enter card number"},
		Validate:  survey.Required,
		Transform: survey.ToLower,
	},
	{
		Name:      "ExpireDate",
		Prompt:    &survey.Input{Message: "Enter expire date"},
		Validate:  survey.Required,
		Transform: survey.ToLower,
	},
	{
		Name:      "CSVCode",
		Prompt:    &survey.Input{Message: "Enter csv code"},
		Validate:  survey.Required,
		Transform: survey.ToLower,
	},
	{
		Name:      "CardHolder",
		Prompt:    &survey.Input{Message: "Enter card holder"},
		Validate:  survey.Required,
		Transform: survey.ToLower,
	},
}

var fileQuest = []*survey.Question{
	{
		Name:     "Filename",
		Prompt:   &survey.Input{Message: "Enter full filename (name + full path to file)"},
		Validate: survey.Required,
	},
}

var dataQuest = [][]*survey.Question{authDataQuest, bankCardQuest, fileQuest}

var descQuest = []*survey.Question{
	{
		Name:      "descripton",
		Prompt:    &survey.Input{Message: "Enter data description"},
		Validate:  survey.Required,
		Transform: survey.ToLower,
	},
}
