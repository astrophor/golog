package golog

import (
	"log"
	"testing"
)

func TestStdLog(t *testing.T) {
	NewStdLog("./", "test", "", uint64(1)<<20, 1)

    log.Println("test")

	log.Println("this is standard log")
}
