package casbin

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/kuops/go-example-app/server/pkg/log"
	"github.com/kuops/go-example-app/server/pkg/response"
	"github.com/kuops/go-example-app/server/pkg/utils/convert"
	"github.com/kuops/go-example-app/server/pkg/utils/jwt"
	"strings"
)

// 拦截器
func Middleware(skipUris []string,Enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _,uri := range skipUris {
			if strings.HasPrefix(c.Request.RequestURI,uri) {
				c.Next()
				return
			}
		}

		claims, _ := c.Get("claims")
		waitUse := claims.(*jwt.CustomClaims)
		// 获取请求的URI
		obj := c.Request.URL.Path
		// 获取请求方法
		act := c.Request.Method
		// 获取用户的角色
		sub := convert.ToString(waitUse.ID)
		// 判断策略中是否存在
		success, err := Enforcer.Enforce(sub, obj, act)
		if err != nil {
			log.Error(err)
			response.FailWithDetailed(gin.H{}, "权限不足", c)
			c.Abort()
			return
		}
		if success {
			c.Next()
		}  else  {
			response.FailWithDetailed(gin.H{}, "权限不足", c)
			c.Abort()
			return
		}
	}
}
