package logkit

type LogLevel int

const (
	TRACE LogLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
)

func (l LogLevel) String() string {
	return []string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR"}[l]
}