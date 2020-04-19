package golog

import "log"

// NewStdLog initialize standard log using LogWriter.
func NewStdLog(logPath, logPrefix, logNameTimeFormat string, maxSize uint64, rotateInterval int) {
	logWriter := New(logPath, logPrefix, logNameTimeFormat, maxSize, rotateInterval)

	log.SetOutput(logWriter)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
