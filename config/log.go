package config

import (
	"fmt"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func LogFormatter() {
	log.SetFormatter(&log.TextFormatter{})
	filename := "./logfile.log"

	_, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Failed to open file : ", err)
	}

	debug, _ := strconv.ParseBool(os.Getenv("IS_LOG_FILE"))
	if debug {
		log.SetOutput(&lumberjack.Logger{
			Filename:   filename,
			MaxSize:    1,
			MaxBackups: 3,
			MaxAge:     28,
		})
	}
}
