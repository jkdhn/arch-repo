package desc

import (
	"encoding/hex"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"
)

type Encoder struct {
	out io.Writer
}

func NewEncoder(out io.Writer) *Encoder {
	return &Encoder{
		out: out,
	}
}

func (e *Encoder) encodeString(value string) (string, error) {
	if value == "" {
		return "", nil
	}
	return value + "\n", nil
}

func (e *Encoder) encodeSlice(value reflect.Value) (string, error) {
	if value.Type().Elem().Kind() == reflect.Uint8 {
		if value.Kind() == reflect.Array {
			value = value.Slice(0, value.Len())
		}
		s := hex.EncodeToString(value.Bytes())
		return s + "\n", nil
	}

	var result strings.Builder
	for i := 0; i < value.Len(); i++ {
		s, err := e.encodeValue(value.Index(i))
		if err != nil {
			return "", err
		}
		result.WriteString(s)
	}
	return result.String(), nil
}

func (e *Encoder) encodeUint64(value uint64) (string, error) {
	return fmt.Sprintf("%v\n", value), nil
}

func (e *Encoder) encodeTime(value time.Time) (string, error) {
	return fmt.Sprintf("%v\n", value.Unix()), nil
}

func (e *Encoder) encodeValue(value reflect.Value) (string, error) {
	switch value.Kind() {
	case reflect.String:
		return e.encodeString(value.String())
	case reflect.Slice, reflect.Array:
		return e.encodeSlice(value)
	case reflect.Uint64:
		return e.encodeUint64(value.Uint())
	case reflect.Struct:
		switch v := value.Interface().(type) {
		case time.Time:
			return e.encodeTime(v)
		}
	}

	return "", fmt.Errorf("unknown type: %v", value.Type())
}

func (e *Encoder) encode(field reflect.StructField, value reflect.Value) error {
	key, ok := field.Tag.Lookup("pkgdesc")
	if !ok {
		if value.Type().Kind() == reflect.Struct {
			return e.encodeStruct(value)
		}
		return nil
	}

	encoded, err := e.encodeValue(value)
	if err != nil {
		return err
	}

	if encoded == "" {
		return nil
	}

	if _, err := fmt.Fprint(e.out, "%", key, "%\n", encoded, "\n"); err != nil {
		return err
	}

	return nil
}

func (e *Encoder) encodeStruct(value reflect.Value) error {
	for i := 0; i < value.NumField(); i++ {
		if err := e.encode(value.Type().Field(i), value.Field(i)); err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) Encode(desc *Description) error {
	return e.encodeStruct(reflect.Indirect(reflect.ValueOf(desc)))
}
