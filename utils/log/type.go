package log

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"

	FormatText = "text"
	FormatJSON = "json"
)

type Config struct {
	Format string `json:"format" yaml:"format" `
	Level  string `json:"level" yaml:"level"`
}
