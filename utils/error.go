package utils

import (
	"encoding/json"
	"fmt"
)

type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (e *Error) Error() string {
	buf, _ := json.Marshal(e)
	return string(buf)
}

func (e *Error) MarshalJSON() ([]byte, error) {
	v := map[string]interface{}{}
	v["code"] = e.Code
	v["message"] = e.Message
	if e.Data != nil {
		v["data"] = e.Data
	}

	return json.Marshal(v)
}

func NewError(msg string, v ...interface{}) error {
	return NewErrorWithCode(10, msg, v...)
}

func NewErrorWithCode(code int, msg string, v ...interface{}) error {
	text := fmt.Sprintf(msg, v...)
	return &Error{Code: code, Message: text}
}

func NewInvalidJsonError() error {
	return NewErrorWithCode(20, "invalid json")
}

func NewInternelError() error {
	return NewErrorWithCode(500, "internel error")
}

func NewParamRequiredError(name string) error {
	return NewError(fmt.Sprintf("param required - %s", name))
}

func NewNotFoundError() error {
	return NewErrorWithCode(404, "not found")
}

func NewLoginRequiredError() error {
	return NewErrorWithCode(401, "login required")
}
func NewNoAccessPermissionError(msg string) error {
	if msg == "" {
		msg = "no  permissions"
	}
	return NewErrorWithCode(403, msg)
}
