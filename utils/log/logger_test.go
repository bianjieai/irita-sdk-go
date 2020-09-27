package log

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	log1 := NewLogger(Config{
		Format: "json",
		Level:  "info",
	})

	log1.Info("Hello World", "foo", "bar")
	log1.Info("Hello World", "foo1", "bar")
	log1.Info("Hello World", "foo2", "bar")
	log1.Info("Hello World", "foo3", "bar")
}
