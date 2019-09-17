package config

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/imroc/req"
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

func SetAPILogger(req *req.Resp, resp *http.Response) {
	log.WithFields(log.Fields{
		"request_info": req,
	}).Info("Request information detail")

	log.WithFields(log.Fields{
		"response_info": resp,
	}).Info("Response information detail")
}
