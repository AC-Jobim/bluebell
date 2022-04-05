package models

// ParamSignUp 注册请求参数
type ParamSignUp struct {
	//binding,进行参数校验：require就是需要这个,必须要有值,否则会报错
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

// ParamLogin 登录请求参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
