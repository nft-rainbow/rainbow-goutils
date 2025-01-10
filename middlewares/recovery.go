package middlewares

import (
	"bytes"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/nft-rainbow/rainbow-goutils/utils/ginutils"
	"github.com/sirupsen/logrus"
)

func Recovery() gin.HandlerFunc {
	var buf bytes.Buffer
	return gin.CustomRecoveryWithWriter(&buf, gin.RecoveryFunc(func(c *gin.Context, err interface{}) {
		defer func() {
			logrus.WithField("recovered", buf.String()).WithField("error", err).Error("panic and recovery")
			buf.Reset()
		}()
		ginutils.RenderError(c, errors.New("internal server error"))
		c.Abort()
	}))
}
