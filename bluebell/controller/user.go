package controller

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// 验证作用逻辑代码
//SignUpHandler 处理注册请求的函数
func SignUpHandler(c *gin.Context){
	// 1. 获取参数和参数校验
	//var p models.ParamSignUp //定义校验的结构体变量
	p := new(models.ParamSignUp)
	// 如果下面的逻辑判断, 如果少字段或者不是JSON数据, 则直接返回错误 这里只是简单的校验, 复杂的自己写!!!!!
	if err := c.ShouldBindJSON(&p); err != nil{
		// 请求参数有误, 直接返回响应
		zap.L().Error("Sigup with invalid param", zap.Error(err)) // 如果有错误记录日志
		errs, ok := err.(validator.ValidationErrors)
		// 判断, 如果不是validation类型的错误就不需要翻译
		if !ok {
			ResponseError(c, CodeInvalidPassword)
			return
		}
		// 如果是, 则对错误进行翻译
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 手动对请求进行详细的业务逻辑校验判断   (如果使用validator库的话就不用手动验证了) 贼高级!!!!
	//if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.RePassword != p.Password{
	//	// 请求参数有误, 直接返回响应
	//	zap.L().Error("Sigup with invalid param") // 如果有错误记录日志
	//	c.JSON(http.StatusOK, gin.H{
	//		"msg": "请求参数有误",
	//	})
	//	return
	//}
	fmt.Println(p)
	// 2. 业务处理
	if err := logic.SignUp(p); err != nil{
		zap.L().Error("logic.SignUp faild", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist){
			ResponseError(c, CodeUserExist)
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应

	ResponseSuccess(c, nil)
}

// LoginHandler 处理登录所用的函数
func LoginHandler(c *gin.Context) {
	//  1. 获取请求参数及参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误, 直接返回响应
		zap.L().Error("Login with invalid param", zap.Error(err))
		// 判断err 是不是validator, ValidationErrors 类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))

		return
	}

	// 2. 业务逻辑代码
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}
	// 3. 返回响应

	ResponseSuccess(c, gin.H{
		"user_id":   fmt.Sprintf("%d", user.UserID), // id值大于1<<53-1  int64类型的最大值是1<<63-1
		"user_name": user.Username,
		"token":     user.Token,
	})
}