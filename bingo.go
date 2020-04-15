package main

import (
	"bingo/pkg/log"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	help    = flag.Bool("h", false, "help")
	debug   = flag.Bool("debug", false, "debug mode")
	verbose = flag.Bool("v", false, "verbose")
)

type Config struct {
	AppName string `json:"app_name,omitempty"`
}

func main() {
	flag.Parse()
	// flag.Usage = flag.Usage()

	if *help {
		flag.Usage()
		os.Exit(1)
	}
	if *debug {
		// TODO: make debug mode.
	}
	if *verbose {
		// TODO: set log to output verbose message.
	}

	viper.SetConfigName("conf/bingo")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.SetDefault("server.port", 8888)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("read config failed: %v", err)
	}
	var c Config
	viper.Unmarshal(&c)

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Notice("read config file changed: %s, op:%s", e.Name, e.Op)
	})

	fmt.Println(viper.Get("app_name"))
	fmt.Println(viper.Get("log_level"))
	fmt.Println("mysql ip: ", viper.Get("mysql.ip"))
	fmt.Println("mysql port: ", viper.Get("mysql.port"))
	fmt.Println("mysql user: ", viper.Get("mysql.user"))
	fmt.Println("mysql password: ", viper.Get("mysql.password"))
	fmt.Println("mysql database: ", viper.Get("mysql.database"))
	fmt.Println("redis ip: ", viper.Get("redis.ip"))
	fmt.Println("redis port: ", viper.Get("redis.port"))

	fmt.Println("bingo")
	fmt.Println("this is a project of testing framework that design for server testing")

	e := gin.Default()
	gin.SetMode(gin.DebugMode)
	e.GET("/x", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"messgae": "success",
		})
	})
	e.Run(":8888")
}
