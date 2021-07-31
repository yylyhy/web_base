package api

import (
	"web-base/global"
	"web-base/internal/service"
	"web-base/pkg/app"
	"web-base/pkg/errcode"
	"fmt"
	"github.com/gin-gonic/gin"
)

func GetAuth(c *gin.Context){
	param:= service.AuthRequest{}
	response:= app.NewResponse(c)
	valid,errs:= app.BindAndValid(c,&param)
	if !valid {
		global.Logger.Errorf("app.BindAndValid err: %v",errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc:= service.New(c.Request.Context())
	err:= svc.CheckAuth(&param)
	if err!= nil{
		global.Logger.Errorf("svc.CheckAut err: %v",err)
		response.ToErrorResponse(errcode.UnauthorizedAuthNotExist)
		return
	}
	fmt.Println(param)
	response.ToResponse(param)
	token,err:= app.GenerateToken(param.AppKey,param.AppSecret)
	if err!= nil{
		global.Logger.Errorf("svc.GenerateToken err: %v",err)
		response.ToErrorResponse(errcode.UnauthorizedokenGenerate)
		return
	}

	response.ToResponse(gin.H{"token":token,})
}
