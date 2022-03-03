package stream

import "strconv"

type Request struct {
	ID    string `mapstructure:"id"`
	Extra string `mapstructure:"extra"`
	Input  map[string]string
}

const InputDataPrefix = `in`

func (r *Request) InputData(i int) string {
	return r.Input[InputDataPrefix + strconv.Itoa(i)]
}