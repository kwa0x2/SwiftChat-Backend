package utils

import (
	"github.com/zishang520/engine.io/utils"
	socketUtils "github.com/zishang520/engine.io/utils"
)

// region "ExtractArgs" extracts the data and callback function from socket arguments.
func ExtractArgs(args []any) (map[string]interface{}, func([]interface{}, error)) {
	// Check if there are enough arguments
	if len(args) < 2 {
		utils.Log().Error(`not enough arguments`) // Log an error if not enough arguments are provided
		return nil, nil
	}

	// Extract data from the first argument and check its type
	data, ok := args[0].(map[string]interface{})
	if !ok {
		utils.Log().Error(`socket message type error`) // Log an error if the data type is incorrect
		return nil, nil
	}

	// Extract the callback function from the second argument and check its type
	callback, ok := args[1].(func([]interface{}, error))
	if !ok {
		utils.Log().Error(`callback function type error`) // Log an error if the callback type is incorrect
		return nil, nil
	}

	return data, callback // Return the extracted data and callback function
}

// endregion

// region Response defines a structured response format for socket communication.
type Response struct {
	Status  string `json:"status"`  // Response status (success/error)
	Message string `json:"message"` // Response message
}

// endregion

// region "SendResponse" sends a structured response back through the callback
func SendResponse(callback func([]interface{}, error), status, message string) {
	response := []interface{}{Response{Status: status, Message: message}} // Create a response object
	callback(response, nil)                                               // Invoke the callback with the response
}

// endregion

// region "LogError" logs an error message and sends an error response
func LogError(callback func([]interface{}, error), message string) {
	socketUtils.Log().Error(message)         // Log the error message
	SendResponse(callback, "error", message) // Send an error response
}

// endregion

// region "LogSuccess" sends a success response
func LogSuccess(callback func([]interface{}, error), message string) {
	SendResponse(callback, "success", message) // Send a success response
}

// endregion
