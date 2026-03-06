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
	i18n          map[string]string
	messageSuffix string
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

func (r *GinError) WithI18n(translations map[string]string) *GinError {
	cloned := *r
	cloned.i18n = translations
	return &cloned
}

func (r *GinError) WithData(data interface{}) *GinError {
	cloned := *r
	cloned.Data = data
	return &cloned
}

func (r *GinError) WithMessage(message string) *GinError {
	cloned := *r
	suffix := ": " + message
	cloned.Message += suffix
	cloned.messageSuffix += suffix
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

func (r *GinError) BodyForAcceptLanguage(acceptLanguage string) GinErrorBody {
	if r == nil {
		return GinErrorBody{}
	}

	body := r.Body()
	localizedBaseMessage := r.localizedBaseMessage(acceptLanguage)
	if localizedBaseMessage == "" {
		return body
	}
	body.Message = localizedBaseMessage + r.messageSuffix
	return body
}

func (r *GinError) Render(c *gin.Context) {
	acceptLanguage := ""
	if c != nil && c.Request != nil {
		acceptLanguage = c.GetHeader("Accept-Language")
	}
	c.JSON(r.HttpStatus, r.BodyForAcceptLanguage(acceptLanguage))
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

func (r *GinError) localizedBaseMessage(acceptLanguage string) string {
	if len(r.i18n) == 0 {
		return ""
	}
	locale := LocaleFromAcceptLanguage(acceptLanguage)
	if msg, ok := r.i18n[locale]; ok {
		return msg
	}
	// Fallback to the first matched locale in i18n map, if Accept-Language doesn't match any locale in i18n map.
	// for example, providing "zh-CN" as WithI18n key but "zh-TW" in Accept-Language, it will fallback to "zh-CN" translation
	for key, msg := range r.i18n {
		if LocaleFromAcceptLanguage(key) == locale {
			return msg
		}
	}
	return ""
}
