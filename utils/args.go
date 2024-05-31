package utils

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type argData struct {
	kind  reflect.Kind
	value any
}

func (a *argData) parse(order []*argData, data string) error {
	switch a.kind {
	case reflect.String:
		if len(order) == 0 {
			a.value = data
			return nil
		}

		end := strings.Index(data, " ")
		if end == -1 {
			a.value = data
			return nil
		}

		a.value = data[:end]
		return order[0].parse(order[1:], data[end+1:])
	case reflect.Int:
		end := strings.Index(data, " ")
		if end == -1 {
			end = len(data)
		}

		var err error
		a.value, err = strconv.Atoi(data[:end])
		if err != nil {
			return err
		}

		if len(order) != 0 {
			return order[0].parse(order[1:], data[end+1:])
		}

		return nil
	}

	return errors.New("unknown kind")
}

type ArgsParser struct {
	args []*argData
}

func (a *ArgsParser) Parse(minCount int, message string) (*ArgsParser, error) {
	if len(a.args) == 0 {
		return a, nil
	}

	start := strings.Index(message, " ")

	if start == -1 {
		if minCount != 0 {
			return a, errors.New("too few arguments")
		}
		return a, nil
	}

	err := a.args[0].parse(a.args[1:], message[start+1:])
	if err != nil {
		return a, err
	}

	for i := 0; i < minCount; i++ {
		if a.args[i].value == nil {
			return a, errors.New("too few arguments")
		}
	}

	return a, nil
}

func (a *ArgsParser) AddInt() *ArgsParser {
	a.args = append(a.args, &argData{
		kind: reflect.Int,
	})

	return a
}

func (a *ArgsParser) AddString() *ArgsParser {
	a.args = append(a.args, &argData{
		kind: reflect.String,
	})

	return a
}

func (a *ArgsParser) Get(i int) any {
	return a.args[i].value
}

func NewArgsParser() *ArgsParser {
	return &ArgsParser{args: make([]*argData, 0)}
}
