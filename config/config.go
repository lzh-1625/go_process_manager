package config

var CF = &configuration{
	LogLevel: "debug",
}

// 只支持 float64、int、int64、bool、string类型
type configuration struct {
	LogLevel                  string  `default:"debug"  describe:"日志等级"`
	Listen                    string  `default:":8797" describe:"监听端口"`
	EsEnable                  bool    `default:"false" describe:"启用Elasticsearch"`
	EsUrl                     string  `default:"" describe:"Elasticsearch url"`
	EsIndex                   string  `default:"server_log_v1" describe:"Elasticsearch index"`
	EsUsername                string  `default:"" describe:"Elasticsearch用户名"`
	EsPassword                string  `default:"" describe:"Elasticsearch密码"`
	FileSizeLimit             float64 `default:"10.0" describe:"文件大小限制（MB）"`
	ProcessInputPrefix        string  `default:">" describe:"进程输入前缀"`
	ProcessRestartsLimit      int     `default:"2" describe:"进程重启次数限制"`
	ProcessMsgCacheLinesLimit int     `default:"50" describe:"std进程缓存消息行数"`
	ProcessMsgCacheBufLimit   int     `default:"4096" describe:"pty进程缓存消息字节长度"`
	ProcessExpireTime         int64   `default:"60" describe:"进程控制权过期时间（秒）"`
	PerformanceInfoListLength int     `default:"30" describe:"性能信息存储长度"`
	PerformanceInfoInterval   int     `default:"1" describe:"监控获取间隔时间（分钟）"`
	UserPassWordMinLength     int     `default:"4" describe:"用户密码最小长度"`
	LogMinLenth               int     `default:"0" describe:"过滤日志最小长度"`
	PprofEnable               bool    `default:"true" describe:"启用pprof分析工具"`
}
