package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/kuops/go-example-app/server/pkg/response"
	"github.com/kuops/go-example-app/server/pkg/store/redis"
	"github.com/kuops/go-example-app/server/pkg/utils/jwt"
	"strings"
)

func Middleware(cache redis.Interface,skipUris []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _,uri := range skipUris {
			if strings.HasPrefix(c.Request.RequestURI,uri) {
				c.Next()
				return
			}
		}

		token := c.Request.Header.Get("x-token")
		if token == "" {
			response.FailWithDetailed(gin.H{"reload": true}, "未登录或非法访问", c)
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(token)
		if err != nil {
			if err == jwt.TokenExpired {
				response.FailWithDetailed(gin.H{"reload": true}, "授权已过期", c)
				c.Abort()
				return
			}
			response.FailWithDetailed(gin.H{"reload": true}, "Token无法正常解析", c)
			c.Abort()
			return
		}

		if _,err := cache.Get(claims.UUID);err != nil {
			response.FailWithDetailed(gin.H{"reload": true}, "未登录或非法访问", c)
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}