package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/geelato/cli/pkg/logger"
	"github.com/geelato/cli/pkg/prompt"
	"github.com/spf13/cobra"
)

func NewApiCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Short: "API管理",
		Long:  `管理 API 脚本，包括创建、测试、运行等操作`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(apiCreateCmd)
	cmd.AddCommand(NewApiTestCmd())
	cmd.AddCommand(NewApiRunCmd())

	return cmd
}

var (
	apiType string
)

var apiCreateCmd = &cobra.Command{
	Use:   "create <api-name>",
	Short: "创建API",
	Long: `创建一个新的 API 脚本

示例:
  geelato api create getUserList
  geelato api create getUserList -t js
  geelato api create saveUser -t python
  geelato api create myHandler -t go`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var apiName string
		if len(args) > 0 {
			apiName = args[0]
		} else {
			input, err := prompt.Input("Enter API name (e.g., getUserList):")
			if err != nil {
				return fmt.Errorf("failed to input API name: %w", err)
			}
			apiName = strings.TrimSpace(input)
		}

		if apiName == "" {
			return fmt.Errorf("API name cannot be empty")
		}

		if apiType == "" {
			apiType = "js"
		}

		if err := createAPI(apiName, apiType); err != nil {
			return fmt.Errorf("failed to create API: %w", err)
		}

		logger.Infof("API '%s' created successfully!", apiName)
		logger.Info("")
		logger.Info("Created file:")
		logger.Info("  api/%s.api.%s", apiName, apiType)

		return nil
	},
}

func init() {
	apiCreateCmd.Flags().StringVarP(&apiType, "type", "t", "", "API type (js, python, go)")
}

func createAPI(apiName, apiType string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	geelatoPath := filepath.Join(cwd, "geelato.json")
	if _, err := os.Stat(geelatoPath); os.IsNotExist(err) {
		return fmt.Errorf("current directory is not a valid Geelato application")
	}

	apiDir := filepath.Join(cwd, "api")
	if err := os.MkdirAll(apiDir, 0755); err != nil {
		return fmt.Errorf("failed to create API directory: %w", err)
	}

	var content string
	switch strings.ToLower(apiType) {
	case "python", "py":
		content = generatePythonAPI(apiName)
	case "go":
		content = generateGoAPI(apiName)
	default:
		content = generateJSAPI(apiName)
	}

	fileName := fmt.Sprintf("%s.api.%s", apiName, apiType)
	filePath := filepath.Join(apiDir, fileName)

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create API file: %w", err)
	}

	return nil
}

func generateJSAPI(apiName string) string {
	return fmt.Sprintf(`/**
 * @api
 * @name %s
 * @path /api/%s
 * @method POST
 * @description API description
 * @version 1.0.0
 */

// @param
// name: param1
// type: String
// required: true
// description: Parameter 1

// @return
// type: Object
// description: Response data

(function() {
    // Get request parameters
    var param1 = $params.param1;

    // TODO: Add your business logic here

    return {
        code: 200,
        message: "success",
        data: {
            result: true
        }
    };
})();
`, apiName, strings.ToLower(apiName))
}

func generatePythonAPI(apiName string) string {
	return fmt.Sprintf(`# -*- coding: utf-8 -*-
"""
@api
@name %s
@path /api/%s
@method POST
@description API description
@version 1.0.0
"""

# @param
# name: param1
# type: String
# required: true
# description: Parameter 1

# @return
# type: Object
# description: Response data


def main(params):
    param1 = params.get('param1')

    # TODO: Add your business logic here

    return {
        'code': 200,
        'message': 'success',
        'data': {
            'result': True
        }
    }
`, apiName, strings.ToLower(apiName))
}

func generateGoAPI(apiName string) string {
	return fmt.Sprintf(`package main

/**
 * @api
 * @name %s
 * @path /api/%s
 * @method POST
 * @description API description
 * @version 1.0.0
 */

import (
    "github.com/geelato/gl"
)

func main() {
    // Get request parameters
    param1 := gl.GetParam("param1")

    // TODO: Add your business logic here

    gl.ReturnJSON(map[string]interface{}{
        "code":    200,
        "message": "success",
        "data": map[string]interface{}{
            "result": true,
        },
    })
}
`, apiName, strings.ToLower(apiName))
}

func NewApiTestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test <api-name>",
		Short: "测试API",
		Long: `本地测试 API 脚本

示例:
  geelato api test getUserList
  geelato api test getUserList --data '{"id": 1}'
  geelato api test saveUser --method POST --data '{"name": "test"}'`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			apiName := args[0]
			data, _ := cmd.Flags().GetString("data")
			method, _ := cmd.Flags().GetString("method")
			return runApiTest(apiName, data, method)
		},
	}

	cmd.Flags().String("data", "{}", "请求数据 (JSON格式)")
	cmd.Flags().String("method", "POST", "HTTP方法")
	cmd.Flags().String("url", "", "覆盖请求URL")

	return cmd
}

func runApiTest(apiName, data, method string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	apiPath := filepath.Join(cwd, "api", apiName+".api.js")
	if _, err := os.Stat(apiPath); os.IsNotExist(err) {
		return fmt.Errorf("API '%s' does not exist", apiName)
	}

	logger.Info("Testing API: %s", apiName)
	logger.Infof("  Method: %s", method)
	logger.Infof("  Data: %s", data)

	content, err := os.ReadFile(apiPath)
	if err != nil {
		return fmt.Errorf("failed to read API file: %w", err)
	}

	logger.Info("")
	logger.Info("API content preview:")
	logger.Info("-------------------")
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if i >= 20 {
			logger.Info("  ...")
			break
		}
		logger.Info("  %s", line)
	}

	logger.Info("")
	logger.Warn("API test requires running Geelato server")
	logger.Info("Please use 'geelato push' to deploy and test on server")

	return nil
}

func NewApiRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run <api-name>",
		Short: "运行API",
		Long: `运行 API 脚本（本地调试）

示例:
  geelato api run getUserList
  geelato api run saveUser --data '{"name": "test"}'`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			apiName := args[0]
			data, _ := cmd.Flags().GetString("data")
			return runApiRun(apiName, data)
		},
	}

	cmd.Flags().String("data", "{}", "请求数据 (JSON格式)")

	return cmd
}

func runApiRun(apiName, data string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	apiPath := filepath.Join(cwd, "api", apiName+".api.js")
	if _, err := os.Stat(apiPath); os.IsNotExist(err) {
		return fmt.Errorf("API '%s' does not exist", apiName)
	}

	logger.Info("Running API: %s", apiName)
	logger.Infof("  Data: %s", data)

	content, err := os.ReadFile(apiPath)
	if err != nil {
		return fmt.Errorf("failed to read API file: %w", err)
	}

	logger.Info("")
	logger.Info("API script:")
	logger.Info("-----------")
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if i >= 30 {
			logger.Info("  ...")
			break
		}
		logger.Info("  %s", line)
	}

	logger.Info("")
	logger.Warn("Local API execution requires Node.js runtime")
	logger.Info("For full testing, please deploy to Geelato server")

	return nil
}
