package models

const (
	Orderbytime  = "time"
	OrderbyScore = "score"
)

// 定义请求体中传入的参数结构体
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_Password" binding:"required,eqfield=Password"`
}

// ParamLogIn登录请求参数
type ParamLogIn struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 投票数据
type ParamVoteData struct {
	PostID    string `json:"post_id" binding:"required"`               //帖子id
	Direction int8   `json:"direction,string" binding:"oneof= 1 0 -1"` //赞成or反对
}

// 获取帖子列表参数
type ParamPostList struct {
	CommunityID int64  `json:"community_id" form:"community_id"` //可以为空
	Page        int64  `form:"page" json:"page"`
	Size        int64  `form:"size" json:"size"`
	Order       string `form:"order" json:"order"`
}

//type ParamCommunityPostList struct {
//	*ParamPostList
//}
