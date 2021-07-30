package logging

import (
	"os"
    "log"
)

/** Utility class used to create a logger. */
type LoggerFactory struct {}

func (loggerFactory LoggerFactory) ConstructLogger(typeName string) *log.Logger {
    const (
        logging_file = "log"
    )

	file, err := os.OpenFile(logging_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatal(err)
    }
    var logger *log.Logger = log.New(file, typeName+" ", log.Ldate|log.Ltime|log.Lshortfile)
    return logger
}