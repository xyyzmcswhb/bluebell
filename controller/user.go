package controller

import (
	"errors"
	"fmt"
	"net/http"
	"web_app/dao/mysql"
	"web_app/logic"
	"web_app/models"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// SignUpHandler处理请求注册的函数
func SignupHandler(c *gin.Context) {
	var p models.ParamSignUp
	//1.获取参数和参数校验
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数错误,直接返回响应
		zap.L().Error("Sign up with invalid parameter", zap.Error(err))
		//判断err是否为validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			ResponseError(c, CodeInvalidParam)
			return
		}
		//翻译错误
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return

	}
	fmt.Println(p)
	//2.业务处理
	if err := logic.Signup(&p); err != nil {
		zap.L().Error("logic.Signup failed", zap.Error(err))
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": "注册失败",
		//})
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, nil)
}

// LoginHandler处理登录的函数
func LoginHandler(c *gin.Context) {
	//1. 获取请求参数及参数校验
	p := new(models.ParamLogIn)
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数错误,直接返回响应
		zap.L().Error("Log in with invalid parameter", zap.Error(err))
		//判断err是否为validator.ValidationErrors类型，如果不是的话则将其转换为该类型变量
		errs, ok := err.(validator.ValidationErrors)
		//判断错误是否为参数验证引起的
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": removeTopStruct(errs.Translate(trans)),
		//})
		//错误为参数验证引起，翻译该错误
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return

	}
	//2.业务逻辑处理
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": "用户名或密码错误",
		//})
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}

	//3.返回响应
	//c.JSON(http.StatusOK, gin.H{
	//	"msg": "登陆成功",
	//})
	ResponseSuccess(c, gin.H{
		"user_id":   fmt.Sprintf("%d", user.UserID), //id值大于1<<53-1	 int64类型的最大值是1<<63-1
		"user_name": user.Username,
		"token":     user.Token,
	})

}
