package stream

type Request struct {
	ID    string      `mapstructure:"id"`
	Extra interface{} `mapstructure:"extra"`
	Data  interface{}
}
