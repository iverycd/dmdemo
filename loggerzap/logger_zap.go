package loggerzap

import (
	"dmdemo/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var (
	globalLogger *zap.Logger
)

// 初始化日志系统
func Init(env string, cfg *config.LogConfig) error {
	var core zapcore.Core
	// 设置编码器
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 不同级别的日志用不同颜色
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	//encoder := zapcore.NewJSONEncoder(encoderConfig) // json格式存储
	encoder := zapcore.NewConsoleEncoder(encoderConfig) // 终端输出转储到日志文件
	// 生产环境配置
	if env == "production" {
		// 日志文件切割配置
		lumberJackLogger := &lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
			LocalTime:  true,
		}

		// 生产环境使用文件+错误级别过滤
		core = zapcore.NewCore(
			encoder,
			zapcore.AddSync(lumberJackLogger),
			getZapLevel(cfg.Level),
		)

	} else {
		// 开发环境使用控制台+彩色输出
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewCore(
			consoleEncoder,
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
			zapcore.DebugLevel,
		)
	}
	// 创建Logger
	globalLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	// 替换zap全局Logger
	zap.ReplaceGlobals(globalLogger)
	return nil
}

// 获取日志级别
func getZapLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// 获取全局Logger
func L() *zap.Logger {
	return globalLogger
}

// 安全关闭
func Sync() error {
	return globalLogger.Sync()
}
