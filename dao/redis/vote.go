package redis

import (
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	ScorePerVote     = 432 //每一票多少分

)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

// redis存储帖子信息
func CreatePost(postid, communityid int64) (err error) {
	pipeline := rdb.TxPipeline()
	//帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZset), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postid,
	})
	//帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZset), redis.Z{
		//Score:  float64(time.Now().Unix()),
		Score:  0,
		Member: postid,
	})
	//把帖子id加入到社区的set
	cKey := getRedisKey(KeyCommunitySetPrefix + strconv.Itoa(int(communityid)))
	pipeline.SAdd(cKey, postid)
	_, err = pipeline.Exec()
	return err
}

func VoteforPost(userid, postid string, value float64) error {
	//1.判断投票限制
	//去redis取帖子发布时间
	postTime := rdb.ZScore(getRedisKey(KeyPostTimeZset), postid).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}

	//2和3需要放进事务中操作

	//2.更新分数
	//先查当前用户给当前帖子的之前的投票记录
	ov := rdb.ZScore(getRedisKey(KeyPostVotedZsetPrefix+postid), userid).Val()
	var op float64
	if value == ov {
		return ErrVoteRepeated
	}
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - value) //计算两次投票的差值
	pipeline := rdb.TxPipeline()
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZset), op*diff*ScorePerVote, postid) //计算分数
	//if ErrVoteTimeExpire != nil {
	//	return err
	//}
	//3.记录用户为该帖子投票的数据
	if value == 0 {
		//取消投票
		pipeline.ZRem(getRedisKey(KeyPostVotedZsetPrefix+postid), userid)
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedZsetPrefix+postid), redis.Z{
			Score:  value, //投的是赞成票还是反对票
			Member: userid,
		})
	}
	_, err := pipeline.Exec()
	return err
}
