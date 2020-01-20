//
// Copyright (c) 2020 Matthew Penner <me@matthewp.io>
//
// This repository is licensed under the MIT License.
// https://github.com/matthewpi/router/blob/master/LICENSE.md
//

package router

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrInvalidBool = errors.New("invalid bool [true/false]")
)

// Parseable represents a parseable argument.
type Parseable interface {
	Parse(arg string) error
}

// ManualParseable represents a manually parseable argument.
type ManualParseable interface {
	ParseContent(arg string) error
}

// Formatter represents an argument that can be formatted for use in a usage string.
type Formatter interface {
	Format(field string) string
}

// argumentValueFn represents an argument value function.
type argumentValueFn func(string) (reflect.Value, error)

// getArgumentValueFn returns an argument value function for the given type
// that handles the type conversion of a string argument.
func getArgumentValueFn(t reflect.Type) (argumentValueFn, error) {
	// IParseable
	if t.Implements(typeIParseable) {
		return parseableArgumentValue(t), nil
	}

	// IManualParseable
	if t.Implements(typeIManualParseable) {
		return manualParseableArgumentValue(t), nil
	}

	var fn argumentValueFn

	switch t.Kind() {
	case reflect.String:
		fn = stringArgumentValue()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fn = intArgumentValue(t)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fn = uintArgumentValue(t)

	case reflect.Float32, reflect.Float64:
		fn = floatArgumentValue(t)

	case reflect.Bool:
		fn = boolArgumentValue()
	}

	if fn == nil {
		return nil, errors.New("invalid type: " + t.String())
	}

	return fn, nil
}

func parseableArgumentValue(t reflect.Type) argumentValueFn {
	mt, ok := t.MethodByName("Parse")
	if !ok {
		panic("router: type IParseable does not implement Parse")
	}

	return func(input string) (reflect.Value, error) {
		v := reflect.New(t.Elem())

		ret := mt.Func.Call([]reflect.Value{
			v, reflect.ValueOf(input),
		})

		if err := errorReturns(ret); err != nil {
			return nilV, err
		}

		return v, nil
	}
}

func manualParseableArgumentValue(t reflect.Type) argumentValueFn {
	mt, ok := t.MethodByName("ParseContent")
	if !ok {
		panic("router: type IManualParseable does not implement ParseContent")
	}

	return func(input string) (reflect.Value, error) {
		v := reflect.New(t.Elem())

		ret := mt.Func.Call([]reflect.Value{
			v, reflect.ValueOf(input),
		})

		if err := errorReturns(ret); err != nil {
			return nilV, err
		}

		return v, nil
	}
}

func stringArgumentValue() argumentValueFn {
	return func(input string) (reflect.Value, error) {
		return reflect.ValueOf(input), nil
	}
}

func intArgumentValue(t reflect.Type) argumentValueFn {
	return func(input string) (reflect.Value, error) {
		i, err := strconv.ParseInt(input, 10, 64)
		return quickRet(i, err, t)
	}
}

func uintArgumentValue(t reflect.Type) argumentValueFn {
	return func(input string) (reflect.Value, error) {
		u, err := strconv.ParseUint(input, 10, 64)
		return quickRet(u, err, t)
	}
}

func floatArgumentValue(t reflect.Type) argumentValueFn {
	return func(input string) (reflect.Value, error) {
		f, err := strconv.ParseFloat(input, 64)
		return quickRet(f, err, t)
	}
}

func boolArgumentValue() argumentValueFn {
	return func(input string) (reflect.Value, error) {
		switch strings.ToLower(input) {
		case "true", "yes", "y", "1":
			return reflect.ValueOf(true), nil
		case "false", "no", "n", "0":
			return reflect.ValueOf(false), nil
		default:
			return nilV, ErrInvalidBool
		}
	}
}
