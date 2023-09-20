package main

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(debug bool, level string) {
	logRotation := &lumberjack.Logger{
		Filename:   "./frontend.log", // Имя файла для записи логов
		MaxSize:    100,              // Максимальный размер файла логов в мегабайтах до ротации
		MaxBackups: 3,                // Максимальное количество ротированных файлов
		MaxAge:     28,               // Максимальное количество дней хранения ротированных файлов
		Compress:   true,             // Архивировать старые лог-файлы
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true, // Включатель цветных логов
		FullTimestamp: true, // Отобразить полное время записи лога
	})
	logrus.SetOutput(logRotation)

	p_level, err := logrus.ParseLevel(level)

	logrus.SetLevel(p_level)

	if err != nil {
		logrus.Error("Cannot parse log level from config!")
	}

	if debug {
		//logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("DEBUG MODE TURNED ON!")
	} else {
		//logrus.SetLevel(logrus.InfoLevel)
		logrus.Info("DEBUG MODE TURNED OFF!")
	}
}
