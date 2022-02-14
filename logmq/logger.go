package logmq

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"runtime"
)

type (
	writeHook struct { //отправлять любое количество уровней логирования в любой писатель
		Writer   []io.Writer
		LogLevel []logrus.Level
	}
)

var (
	line string
	err  error
)

/*type Logger struct {
	*logrus.Entry
}
var entrLog *logrus.Entry

func GetLogger() Logger {
	return Logger{entrLog}
}

func (l *Logger) GetLoggerWithFild(k string, v interface{}) Logger {
	return Logger{l.WithField(k, v)}

}*/

func (hook *writeHook) Fire(entry *logrus.Entry) error {
	line, err = entry.String()
	if err != nil {
		fmt.Printf("func Fire Strings error: %v", err) //???
		return err
	}

	for _, w := range hook.Writer { //перебераем входные писатели
		_, err = w.Write([]byte(line))
		if err != nil {
			fmt.Printf("Writer error: %v", err)
			return err
		}
	}
	return nil
}

func (hook *writeHook) Levels() []logrus.Level { // возвращает уровни из хука
	return hook.LogLevel
}

func InitializationLogger() *logrus.Entry {

	var (
		file   *os.File
		logger = logrus.New()
	)

	logger.SetReportCaller(true)

	logger.Formatter = &logrus.JSONFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) { //внутри происходит логирование
			filename := path.Base(frame.File)                                                    // информация о файле
			return fmt.Sprintf("%s", frame.Function), fmt.Sprintf("%s %d", filename, frame.Line) //название файла, номер строчки
		},
	}

	file, err = os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {

		logger.Errorf("error opening the log file: %v", err)
		return nil
	}

	logger.SetOutput(io.Discard)

	logger.AddHook(&writeHook{
		Writer:   []io.Writer{file, os.Stdout},
		LogLevel: logrus.AllLevels,
	})

	logger.SetLevel(logrus.TraceLevel)

	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	return logrus.NewEntry(logger)
}

