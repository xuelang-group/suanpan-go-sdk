package stream

type Request struct {
	ID    string `mapstructure:"id"`
	Extra string `mapstructure:"extra"`
	Data  string
}
