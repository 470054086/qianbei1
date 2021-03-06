package config

type Config struct {
	App   app   `yaml:"app" json:"app"`
	Log   log   `yaml:"log" json:"log"`
	Mysql mysql `yaml:"mysql" json:"mysql"`
	Redis redis `json:"redis" yaml:"json"`
}

type app struct {
	// 域名
	Domain        string `json:"domain" yaml:"domain"`
	Port          string `json:"port" yaml:"port"`
	DbType        string `json:"dbtype" yaml:"dbtype"`
	Schedule      string `json:"schedule" yaml:"schedule"`
	CronRpcPort   string `json:"cronrpcport" yaml:"cronrpcport"`
	CronRpcHost   string `json:"cronrpchost" yaml:"cronrpchost"`
	DaemonRpcHost string `json:"daemonrpchost" yaml:"daemonrpchost"`
	DaemonRpcPort string `json:"daemonrpcport" yaml:"daemonrpcport"`
}
type log struct {
	// 日志地址
	Dir string `json:"dir" yaml:"dir"`
	// 日志格式
	Format string `json:"format" yaml:"format"`
}

type mysql struct {
	Username     string ` json:"username" yaml:"username"`
	Password     string ` json:"password" yaml:"password"`
	Path         string ` json:"path" yaml:"path"`
	Dbname       string `json:"dbname" yaml:"dbname"`
	Config       string ` json:"config" yaml:"config"`
	MaxIdleConns int    ` json:"maxIdleConns" yaml:"maxidleconns"`
	MaxOpenConns int    ` json:"maxOpenConns" yaml:"maxopenconns"`
	LogMode      bool   `json:"logMode" yaml:"logmode"`
}

type redis struct {
	Addr     string `json:"addr" yaml:"addr"`
	Password string `json:"password" yaml:"password"`
	Db       int    `json:"db" yaml:"db"`
}
