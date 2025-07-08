package config

type LogConfig struct {
	Level      string `yaml:"level"`       // 日志级别: debug, info, warn, error
	FilePath   string `yaml:"file_path"`   // 日志文件路径
	MaxSize    int    `yaml:"max_size"`    // 单个文件最大大小(MB)
	MaxBackups int    `yaml:"max_backups"` // 保留旧文件最大数量
	MaxAge     int    `yaml:"max_age"`     // 保留旧文件最大天数
	Compress   bool   `yaml:"compress"`    // 是否压缩旧文件
}
