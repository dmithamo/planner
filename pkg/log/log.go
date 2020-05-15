package log

import (
	"log"
	"os"
)

// Log exposes the functionality available from this pkg
type Log struct {
	LogFile     *os.File
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
}

// AttachLogFile links this pkg's loggers to the logfile provided.
// Defaults to stdout, stderr
func (l *Log) AttachLogFile() {
	var infoLogOutput *os.File = l.LogFile
	if infoLogOutput == nil {
		infoLogOutput = os.Stdout
		log.Println("Loging info to stdout")
	} else {
		log.Printf("logging info to file %v", &l.LogFile)
	}
	l.InfoLogger = log.New(infoLogOutput, "INFO \t", log.Ldate|log.Ltime)

	var errLogOutput *os.File = l.LogFile
	if errLogOutput == nil {
		errLogOutput = os.Stdout
		log.Println("Loging errs to stderr")
	} else {
		log.Printf("logging errs to file %v", &l.LogFile)
	}
	l.ErrorLogger = log.New(errLogOutput, "ERROR\t", log.Ldate|log.Ltime)
}
