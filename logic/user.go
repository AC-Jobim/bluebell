package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
)

func SignUp(p *models.ParamSignUp) error {
	// 1.检查用户是否存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		// 数据库查询出错或者用户名已存在
		return err
	}
	// 2.生成UID
	var user *models.User = &models.User{
		UserID:   snowflake.GenID(), //生成UID
		Username: p.Username,
		Password: p.Password,
	}
	// 3.保存进数据库
	return mysql.InsertUser(user)
}

func Login(p *models.ParamLogin) (token string, err error) {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	// 传递的是指针,可以拿到user.userID
	if err := mysql.Login(user); err != nil {
		// 登陆失败
		return "", err
	}
	return jwt.GenToken(user.UserID, user.Username)
}
