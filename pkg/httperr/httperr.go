package httperr

import "fmt"

type Error struct {
	Message string `json:"error"`
	Code    int    `json:"code"`
}

func (e *Error) Error() string { return fmt.Sprintf("%d: %s", e.Code, e.Message) }

func New(code int, msg string) *Error { return &Error{Message: msg, Code: code} }

func BadRequest(msg string) *Error   { return New(400, msg) }
func Unauthorized(msg string) *Error { return New(401, msg) }
func Forbidden(msg string) *Error    { return New(403, msg) }
func NotFound(msg string) *Error     { return New(404, msg) }
func Conflict(msg string) *Error     { return New(409, msg) }
func Internal(msg string) *Error     { return New(500, msg) }
