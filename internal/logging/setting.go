package logging

import (
	"MediaTools/constants"
	"MediaTools/utils"
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type serviceLoggerSetting struct{}

func (s *serviceLoggerSetting) Format(entry *logrus.Entry) ([]byte, error) {
	// 根据日志级别设置颜色
	var colorCode uint8
	switch entry.Level {
	case logrus.DebugLevel:
		colorCode = constants.ColorBlue
	case logrus.InfoLevel:
		colorCode = constants.ColorGreen
	case logrus.WarnLevel:
		colorCode = constants.ColorYellow
	case logrus.ErrorLevel:
		colorCode = constants.ColorRed
	default:
		colorCode = constants.ColorGray
	}

	// 设置文本Buffer
	var b *bytes.Buffer
	if entry.Buffer == nil {
		b = &bytes.Buffer{}
	} else {
		b = entry.Buffer
	}
	// 时间格式化
	formatTime := entry.Time.Format("2006-01-02 15:04:05")

	fmt.Fprintf(
		b,
		"\033[3%dm【%s】\033[0m %s | %s - %s\n", // 长度需要算是上控制字符的长度
		colorCode,
		strings.ToUpper(entry.Level.String()),
		formatTime,
		entry.Caller.Function,
		entry.Message,
	)
	return b.Bytes(), nil
}

func (s *serviceLoggerSetting) Levels() []logrus.Level {
	return []logrus.Level{logrus.ErrorLevel, logrus.WarnLevel, logrus.InfoLevel}
}

// HOOK
//
// 将日志写入文件
func (s *serviceLoggerSetting) Fire(entry *logrus.Entry) error {
	const logFileDir = "logs"
	if err := os.MkdirAll(logFileDir, os.ModePerm); err != nil {
		return err
	}
	logFilePath := logFileDir + "/" + time.Now().Local().Format("2006-01-02") + ".log"
	serviceLogFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer serviceLogFile.Close()

	line, err := entry.String()
	if err != nil {
		return err
	}
	serviceLogFile.WriteString(utils.RemoveColorCodes(line))
	return nil
}
