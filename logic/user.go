package logic

import (
	"web_app/dao/mysql"
	"web_app/models"
	"web_app/pkg/jwt"
	"web_app/pkg/snowflake"
)

// 存放业务逻辑的代码
func Signup(p *models.ParamSignUp) (err error) {
	//判断用户存不存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		//数据库查询出错
		return err
	}
	//if !exist {
	//	//用户已存在的错误
	//	return errors.New("用户已存在")
	//}
	//利用雪花算法生成UID
	userID := snowflake.GenID()
	//构造一个User实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	//密码加密
	//保存进数据库
	err = mysql.InsertUser(user)
	return err
}

func Login(p *models.ParamLogIn) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	//传递的是指针
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	//生成JWT tokon并返回
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return
	}
	user.Token = token
	return
}
