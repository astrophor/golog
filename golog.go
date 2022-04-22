//
// Package golog provides a rolling logger.
//
// golog plays well with any logging package that can write to an
// io.Writer, including the standard library's log package.
//
package golog

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	defaultLogPath        = "./"
	defaultMaxFileSize    = uint64(500) << 20
	defaultRotateInterval = 1
	defaultNameTimeFormat = "20060102150405"
)

// LogWriter is an io.WriteCloser that writes logs to rotate files.
type LogWriter struct {
	prefix         string
	path           string
	maxSize        uint64 // max size of each log file, in bytes
	rotateInterval int    // max time interval to rotate file, in time.Hour
	nameTimeFormat string

	logID        int
	logSize      uint64
	logFd        *os.File
	logStartTime time.Time

	mu sync.Mutex
}

// New returns an initialized LogWriter.
func New(logPath, logPrefix, logNameTimeFormat string, maxSize uint64, rotateInterval int) *LogWriter {
	return new(LogWriter).set(logPath, logPrefix, logNameTimeFormat, maxSize, rotateInterval)
}

// Write implements io.Writer.
// if the length of the writes is greater than maxSize, an error is returned
// if a write would cause the log file to be larger than maxSize, it will rotate to a new file
// if current time reaches rotate interval, it also will rotate to a new file
// if compress is set, old file will compress using gzip
func (l *LogWriter) Write(p []byte) (n int, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	writeLen := uint64(len(p))
	if writeLen > l.maxSize {
		return 0, fmt.Errorf("write length %v is max than max file size %v", writeLen, l.maxSize)
	}

	if l.logFd == nil {
		if err := l.open(); err != nil {
			return 0, fmt.Errorf("open log file failed: %v", err)
		}
	}

	if l.logSize+writeLen > l.maxSize || time.Since(l.logStartTime) > time.Duration(l.rotateInterval)*time.Hour {
		if err := l.rotate(); err != nil {
			return 0, fmt.Errorf("rotate to new file failed: %v", err)
		}
	}

	n, err = l.logFd.Write(p)
	l.logSize += uint64(n)

	return n, err
}

// Close implements io.Closer.
func (l *LogWriter) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.close()
}

func (l *LogWriter) set(logPath, logPrefix, logNameTimeFormat string, maxSize uint64, rotateInterval int) *LogWriter {
	l.prefix = logPrefix
	l.path = logPath
	l.maxSize = maxSize
	l.rotateInterval = rotateInterval
	l.nameTimeFormat = logNameTimeFormat
	l.logID = 0

	if l.prefix == "" {
		l.prefix = filepath.Base(os.Args[0])
	}

	if l.path == "" {
		l.path = defaultLogPath
	}

	if l.nameTimeFormat == "" {
		l.nameTimeFormat = defaultNameTimeFormat
	}

	if l.maxSize == 0 {
		l.maxSize = defaultMaxFileSize
	}

	if l.rotateInterval <= 0 {
		l.rotateInterval = defaultRotateInterval
	}

	return l
}

func (l *LogWriter) open() error {
	l.logSize = uint64(0)
	l.logStartTime = time.Now()
	fileName := fmt.Sprintf("%v/%v_%v_%v.log", l.path, l.prefix, l.logStartTime.Format(l.nameTimeFormat), l.logID)

	var err error
	if l.logFd, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(0644)); err != nil {
		return fmt.Errorf("open new log file: %v failed: %v", fileName, err)
	}

	return nil
}

func (l *LogWriter) rotate() error {
	if err := l.close(); err != nil {
		return err
	}

	l.logID++
	return l.open()
}

func (l *LogWriter) close() error {
	if l.logFd == nil {
		return nil
	}

	err := l.logFd.Close()
	l.logFd = nil

	return err
}
