package log

import (
	"testing"
)

func TestLogA(t *testing.T) {
	//logger, _ := zap.NewDevelopment()
	//Default.logger = logger

	Default.Info("hei this is Debug")

	Default.Info("hei", "ff", "mm")
	Default.Debug("hei", "ff", "mm")

	la := New("Test")
	la.Info("im the king of the world!!")
	la.Warn("aa")
	la.Error(nil, "aa")
	la.Info("aa", 11)

	// defer logger.Sync()
	// logger.Info("Hello Zap!")
	// logger.Warn("Beware of getting Zapped! (Pun)")
	//logger.Error("I'm out of Zap joke!!")
}
