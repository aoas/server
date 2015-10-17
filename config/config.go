package config

// 数据库相关配置
type Database struct {
	Type         string // mysql/postgres/sqlite3/mssql
	Host         string
	UserName     string
	Password     string
	DatabaseName string
	Path         string // 如Sqlite等嵌入式数据库时存放的位置
	SSLMode      bool

	LogPath string

	ShowSQL   bool
	ShowInfo  bool
	ShowDebug bool
	ShowErr   bool
	ShowWarn  bool

	EnableCache bool
}

type App struct {
	Host string
	Port int
}

type Email struct {
	Host     string
	Port     int
	UserName string
	Password string
	// 发件邮箱
	Sender string
}

type File struct {
	UploadPath string
}

// 配置主结构
type Config struct {
	// Token加密密钥 (AES)
	TokenSecret    string
	TokenExpiredIn int64

	LogPath  string
	Database Database
	App      App
	File     File
}
