package logservice

import (
	"log"
	"os"
)

// InforLogger logs all non-error output
var InforLogger *log.Logger

// ErrorLogger logs all error output
var ErrorLogger *log.Logger

// AttachLogFile points the loggers to a file to log to
func AttachLogFile(f *os.File) {
	InforLogger = log.New(f, "INFOR\t", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(f, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
}
