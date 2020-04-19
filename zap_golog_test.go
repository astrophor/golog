package golog

import (
	"testing"

	"go.uber.org/zap"
)

func TestZapLog(t *testing.T) {
	logger := NewZapLog("./", "zap_test", "", uint64(1)<<20, 1)

	logger.Info("this is a info log",
		zap.String("place_id", "test"),
		zap.Int("value", 12))
}
