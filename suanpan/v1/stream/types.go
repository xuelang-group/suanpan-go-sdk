package stream

type Request struct {
	ID    string `mapstructure:"id"`
	Extra string `mapstructure:"extra"`
	Input map[string]string
}

const (
	InputDataPrefix = `in`
	OutputDataPrefix = `out`
)
