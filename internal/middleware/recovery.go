package middleware

import (
	"web-base/global"
	"web-base/pkg/app"
	"web-base/pkg/errcode"
	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc{
	return func(c *gin.Context) {
		defer func(){
			if err:= recover();err!= nil{
				global.Logger.WithCallersFrames().Errorf("panic recover err : %v",err)
				app.NewResponse(c).ToErrorResponse(errcode.ServerError)
				c.Abort()
			}
		}()
	}
}
