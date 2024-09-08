package controller

import "web_app/models"

//专门用来放接口文档专门用来的model,接口文档返回的数据是一致的，具体data类型不一致

type _ResponsePostList struct {
	Code    ResCode                 `json:"code"`    //业务响应状态码
	Message string                  `json:"message"` //提示信息
	Data    []*models.ApiPostDetail `json:"data"`    //数据
}
