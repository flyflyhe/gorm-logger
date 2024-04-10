package configs

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
	"sync"
)

var (
	app  Config //全局配置
	once sync.Once
)

type Config struct {
	Web             web
	Log             logConfig
	Aes             aes
	IdentifyAppCode string `mapstructure:"identify_app_code"`
	OssUrl          string `mapstructure:"oss_url"`
	WebAuthUrl      string `mapstructure:"web_auth_url"`
	Cron            bool
	Storage         struct {
		Upload   string
		Download string
	}
	Mysql MysqlConfig
	Redis RedisConfig
}

type web struct {
	Port      int    `mapstructure:"port"`
	Ip        string `mapstructure:"ip"`
	PortAdmin int    `mapstructure:"port_admin"`
	IpAdmin   string `mapstructure:"ip_admin"`
	Url       string
	Debug     bool
	RunMode   string
}

type logConfig struct {
	Info struct {
		Filename string
	}
	Error struct {
		Filename string
	}
}

type aes struct {
	Key string
	Iv  string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	Db       int
}

type MysqlConfig struct {
	Mysql         mysqlConfigInternal
	MysqlReadList []mysqlConfigInternal `mapstructure:"mysql_read_list"`
}

type mysqlConfigInternal struct {
	Host     string
	User     string
	Port     int
	Password string
	Database string
}

func (c *mysqlConfigInternal) GetDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		c.User, c.Password, c.Host, c.Port, c.Database,
	) + "&loc=Asia%2fShanghai"
}

func InitConfig(taxConfigFile string) {
	log.Printf(taxConfigFile)

	once.Do(func() {
		viper.SetConfigFile(taxConfigFile)

		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}

		if err := viper.Unmarshal(&app); err != nil {
			panic(err)
		}

		app.Web.Url = fmt.Sprintf("%s:%d", app.Web.Ip, app.Web.Port)
		log.Println(app.Web.Port)
		log.Println("key", app.Aes.Key, "iv", app.Aes.Iv)
	})
}

func InitConfigAuto(filename string) {
	dirs := []string{"./", "../", "../../", "../../../", "../../../../"}

	for _, dir := range dirs {
		if _, err := os.Stat(path.Join(dir, filename)); err != nil {
			continue
		} else {
			InitConfig(path.Join(dir, filename))
			break
		}
	}
}

func GetApp() Config {
	return app
}
