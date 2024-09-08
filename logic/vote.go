package logic

import (
	"strconv"
	"web_app/dao/redis"
	"web_app/models"

	"go.uber.org/zap"
)

//投票功能：
//1. 用户投票的数据确定下来

//投一票就加432分 86400/200->需要两百张赞成票才能给帖子续一天
/*
direction=1时，有两种情况：
	1.之前没有投过票，现在改投赞成票 -->更新分数和投票记录 差值的绝对值：1 +432
	2.之前投反对票，现在改投赞成票  -->更新分数和投票记录 差值的绝对值：2 +432*2
direction = 0,：
	1.之前投赞成票，现在取消投票  差值的绝对值：1 -432
	2.之前投过反对票，现在取消投票 差值的绝对值：1 +432
direciton=-1
	1.之前没有投过票，现在投反对票 差值的绝对值：1 -432
	2.之前投赞成票，现在改投反对票 差值的绝对值：1 -432*2

投票限制：
每个帖子发表之日起，一周内允许用户投票，超过一周不允许用户投票
	1.到期后将redis中保存的赞成票和反对票中存储到mysql
	2。到期后删除KeyPostVotedZsetPrefix
*/
//为帖子投票
func VoteforPost(userId int64, p *models.ParamVoteData) error {
	//redis.KeyPostTimeZset
	zap.L().Debug("vote for post",
		zap.Int64("userId", userId),
		zap.String("postid", p.PostID),
		zap.Int8("direction", p.Direction))
	userid := strconv.Itoa(int(userId))
	return redis.VoteforPost(userid, p.PostID, float64(p.Direction))

}
