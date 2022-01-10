package logger

// 日志库的设计
import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
	"web1/settings"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Init 初始化Logger,将配置文件中的LogConfig传进来
func Init(cfg *settings.LogConfig, mode string) (err error) {
	//定制logger，将logger写入文件中
	//1、创建编码器Encoder
	encoder := zapEncoder()
	//2、日志写入器（日志写到哪里）
	writeSync := zapLogWrite(
		cfg.Filename,
		cfg.MaxSize,
		cfg.MaxBackups,
		cfg.MaxAge,
	)
	//3、定义日志级别zapcore.Level
	l := new(zapcore.Level)
	//log.level是string类型，序列化成zapcore.Level类型
	// 这样设计的好处是，可以直接在配置文件中修改那种级别作为日志的级别（cfg.Level）
	err = l.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return
	}
	// 4.0、方便改变模式，主要用来改变发布模式和开发模式
	var core zapcore.Core
	if mode == "dev" {
		// 进入开发模式，日志输出到终端
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		// 设置同时向多方面写日志信息，包括文件，终端等
		core = zapcore.NewTee(
			// 把日志文件写入日志文件中，以下面的发布模式一样
			zapcore.NewCore(encoder, writeSync, l),
			// 将日志信息向终端输出，使用os.Stdout标准输出
			// 这里需要将os.Stdout转换成zapcore.WriteSyncer格式，使用zapcore.Lock进行转化
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
	} else {
		// 进入发布模式，将日志信息记录到文件中
		//4、创建定制的logger
		core = zapcore.NewCore(encoder, writeSync, l)
	}
	lg := zap.New(core, zap.AddCaller())
	//5、 替换zap库中全局的logger，当需要使用logger的时候调用zap.L()即可
	zap.ReplaceGlobals(lg)
	return
}

// 编码器
func zapEncoder() zapcore.Encoder {
	//定制自己的日志输出格式，不使用默认的日志输出格式
	//简单点：就是在默认的日志输出格式基础上修改一下
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	//将NewJSONEncoder修改成NewConsoleEncoder，就可以将JSON Encoder更改为普通的Log Encoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

// 日志写入文件
func zapLogWrite(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	//1、以什么格式写到那个文件中
	//利用第三方库lumberjack进行日志分割归档
	//由于lumberJack需要四个参数，所以函数应该接收四个参数
	lumberJackLogger := lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	//2、写入文件
	return zapcore.AddSync(&lumberJackLogger)
}

// 在routes中作为中间件使用
// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()             // 时间
		path := c.Request.URL.Path      // 请求的路径
		query := c.Request.URL.RawQuery // 请求的参数
		c.Next()                        // 执行后续的中间件

		cost := time.Since(start) //计时（运行的总时长），结合前面的start
		zap.L().Info(path,
			zap.Int("status", c.Writer.Status()),                                 // 请求的状态
			zap.String("method", c.Request.Method),                               // 请求的方法
			zap.String("path", path),                                             // 请求的路径
			zap.String("query", query),                                           // 请求的参数
			zap.String("ip", c.ClientIP()),                                       // 请求的ip
			zap.String("user-agent", c.Request.UserAgent()),                      // 调用者
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()), // 错误信息
			zap.Duration("cost", cost),                                           // 总耗时
		)
	}
}

// 在routes中作为中间件使用
// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					zap.L().Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
