package middlewares

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ApiLogMiddleware(bodyIgnoredPaths []string) gin.HandlerFunc {
	bodyIgnoredPathsMap := make(map[string]bool)
	for _, path := range bodyIgnoredPaths {
		bodyIgnoredPathsMap[strings.ToLower(path)] = true
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// record body
		var body []byte
		if !bodyIgnoredPathsMap[strings.ToLower(path)] && c.Request.ContentLength < 5*1024 {
			body, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewReader(body))
		}

		// Process request
		c.Next()

		param := gin.LogFormatterParams{
			Request: c.Request,
			Keys:    c.Keys,
		}

		// Stop timer
		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)

		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.String()

		param.BodySize = c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		param.Path = path

		entry := logrus.WithFields(logrus.Fields{
			"status code":    param.StatusCode,
			"latency":        fmt.Sprintf("%13v", param.Latency),
			"client ip":      fmt.Sprintf("%15s", param.ClientIP),
			"method":         param.Method,
			"path":           param.Path,
			"full path":      c.FullPath(),
			"body":           string(body),
			"content length": c.Request.ContentLength,
		})

		if param.ErrorMessage != "" {
			entry = entry.
				WithField("errors", param.ErrorMessage).
				WithField("stack", c.GetString("error_stack"))

			if c.GetString("error_stack") != "" {
				fmt.Printf("Request error %v\n%v\n", param.ErrorMessage, c.GetString("error_stack"))
			}
		}
		entry.Info("[ApiLogMiddleware] Request")
	}
}
