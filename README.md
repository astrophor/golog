#usage

## use standard log
```
package main

import (
	"log"

	"github.com/astrophor/golog"
)

func main() {
	golog.NewStdLog("./", "test", "", uint64(1)<<20, 1)
	log.Println("this is standard log")
}
```

## use zap
```
package main

import (
	"github.com/astrophor/golog"
	"go.uber.org/zap"
)

func main() {
	logger := golog.NewZapLog("./", "zap_test", "", uint64(1)<<20, 1)

	logger.Info("this is a info log",
		zap.String("place_id", "test"),
		zap.Int("value", 12))
}
```