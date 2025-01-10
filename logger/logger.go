package logger

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/mcuadros/go-defaults"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type LogConfig struct {
	Level  logrus.Level `yaml:"level" default:"4"`
	Folder string       `yaml:"folder" default:".log"`
	Caller bool         `yaml:"caller" default:"false"`
	Format string       `yaml:"format" default:"text"` // 可选 text, json
}

var logConfig LogConfig

func Init(config ...LogConfig) {
	defaults.SetDefaults(&logConfig)
	if len(config) > 0 {
		logConfig = config[0]
	}
	fmt.Printf("logConfig: %+v\n", logConfig)
	logrus.SetLevel(logConfig.Level)
	logrus.SetReportCaller(logConfig.Caller)
	logrus.SetFormatter(formatter())
	logrus.SetOutput(output())
	logrus.AddHook(&ConsoleHook{consoleFormatter()})
	logrus.Info("init logrus done")

	refreshOutputDaily()
}

func refreshOutputDaily() {
	c := cron.New()

	_, err := c.AddFunc("@daily", func() {
		refreshOutput()
	})
	if err != nil {
		panic(err)
	}

	c.Start()
}

func refreshOutput() {
	logrus.SetOutput(output())
}

func output() io.Writer {
	dirPath := logConfig.Folder
	if dirPath == "" {
		logrus.Panic("not set log folder path")
	}
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0777)
		if err != nil {
			logrus.Panic(err)
		}
	}

	filePath := path.Join(dirPath, time.Now().Format("2006_01_02.log"))
	writer, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		logrus.Panic(err)
	}
	return writer
}

func formatter() logrus.Formatter {
	if logConfig.Format == "json" {
		return &logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "ztimestamp",
				logrus.FieldKeyLevel: "zlevel",
				logrus.FieldKeyMsg:   "@message",
				logrus.FieldKeyFunc:  "zcaller",
				logrus.FieldKeyFile:  "zfile",
			},
			CallerPrettyfier: callerPrettyfier,
		}
	}
	return &logrus.TextFormatter{
		CallerPrettyfier: callerPrettyfier,
	}
}

func consoleFormatter() logrus.Formatter {
	return &logrus.TextFormatter{
		CallerPrettyfier: callerPrettyfier,
		ForceColors:      true, // 强制使用颜色
		// FullTimestamp:    true, // 显示完整时间戳
	}
}

func callerPrettyfier(f *runtime.Frame) (function string, file string) {
	s := strings.Split(f.Function, ".")
	funcname := s[len(s)-1]
	_, filename := path.Split(f.File)
	return funcname, filename
}
