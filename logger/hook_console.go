package logger

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type ConsoleHook struct {
	Formatter logrus.Formatter
}

func (h *ConsoleHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *ConsoleHook) Fire(entry *logrus.Entry) error {
	// 格式化日志
	line, err := h.Formatter.Format(entry)
	if err != nil {
		return err
	}

	str := string(line)
	fmt.Println(strings.ReplaceAll(str, "\n", ""))
	return nil
}
