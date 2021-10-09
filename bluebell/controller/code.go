package controller

type ResCode int64
const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy

	// 中间件

	CodeInvalidToken
	CodeNeedLogin
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess: "success",
	CodeInvalidParam: "请求参数错误",
	CodeUserExist: "用户名已存在",
	CodeUserNotExist: "用户名不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy: "服务器繁忙",
	CodeNeedLogin : "需要登录",
	CodeInvalidToken: "无效的token",

}

func (c ResCode) Msg() string{
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}