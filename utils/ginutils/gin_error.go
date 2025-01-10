package ginutils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var HTTP_STATUS_BUISINESS = 599
var ERR_CORDE_NORMAL = 100

type GinErrorBody struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type GinError struct {
	HttpStatus int
	GinErrorBody
}

func NewGinError(httpStatus int, code int, message string) *GinError {
	var e GinError
	e.HttpStatus = httpStatus
	e.Code = code
	e.Message = message
	return &e
}

func NewBusinessGinError(code int, message string, data ...interface{}) *GinError {
	e := NewGinError(HTTP_STATUS_BUISINESS, code, message)
	if len(data) > 0 {
		e = e.WithData(data[0])
	}
	return e
}

func NewBadRequestGinError(code int, message string, data ...interface{}) *GinError {
	e := NewGinError(http.StatusBadRequest, code, message)
	if len(data) > 0 {
		e = e.WithData(data[0])
	}
	return e
}

func NewConflictGinError(code int, message string, data ...interface{}) *GinError {
	e := NewGinError(http.StatusConflict, code, message)
	if len(data) > 0 {
		e = e.WithData(data[0])
	}
	return e
}

func NewBusinessNormalGinError(message string, data ...interface{}) *GinError {
	return NewBusinessGinError(ERR_CORDE_NORMAL, message, data...)
}

func NewBadRequestNormalGinError(message string, data ...interface{}) *GinError {
	return NewBadRequestGinError(ERR_CORDE_NORMAL, message, data...)
}

func NewConflictNormalGinError(message string, data ...interface{}) *GinError {
	return NewConflictGinError(ERR_CORDE_NORMAL, message, data...)
}

func (r *GinError) WithData(data interface{}) *GinError {
	cloned := *r
	cloned.Data = data
	return &cloned
}

func (r *GinError) WithMessage(message string) *GinError {
	cloned := *r
	cloned.Message += ": " + message
	return &cloned
}

func (r *GinError) Error() string {
	if r.Data != nil {
		return fmt.Sprintf("%v: %v. Data: %v", r.Code, r.Message, r.Data)
	}
	return fmt.Sprintf("%v: %v", r.Code, r.Message)
}

func (r *GinError) Body() GinErrorBody {
	return r.GinErrorBody
}

func (r *GinError) Render(c *gin.Context) {
	c.JSON(r.HttpStatus, r.Body())
}

func (r *GinError) IsSameCode(err error) bool {
	causeErr := errors.Cause(err)
	if causeErr == nil {
		return false
	}

	if ginErr, ok := causeErr.(*GinError); ok {
		return ginErr.Code == r.Code
	}
	return false
}
