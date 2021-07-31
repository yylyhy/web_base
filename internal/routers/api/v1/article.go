package v1

import (
	"web-base/global"
	"web-base/internal/service"
	"web-base/pkg/app"
	"web-base/pkg/convert"
	"web-base/pkg/errcode"
	"github.com/gin-gonic/gin"
)

type Article struct{}


func NewArticle() Article {
	return Article{}
}

func (t Article) Get(c *gin.Context)    {
	param:= service.ArticleRequest{ID: convert.StrTo(c.Param("id")).MustUInt32()}
	response:= app.NewResponse(c)
	valid,errs:= app.BindAndValid(c,&param)
	if !valid{
		global.Logger.Errorf("app.BindAndValid errs: %v",errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	article,err:= svc.GetArticle(&param)
	if err!= nil{
		global.Logger.Errorf("svc.GetArticle err: %v",err)
		response.ToErrorResponse(errcode.ErrorGetArticleFail)
		return
	}

	response.ToResponse(article)
	return
}
func (t Article) List(c *gin.Context)   {
	param:= service.ArticleListRequest{}
	response:= app.NewResponse(c)
	valid,errs:= app.BindAndValid(c,&param)
	if !valid{
		global.Logger.Errorf("app.BindAndValid errs: %v",errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc:= service.New(c.Request.Context())
	pager:= app.Pager{Page: app.GetPage(c),PageSize: app.GetPageSize(c)}
	articles,totalRows,err:= svc.GetArticleList(&param,&pager)
	if err != nil {
		global.Logger.Errorf("svc.GetArticleList err: %v",err)
		response.ToErrorResponse(errcode.ErrorGetArticlesFail)
		return
	}
	response.ToResponseList(articles,totalRows)

}
func (t Article) Create(c *gin.Context) {
	param:= service.CreateArticleRequest{}
	response:= app.NewResponse(c)
	valid,errs:= app.BindAndValid(c,&param)
	if !valid{
		global.Logger.Errorf("app.BindAndValid errs: %v",errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc:= service.New(c.Request.Context())
	err:= svc.CreateArticle(&param)
	if err != nil {
		global.Logger.Errorf("svc.CreateArticle err: %v",err)
		response.ToErrorResponse(errcode.ErrorCreateArticleFail)
		return
	}
	response.ToResponse(gin.H{})
	return
}
func (t Article) Update(c *gin.Context) {
	param:= service.UpdateArticleRequest{ID:convert.StrTo(c.Param("id")).MustUInt32()}
	response:= app.NewResponse(c)
	valid,errs:= app.BindAndValid(c,&param)
	if !valid{
		global.Logger.Errorf("app.BindAndValid errs: %v",errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc:= service.New(c.Request.Context())
	err:= svc.UpdateArticle(&param)
	if err != nil {
		global.Logger.Errorf("svc.UpdateArticle err: %v",err)
		response.ToErrorResponse(errcode.ErrorUpdateArticleFail)
		return
	}

	response.ToResponse(gin.H{})
	return
}
func (t Article) Delete(c *gin.Context) {
	param:= service.DeleteArticleRequest{ID:convert.StrTo(c.Param("id")).MustUInt32()}
	response:= app.NewResponse(c)
	valid,errs:= app.BindAndValid(c,&param)
	if !valid{
		global.Logger.Errorf("app.BindAndValid errs: %v",errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc:= service.New(c.Request.Context())
	err:= svc.DeleteArticle(&param)
	if err != nil {
		global.Logger.Errorf("svc.DeleteArticle err: %v",err)
		response.ToErrorResponse(errcode.ErrorDeleteArticleFail)
		return
	}

	response.ToResponse(gin.H{})
	return

}
