package logging

import (
	"log"
	"os"
	"github.com/microsoft/BladeMonRT/configs"
)

/** Utility class used to create a logger. */
type LoggerFactory struct{}

func (loggerFactory LoggerFactory) ConstructLogger(typeName string) log.Logger {
	file, err := os.OpenFile(configs.LOGGING_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	var logger log.Logger = *log.New(file, typeName+" ", log.Ldate|log.Ltime|log.Lshortfile)
	return logger
}
