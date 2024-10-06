package utils

import (
	"github.com/zishang520/engine.io/utils"
	socketUtils "github.com/zishang520/engine.io/utils"
)

func ExtractArgs(args []any) (map[string]interface{}, func([]interface{}, error)) {
	if len(args) < 2 {
		utils.Log().Error(`not enough arguments`)
		return nil, nil
	}

	data, ok := args[0].(map[string]interface{})
	if !ok {
		utils.Log().Error(`socket message type error`)
		return nil, nil
	}

	callback, ok := args[1].(func([]interface{}, error))
	if !ok {
		utils.Log().Error(`callback function type error`)
		return nil, nil
	}

	return data, callback
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// SendResponse sends a structured response back through the callback
func SendResponse(callback func([]interface{}, error), status, message string) {
	response := []interface{}{Response{Status: status, Message: message}}
	callback(response, nil)
}

// LogError logs an error message and sends an error response
func LogError(callback func([]interface{}, error), message string) {
	socketUtils.Log().Error(message)
	SendResponse(callback, "error", message)
}

// LogSuccess sends a success response
func LogSuccess(callback func([]interface{}, error), message string) {
	SendResponse(callback, "success", message)
}
