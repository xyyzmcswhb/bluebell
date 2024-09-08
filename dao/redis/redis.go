package redis

import (
	"fmt"
	"web_app/settings"

	"github.com/go-redis/redis"
)

var (
	rdb *redis.Client
	//Nil = redis.Nil
)

func Init(cfg *settings.RedisConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			cfg.Host,
			cfg.Port),
		Password: cfg.Password, // 密码
		DB:       cfg.DB,       // 数据库
		PoolSize: cfg.PoolSize, // 连接池大小
	})

	//opt, err := redis.ParseURL("redis://<user>:<pass>@localhost:6379/<db>")
	//if err != nil {
	//	panic(err)
	//}
	//
	//rdb := redis.NewClient(opt)
	_, err = rdb.Ping().Result() //连接redis服务
	return err
}

func Close() {
	_ = rdb.Close()
}
