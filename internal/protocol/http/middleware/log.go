package middleware

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLoggerMiddleware(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody []byte
		if c.Request.Body != nil {
			reqBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(
				bytes.NewBuffer(reqBody),
			)
		}

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		fields := []zapcore.Field{
			zap.String("request_id", func() string {
				if reqId := c.Writer.Header().Get("X-Request-Id"); reqId != "" {
					return reqId
				}
				return uuid.New().String()
			}()),
			zap.String("request_body", string(reqBody)),
			zap.String("response_body", blw.body.String()),
		}

		log.With(fields...).Info("handled request")
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
