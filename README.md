#usage

## use standard log
```
package golog

import (
	"log"
	"testing"
)

func main() {
	NewStdLog("./", "test", "", uint64(1)<<20, 1)

	log.Println("this is standard log")
}
```

## use zap
```
package golog

import (
	"testing"

	"go.uber.org/zap"
)

func main() {
	logger := NewZapLog("./", "zap_test", "", uint64(1)<<20, 1)

	logger.Info("this is a info log",
		zap.String("place_id", "test"),
		zap.Int("value", 12))
}
```