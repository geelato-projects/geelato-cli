package prompt

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/geelato/cli/pkg/logger"
)

type SelectOption struct {
	Name  string
	Value string
}

func Select(message string, options []SelectOption, defaultValue ...string) (string, error) {
	opts := make([]string, len(options))
	optionMap := make(map[string]string)
	for i, opt := range options {
		opts[i] = opt.Name
		optionMap[opt.Name] = opt.Value
	}

	var answer string
	prompt := &survey.Select{
		Message: message,
		Options: opts,
	}

	if len(defaultValue) > 0 {
		prompt.Default = defaultValue[0]
	}

	if err := survey.AskOne(prompt, &answer); err != nil {
		logger.Errorf("选择失败: %v", err)
		return "", err
	}

	return optionMap[answer], nil
}

func MultiSelect(message string, options []SelectOption) ([]string, error) {
	opts := make([]string, len(options))
	for i, opt := range options {
		opts[i] = opt.Name
	}

	var answers []string
	prompt := &survey.MultiSelect{
		Message: message,
		Options: opts,
	}

	if err := survey.AskOne(prompt, &answers); err != nil {
		logger.Errorf("多选失败: %v", err)
		return nil, err
	}

	return answers, nil
}

func Confirm(message string, defaultValue ...bool) (bool, error) {
	var answer bool
	prompt := &survey.Confirm{
		Message: message,
	}

	if len(defaultValue) > 0 {
		prompt.Default = defaultValue[0]
	}

	if err := survey.AskOne(prompt, &answer); err != nil {
		logger.Errorf("确认失败: %v", err)
		return false, err
	}

	return answer, nil
}

func Input(message string, defaultValue ...string) (string, error) {
	var answer string
	prompt := &survey.Input{
		Message: message,
	}

	if len(defaultValue) > 0 {
		prompt.Default = defaultValue[0]
	}

	if err := survey.AskOne(prompt, &answer); err != nil {
		logger.Errorf("输入失败: %v", err)
		return "", err
	}

	return answer, nil
}

func Password(message string) (string, error) {
	var answer string
	prompt := &survey.Password{
		Message: message,
	}

	if err := survey.AskOne(prompt, &answer); err != nil {
		logger.Errorf("密码输入失败: %v", err)
		return "", err
	}

	return answer, nil
}
