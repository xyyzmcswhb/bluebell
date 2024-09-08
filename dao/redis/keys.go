package redis

//redis key

// redis key注意尽量使用命名空间的方式，区分不同的key，方便查询和拆分
const (
	KeyPrefix              = "web_app:"
	KeyPostTimeZset        = "post:time"   //发帖时间
	KeyPostScoreZset       = "post:score"  //投票分数
	KeyPostVotedZsetPrefix = "post:voted:" //记录用户及投票类型，参数是post_id
	KeyCommunitySetPrefix  = "commmunity:" //set保存每个分区下每个帖子的id
)

// 给redis KEY加上前缀
func getRedisKey(key string) string {
	return KeyPrefix + key
}
