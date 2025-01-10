package enumutils

import (
	"fmt"
)

type EnumBase[T comparable] struct {
	Val2StrMap map[T]string
	Str2ValMap map[string]T
	EnumName   string
}

func NewEnumBase[T comparable](enumName string, val2StrMap map[T]string) EnumBase[T] {
	var enumBase EnumBase[T]

	enumBase.EnumName = enumName
	enumBase.Val2StrMap = val2StrMap
	enumBase.Str2ValMap = make(map[string]T)
	for k, v := range val2StrMap {
		enumBase.Str2ValMap[v] = k
	}
	return enumBase
}

func (e EnumBase[T]) String(val T) string {
	v, ok := e.Val2StrMap[val]
	if ok {
		return v
	}
	return "unkown " + e.EnumName
}

func (e EnumBase[T]) MarshalText(val T) ([]byte, error) {
	return []byte(e.String(val)), nil
}

func (e EnumBase[T]) UnmarshalText(data []byte) (T, error) {
	v, ok := e.Str2ValMap[string(data)]
	if ok {
		return v, nil
	}
	return v, fmt.Errorf("unknown %s: %v", e.EnumName, string(data))
}

func (e EnumBase[T]) Parse(str string) (T, error) {
	return e.UnmarshalText([]byte(str))
}
