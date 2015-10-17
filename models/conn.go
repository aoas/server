package models

import (
	"fmt"
	"os"
	"path"

	// _ "github.com/denisenkom/go-mssqldb"
	"github.com/aoas/server/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	// _ "github.com/lib/pq"
	//_ "github.com/mattn/go-sqlite3"
)

var (
	x *xorm.Engine
)

func CreateEngine(cfg *config.Config) (engine *xorm.Engine, err error) {
	x, err = getEngine(cfg)
	if err != nil {
		return nil, fmt.Errorf("Fail to connect to database: %v", err)
	}

	if cfg.Database.LogPath != "" {
		logPath := path.Join(cfg.Database.LogPath, "xorm.log")
		os.MkdirAll(path.Dir(logPath), os.ModePerm)

		file, err := os.Create(logPath)
		if err != nil {
			return nil, fmt.Errorf("Fail to create xorm.log file: %v", err)
		}
		x.SetLogger(xorm.NewSimpleLogger(file))
	}

	x.ShowDebug = cfg.Database.ShowDebug
	x.ShowInfo = cfg.Database.ShowInfo
	x.ShowSQL = cfg.Database.ShowSQL
	x.ShowWarn = cfg.Database.ShowWarn
	x.ShowErr = cfg.Database.ShowErr

	// 如果有其他orm相关设置, 请直接设置即可

	return x, nil
}

func getEngine(cfg *config.Config) (*xorm.Engine, error) {
	connString := ""

	protocol := "tcp"
	if cfg.Database.Host[0] == '/' {
		protocol = "unix"
	}

	connString = fmt.Sprintf("%s:%s@%s(%s)/%s?charset=utf8&parseTime=true",
		cfg.Database.UserName, cfg.Database.Password, protocol, cfg.Database.Host, cfg.Database.DatabaseName)
	// switch cfg.Database.Type {
	// case "mysql":
	// 	protocol := "tcp"
	// 	if cfg.Database.Host[0] == '/' {
	// 		protocol = "unix"
	// 	}

	// 	connString = fmt.Sprintf("%s:%s@%s(%s)/%s?charset=utf-8&parseTime=true",
	// 		cfg.Database.UserName, cfg.Database.Password, protocol, cfg.Database.Host, cfg.Database.DatabaseName)

	// case "postgres":
	// 	var host, port = "127.0.0.1", "5432"
	// 	fields := strings.Split(cfg.Database.Host, ":")
	// 	if len(fields) > 0 && len(strings.TrimSpace(fields[0])) > 0 {
	// 		host = fields[0]
	// 	}
	// 	if len(fields) > 1 && len(strings.TrimSpace(fields[1])) > 0 {
	// 		port = fields[1]
	// 	}
	// 	connString = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
	// 		url.QueryEscape(cfg.Database.UserName), url.QueryEscape(cfg.Database.Password),
	// 		host, port, cfg.Database.DatabaseName, cfg.Database.SSLMode)
	// case "sqlite3":
	// 	if err := os.MkdirAll(path.Dir(cfg.Database.Path), os.ModePerm); err != nil {
	// 		return nil, fmt.Errorf("Fail to create directories: %v", err)
	// 	}
	// 	connString = "file:" + cfg.Database.Path + "?cache=shared&mode=rwc"
	// case "mssql":
	// 	var host, port = "127.0.0.1", "5432"
	// 	fields := strings.Split(cfg.Database.Host, ":")
	// 	if len(fields) > 0 && len(strings.TrimSpace(fields[0])) > 0 {
	// 		host = fields[0]
	// 	}
	// 	if len(fields) > 1 && len(strings.TrimSpace(fields[1])) > 0 {
	// 		port = fields[1]
	// 	}
	// 	connString = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s",
	// 		host, cfg.Database.UserName, cfg.Database.Password, port, cfg.Database.DatabaseName)
	// default:
	// 	return nil, fmt.Errorf("Invalid database type:　%s", cfg.Database.Type)
	// }

	//fmt.Printf("database connection string - %s", connString)

	return xorm.NewEngine(cfg.Database.Type, connString)

}
