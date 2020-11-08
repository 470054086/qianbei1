package util

import (
	"fmt"
	"net/http"
	"qianbei.com/constat"
	"qianbei.com/core"
	"strings"
)

type response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const (
	DEFAULT_SUCCESS_CODE    = http.StatusOK
	DEFAULT_SUCCESS_MESSAGE = "请求成功"
	DEFAULT_ERROR_MESSAGE   = "请求错误"
	PARAMS_ERROR_CODE       = 400
	PARAMS_ERROR_MESSAGE    = "参数错误"
)

// 默认请求成功
func Success(data interface{}) *response {
	if data == nil {
		data = map[string]interface{}{}
	}
	return &response{
		Code:    DEFAULT_SUCCESS_CODE,
		Message: DEFAULT_SUCCESS_MESSAGE,
		Data:    data,
	}
}

func SuccessCodeMessage(code int, message string, data interface{}) *response {
	return &response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func Error() *response {

	return &response{
		Code:    http.StatusInternalServerError,
		Message: DEFAULT_ERROR_MESSAGE,
		Data:    map[string]interface{}{},
	}
}

//  默认错误 并且记录日志
func ErrorLog(err error) *response {
	addLog(err)
	var code int
	var message string
	switch err.(type) {
	// 抛出的异常
	case constat.ShowErrorMsg:
		code = err.(constat.ShowErrorMsg).Code
		message = err.(constat.ShowErrorMsg).Message
	default:
		// 服务器异常
		code = http.StatusInternalServerError
		message = DEFAULT_ERROR_MESSAGE
	}
	return &response{
		Code:    code,
		Message: message,
		Data:    map[string]interface{}{},
	}
}

// 参数错误
func ErrorParams() *response {
	return &response{
		Code:    PARAMS_ERROR_CODE,
		Message: PARAMS_ERROR_MESSAGE,
		Data:    map[string]interface{}{},
	}
}

func ErrorCode(code int) *response {
	return &response{
		Code:    code,
		Message: DEFAULT_ERROR_MESSAGE,
		Data:    map[string]interface{}{},
	}
}

// code和message的错误
func ErrorCodeMessage(code int, message string) *response {
	return &response{
		Code:    code,
		Message: message,
		Data:    map[string]interface{}{},
	}
}

func addLog(err error) {
	switch err.(type) {
	// 抛出的异常
	case constat.ShowErrorMsg:
		return
	default:
		// 服务器异常
		sprintf := fmt.Sprintf("%+v", err)
		replace := strings.Replace(sprintf, "\n", "\r\n", len(sprintf))
		core.QLog().Info(replace)
	}

}
