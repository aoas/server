package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/aoas/server/config"
	"github.com/aoas/server/models"
	"github.com/aoas/server/routers"
	"github.com/aoas/server/utils"
)

func main() {
	var conf config.Config

	if _, err := toml.DecodeFile("./config.toml", &conf); err != nil {
		log.Fatalf("invalid config file - %s", err.Error())
	}

	fmt.Println("origins :", conf.CORS.Origins)
	// 创建数据库ORM
	engine, err := models.CreateEngine(&conf)
	if err != nil {
		log.Fatalf("create database engine failed - %s", err.Error())
	}
	// 同步相间结构到数据库
	if err := models.SyncTables(); err != nil {
		log.Fatalf("sync model to database failed - %s", err.Error())
	}

	// 初始化日志
	logFile, err := os.Create(conf.LogPath)
	defer logFile.Close()
	if err != nil {
		log.Fatalf("Create log file failed - %s", err.Error())
	}
	logger := utils.NewSimpleLogger(logFile)

	// 缓存用户及权限相关表
	if conf.Database.EnableCache {
		models.EnableCache()
	}

	router, err := routers.New(engine, conf, logger)
	if err != nil {
		log.Fatalf("init router failed - %s", err.Error())
	}
	// cors

	router.Run(fmt.Sprintf("%s:%d", conf.App.Host, conf.App.Port))

}
