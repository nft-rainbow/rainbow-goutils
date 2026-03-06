package ginutils

import (
	"fmt"
	"net/http"
	"strings"

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

type ginErrorMessageSegment struct {
	raw  string
	i18n map[string]string
}

type GinError struct {
	HttpStatus      int
	Code            int
	Data            interface{}
	messageSegments []ginErrorMessageSegment
}

func NewGinError(httpStatus int, code int, message string) *GinError {
	var e GinError
	e.HttpStatus = httpStatus
	e.Code = code
	e.messageSegments = []ginErrorMessageSegment{{raw: message}}
	return &e
}

func NewLocalizedGinError(httpStatus int, code int, message string, translations map[string]string) *GinError {
	e := NewGinError(httpStatus, code, message)
	e.messageSegments[0].i18n = translations
	return e
}

func NewBusinessGinError(code int, message string, data ...interface{}) *GinError {
	e := NewGinError(HTTP_STATUS_BUISINESS, code, message)
	if len(data) > 0 {
		e = e.WithData(data[0])
	}
	return e
}

func NewLocalizedBusinessGinError(code int, message string, translations map[string]string, data ...interface{}) *GinError {
	e := NewLocalizedGinError(HTTP_STATUS_BUISINESS, code, message, translations)
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

func NewLocalizedBadRequestGinError(code int, message string, translations map[string]string, data ...interface{}) *GinError {
	e := NewLocalizedGinError(http.StatusBadRequest, code, message, translations)
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

func NewLocalizedConflictGinError(code int, message string, translations map[string]string, data ...interface{}) *GinError {
	e := NewLocalizedGinError(http.StatusConflict, code, message, translations)
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

func (r *GinError) WithMessage(message string, translations ...map[string]string) *GinError {
	cloned := *r
	cloned.messageSegments = cloneGinErrorMessageSegments(r.messageSegments)
	segment := ginErrorMessageSegment{raw: message}
	if len(translations) > 0 {
		segment.i18n = translations[0]
	}
	cloned.messageSegments = append(cloned.messageSegments, segment)
	return &cloned
}

func (r *GinError) Error() string {
	message := renderMessageSegments(r.messageSegments, "")
	if r.Data != nil {
		return fmt.Sprintf("%v: %v. Data: %v", r.Code, message, r.Data)
	}
	return fmt.Sprintf("%v: %v", r.Code, message)
}

func (r *GinError) Body() GinErrorBody {
	if r == nil {
		return GinErrorBody{}
	}

	return GinErrorBody{
		Code:    r.Code,
		Message: renderMessageSegments(r.messageSegments, ""),
		Data:    r.Data,
	}
}

func (r *GinError) BodyForAcceptLanguage(acceptLanguage string) GinErrorBody {
	if r == nil {
		return GinErrorBody{}
	}

	body := r.Body()
	body.Message = renderMessageSegments(r.messageSegments, acceptLanguage)
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

func cloneGinErrorMessageSegments(segments []ginErrorMessageSegment) []ginErrorMessageSegment {
	if len(segments) == 0 {
		return nil
	}
	cloned := make([]ginErrorMessageSegment, len(segments))
	copy(cloned, segments)
	return cloned
}

func renderMessageSegments(segments []ginErrorMessageSegment, acceptLanguage string) string {
	if len(segments) == 0 {
		return ""
	}

	rendered := make([]string, 0, len(segments))
	for _, segment := range segments {
		message := localizedSegmentMessage(segment, acceptLanguage)
		if message == "" {
			continue
		}
		rendered = append(rendered, message)
	}
	return strings.Join(rendered, ": ")
}

func localizedSegmentMessage(segment ginErrorMessageSegment, acceptLanguage string) string {
	if len(segment.i18n) == 0 {
		return segment.raw
	}

	locale := LocaleFromAcceptLanguage(acceptLanguage)
	if msg, ok := segment.i18n[locale]; ok {
		return msg
	}
	// Fallback to the first matched locale in i18n map, if Accept-Language doesn't match any locale in i18n map.
	// for example, providing "zh-CN" as WithI18n key but "zh-TW" in Accept-Language, it will fallback to "zh-CN" translation
	for key, msg := range segment.i18n {
		if LocaleFromAcceptLanguage(key) == locale {
			return msg
		}
	}
	return segment.raw
}
