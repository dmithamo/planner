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

// Initialize starts the logger
func (l *Log) Initialize() {
	var output *os.File = l.LogFile
	if output == nil {
		output = os.Stdout
		log.Println("loging to stdout")
	} else {
		log.Printf("logging to file %v", l.LogFile.Name())
	}

	l.InfoLogger = log.New(output, "INFO\t", log.Ldate|log.Ltime)
	l.ErrorLogger = log.New(output, "ERRO\t", log.Ldate|log.Ltime|log.Lshortfile)
}
