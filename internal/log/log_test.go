package log

import (
	"testing"

	"github.com/lack-io/cirrus/config"
)

func TestPrint(t *testing.T) {
	Init(nil)

	Debug("debug...")

	Info("info...")

	Warn("warn...")

	Error("error...")

	Fatal("fatal...")
}

func TestInit(t *testing.T) {
	cfg := config.Logger{
		Filename:   "test.log",
		MaxSize:    10,
		MaxAge:     1,
		MaxBackups: 1,
		LocalTime:  true,
		Compress:   false,
	}
	if err := Init(&cfg); err != nil {
		t.Fatal(err)
	}
	defer Sync()

	Debug("debug...")

	Info("info...")

	Warn("warn...")

	Error("error...")

	Fatal("fatal...")
}
