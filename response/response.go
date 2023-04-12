package response

import "strings"

// Response is used for static shape json return
type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResError struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
}

// BuildResponse method is to inject data value to dynamic success response
func ResponseSuccess(message string, data interface{}) Response {
	res := Response{
		Status:  true,
		Message: message,
		Data:    data,
	}
	return res
}

// BuildErrorResponse method is to inject data value to dynamic failed response
func ResponseError(message string, err string) ResError {
	splittedError := strings.Split(err, "\n")
	res := ResError{
		Status:  false,
		Message: message,
		Errors:  splittedError,
	}
	return res
}
