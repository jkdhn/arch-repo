package info

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Decoder struct {
	scanner *bufio.Scanner
}

func NewDecoder(in io.Reader) *Decoder {
	return &Decoder{
		scanner: bufio.NewScanner(in),
	}
}

func (d *Decoder) decodeString(value reflect.Value, input string) error {
	value.SetString(input)
	return nil
}

func (d *Decoder) decodeUint64(value reflect.Value, input string) error {
	n, err := strconv.ParseUint(input, 10, 64)
	if err != nil {
		return err
	}

	value.SetUint(n)
	return nil
}

func (d *Decoder) decodeTime(value reflect.Value, input string) error {
	n, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return err
	}

	value.Set(reflect.ValueOf(time.Unix(n, 0)))
	return nil
}

func (d *Decoder) decodeSlice(value reflect.Value, input string) error {
	value.Set(reflect.Append(value, reflect.ValueOf(input)))
	return nil
}

func (d *Decoder) decode(value reflect.Value, input string) error {
	switch value.Kind() {
	case reflect.String:
		return d.decodeString(value, input)
	case reflect.Uint64:
		return d.decodeUint64(value, input)
	case reflect.Slice:
		return d.decodeSlice(value, input)
	case reflect.Struct:
		switch value.Interface().(type) {
		case time.Time:
			return d.decodeTime(value, input)
		}
	}

	return fmt.Errorf("unknown type: %v", value.Type())
}

func (d *Decoder) Decode() (*Info, error) {
	result := new(Info)
	t := reflect.TypeOf(result).Elem()
	v := reflect.ValueOf(result).Elem()
	fields := make(map[string]int)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		key, ok := field.Tag.Lookup("pkginfo")
		if !ok {
			continue
		}
		fields[key] = i
	}

	for d.scanner.Scan() {
		line := strings.SplitN(d.scanner.Text(), "=", 2)
		if line[0] == "" || line[0][0] == '#' {
			continue
		}

		key := strings.TrimSpace(line[0])
		value := strings.TrimSpace(line[1])
		i, ok := fields[key]
		if !ok {
			return nil, fmt.Errorf("unknown key: \"%v\"", key)
		}

		if err := d.decode(v.Field(i), value); err != nil {
			return nil, err
		}
	}
	if err := d.scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
