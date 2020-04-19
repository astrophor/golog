package golog

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const logTimeLayout = `2006-01-02 15:04:05`

// NewZapLog initialize zap log using LogWriter.
func NewZapLog(logPath, logPrefix, logNameTimeFormat string, maxSize uint64, rotateInterval int) *zap.Logger {
	logWriter := New(logPath, logPrefix, logNameTimeFormat, maxSize, rotateInterval)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "pos",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     logTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.InfoLevel)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(logWriter)),
		atomicLevel,
	)

	caller := zap.AddCaller()

	return zap.New(core, caller)
}

func logTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(logTimeLayout))
}
