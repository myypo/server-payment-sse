package middleware

import "github.com/gin-gonic/gin"

func NewStreamHeadersMiddleware() gin.HandlerFunc {
	return func(g *gin.Context) {
		g.Writer.Header().Set("Content-Type", "text/event-stream")
		g.Writer.Header().Set("Cache-Control", "no-cache")
		g.Writer.Header().Set("Connection", "keep-alive")
		g.Writer.Header().Set("Transfer-Encoding", "chunked")
		g.Writer.Header().Set("X-Accel-Buffering", "no")
		g.Next()
	}
}
