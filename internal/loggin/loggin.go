package loggin

import (
	"log"
	"os"
)

var (
	Warning *log.Logger
	Info    *log.Logger
	Error   *log.Logger
	Debug   *log.Logger
)

func init() {
	log_path, exists := os.LookupEnv("LOG_FILE")

	if !exists {
		log_path = "sfp.log"
	}

	file, err := os.OpenFile(log_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		log.Fatal(err)
	}

	log_args := log.Ldate | log.Ltime | log.Lshortfile

	Info = log.New(file, "INFO: ", log_args)
	Warning = log.New(file, "WARNING: ", log_args)
	Error = log.New(file, "ERROR: ", log_args)
	Debug = log.New(file, "DEBUG: ", log_args)
}
