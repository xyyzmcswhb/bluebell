package mysql

import "errors"

var (
	ErrorUserExist    = errors.New("用户已存在")
	ErrorUserNotExist = errors.New("用户不存在")
	ErrorPassword     = errors.New("用户名或密码错误")
	ErrorInvalidID    = errors.New("无效的社区ID")
	ErrorInvalidPid   = errors.New("无效的帖子ID")
)
