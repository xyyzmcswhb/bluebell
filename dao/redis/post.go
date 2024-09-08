package redis

import (
	"strconv"
	"time"
	"web_app/models"

	"github.com/go-redis/redis"
)

func GetIDSfromKey(key string, page, size int64) ([]string, error) {
	//确定查询索引的起始位置
	start := (page - 1) * size
	end := start + size - 1
	//按分数从大到小查询指定数量的元素
	return rdb.ZRevRange(key, start, end).Result()
}

func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	//从redis获取id
	//根据用户中请求的携带的order参数选用相应的key
	key := getRedisKey(KeyPostTimeZset)
	if p.Order == models.OrderbyScore {
		key = getRedisKey(KeyPostScoreZset)
	}
	//确定查询的索引起始点
	return GetIDSfromKey(key, p.Page, p.Size)
}

// 根据idlist查询每篇帖子赞成票的数据
func GetPostVoteData(idlist []string) (data []int64, err error) {
	//for _, id := range idlist {
	//	key := getRedisKey(KeyPostVotedZsetPrefix + id)
	//	//统计每篇帖子的赞成票数量
	//	v := rdb.ZCount(key, "min", "1").Val()
	//	data = append(data, v)
	//}
	pipeline := rdb.Pipeline()
	for _, id := range idlist {
		key := getRedisKey(KeyPostVotedZsetPrefix + id)
		pipeline.ZCount(key, "1", "1") //查询投票类型为1即赞成票的数据
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// GetPostVoteNum 根据id查询每篇帖子的投赞成票的数据
func GetPostVoteNum(ids int64) (data int64, err error) {
	key := KeyPostVotedZsetPrefix + strconv.Itoa(int(ids))
	// 查找key中分数是1的元素数量 -> 统计每篇帖子的赞成票的数量
	data = rdb.ZCount(key, "1", "1").Val()
	return data, nil
}

// 按社区根据社区id查询帖子列表id
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	//使用zinterstore把分区的帖子set与帖子分数的zset生成一个新的zset，计算两个结果的交集然后存在新的key中
	//针对新的zset按之前的逻辑取数据
	//社区的key
	//根据用户中请求的携带的order参数选用相应的key
	orderKey := getRedisKey(KeyPostTimeZset)
	if p.Order == models.OrderbyScore { //按照分数请求
		orderKey = getRedisKey(KeyPostScoreZset)
	}
	//社区的key
	ckey := getRedisKey(KeyCommunitySetPrefix + strconv.Itoa(int(p.CommunityID)))
	//利用缓存key减少zinterstore执行的次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if rdb.Exists(key).Val() < 1 {
		//不存在，需要计算
		pipeline := rdb.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX", //将两个ZSet聚合的时候，求最大值
		}, ckey, orderKey)
		pipeline.Expire(key, 60*time.Second) //设置超时时间60s
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	//存在的话，返回根据key查询到的idlist
	return GetIDSfromKey(key, p.Page, p.Size)

}
