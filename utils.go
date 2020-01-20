//
// Copyright (c) 2020 Matthew Penner <me@matthewp.io>
//
// This repository is licensed under the MIT License.
// https://github.com/matthewpi/router/blob/master/LICENSE.md
//

package router

import (
	"reflect"
)

// callWith calls a caller using the specified arguments.
func callWith(caller reflect.Value, ev interface{}, values ...reflect.Value) error {
	return errorReturns(
		caller.Call(
			append(
				[]reflect.Value{reflect.ValueOf(ev)},
				values...,
			),
		),
	)
}

// errorReturns handles the reflection of fetching an error from a method call's return values.
func errorReturns(returns []reflect.Value) error {
	// WARNING: This assumes that the first return value is an error!
	v := returns[0].Interface()

	if v == nil {
		return nil
	}

	return v.(error)
}

// quickRet does something majestic.
func quickRet(v interface{}, err error, t reflect.Type) (reflect.Value, error) {
	if err != nil {
		return nilV, err
	}

	rv := reflect.ValueOf(v)

	if t == nil {
		return rv, nil
	}

	return rv.Convert(t), nil
}
