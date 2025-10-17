package requestid

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const Header = "X-Request-ID"

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(Header)
		if id == "" {
			id = uuid.NewString()
			c.Request.Header.Set(Header, id)
		}
		c.Writer.Header().Set(Header, id)
		c.Set(Header, id)
		c.Next()
	}
}
