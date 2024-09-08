package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"web_app/models"
)

// 把每一步数据库操作封装成函数
// 待logic层根据业务需求调用
const secret = "hb&wyl"

// CheckUserExist检查指定用户名的用户是否存在
func CheckUserExist(username string) (err error) {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

// InsertUser向数据库中插入新的用户记录
func InsertUser(user *models.User) (err error) {
	//对密码进行加密
	user.Password = encryptPassword(user.Password)
	//执行sql语句入库
	sqlStr := `Insert into user(user_id,username,password) values(?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return
}

// 利用哈希算法给密码加密
func encryptPassword(opassword string) string {
	h := md5.New()
	h.Write([]byte(secret))                             //写入字节
	return hex.EncodeToString(h.Sum([]byte(opassword))) //将字节转换成 16进制字符串
}

// Login登陆，向数据库中根据用户名查询用户信息
func Login(user *models.User) (err error) {
	oPassword := user.Password //用户登录时所传入的密码
	//查询用户是否在数据库中
	sqlStr := `select user_id, username, password from user where username = ?`
	err = db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}

	if err != nil {
		//查询数据库失败
		return err
	}
	//判断密码是否正确
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrorPassword
	}
	return
}

// 根据ID(帖子作者id)获取用户信息
func GetUserByID(uid int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id,username 
			   from user where user_id = ?`
	err = db.Get(user, sqlStr, uid)
	return
}
