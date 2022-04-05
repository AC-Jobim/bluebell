package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
)

// 把每一步数据库操作封装成函数

const secret = "zhaobin"

// CheckUserExist 检查用户是否存在
func CheckUserExist(username string) (err error) {
	sqlStr := "select count(user_id) from user where username = ?"
	var count int
	if err = db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

// InsertUser 向数据库中插入新的用户
func InsertUser(user *models.User) error {
	sqlStr := "insert into user(user_id, username, password) values(?, ?, ?)"
	_, err := db.Exec(sqlStr, user.UserID, user.Username, encryptPassword(user.Password)) //对密码进行加密
	return err
}

func Login(user *models.User) (err error) {
	oPassword := user.Password // 用户登录的密码
	sqlStr := "select user_id, username, password from user where username = ?"
	err = db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	if err != nil {
		return err
	}
	// 判断密码是否正确
	if user.Password != encryptPassword(oPassword) {
		return ErrorInvalidPassword
	}
	return
}

func encryptPassword(password string) (result string) {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(password)))
}
