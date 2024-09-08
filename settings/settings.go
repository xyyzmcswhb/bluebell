package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"

	//"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 全局变量原用来保存程序的所有配置信息
var Conf = new(multipleConfig)

type multipleConfig struct {
	*AppConfig   `mapstructure:"app"`
	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type AppConfig struct {
	Name      string `mapstructure:"name"` //tag中标签表示工作模式
	Mode      string `mapstructure:"mode"`
	Version   string `mapstructure:"version"`
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
	Port      string `mapstructure:"port"`

	//*LogConfig   `mapstructure:"log"`
	//*MySQLConfig `mapstructure:"mysql"`
	//*RedisConfig `mapstructure:"Redis"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"` //tag中标签表示工作模式
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
}

type MySQLConfig struct {
	Host           string `mapstructure:"host"` //tag中标签表示工作模式
	Port           int    `mapstructure:"port"`
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	Dbname         string `mapstructure:"dbname"`
	MaxConnections int    `mapstructure:"max_connections"`
	MaxIdlecolumns int    `mapstructure:"max_idle_col"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"` //tag中标签表示工作模式
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

func Init(file string) (err error) {
	//viper.SetConfigName("config") // 指定配置文件路径
	//viper.SetConfigType("yaml")   // 指定配置文件路径
	//viper.AddConfigPath(".")
	viper.SetConfigFile(file)
	err = viper.ReadInConfig() // 读取配置信息
	if err != nil {            // 读取配置信息失败
		fmt.Printf("Fatal error config file: %s \n", err)
		return
	}

	if err := viper.Unmarshal(Conf); err != nil { //将配置信息反序列化至结构体变量中
		fmt.Printf("viper.Unmarshal failed, err:#{err}\n")
	}
	// 监控配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了")
		if err := viper.Unmarshal(Conf); err != nil { //将配置信息反序列化至结构体变量中
			fmt.Printf("viper.Unmarshal failed, err:#{err}\n")
		}
	})

	return
}
