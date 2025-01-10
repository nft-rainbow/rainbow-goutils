package ginutils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func RenderError(c *gin.Context, err error) {
	ginErr := ConvertToGinError(err)
	renderGinError(c, ginErr)
}

func renderGinError(c *gin.Context, err *GinError) {
	c.Error(err)
	c.Set("error_stack", fmt.Sprintf("%+v", errors.WithStack(err)))
	err.Render(c)
}

func RenderSuccess(c *gin.Context, obj interface{}) {
	if obj == nil {
		obj = make(map[string]interface{})
	}
	c.JSON(http.StatusOK, obj)
}

func RenderResponse(c *gin.Context, obj interface{}, err error) {
	if err != nil {
		ginErr, ok := errors.Cause(err).(*GinError)
		if ok {
			ginErr.Render(c)
		} else {
			NewBusinessNormalGinError(err.Error()).Render(c)
		}
		return
	}

	RenderSuccess(c, obj)
}

func ConvertToGinError(err error) *GinError {
	ginErr, ok := errors.Cause(err).(*GinError)
	if ok {
		return ginErr
	}
	return NewBusinessNormalGinError(err.Error())
}
