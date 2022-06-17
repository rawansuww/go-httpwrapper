package types

type Field struct {
	Key   string
	Value string
}

func NewField(key string, value string) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}
