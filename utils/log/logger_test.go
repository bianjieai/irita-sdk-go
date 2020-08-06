package log

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	log1 := NewLogger("info")

	log1.Info().Str("foo", "bar").Msg("Hello World")
	log1.Info().Str("foo1", "bar").Msg("Hello World")
	log1.Info().Str("foo2", "bar").Msg("Hello World")
	log1.Info().Str("foo3", "bar").Msg("Hello World")

}
