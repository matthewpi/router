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

// Registrar represents a Command Registrar.
type Registrar interface {
	Descriptions() map[string]string
	Arguments() map[string][]string
}

var ignoredRegistrarMethods []string

func init() {
	ignoredRegistrarMethods = getIgnoredRegistrarMethods()
}

func getIgnoredRegistrarMethods() []string {
	typeIRegistrar := reflect.TypeOf((*Registrar)(nil)).Elem()

	var methods []string
	for i := 0; i < typeIRegistrar.NumMethod(); i++ {
		method := typeIRegistrar.Method(i)
		methods = append(methods, method.Name)
	}

	return methods
}

func isRegistrarMethodIgnored(name string) bool {
	for _, method := range ignoredRegistrarMethods {
		if method == name {
			return true
		}
	}

	return false
}

func getRegistrarMethods(i Registrar) ([]reflect.Value, []reflect.Method, error) {
	t := reflect.TypeOf(i)

	// Check if the interface is not a pointer.
	if t.Kind() != reflect.Ptr {
		return nil, nil, ErrInterfaceIsNotAPointer
	}

	values := make([]reflect.Value, 0)
	methods := make([]reflect.Method, 0)

	methodCount := t.NumMethod()
	if methodCount < 3 {
		// "i" has no commands, only has the required Registrar interface methods.
		return values, methods, nil
	}

	v := reflect.ValueOf(i)
	for i := 0; i < methodCount; i++ {
		value := v.Method(i)

		if !value.CanInterface() {
			continue
		}

		method := t.Method(i)

		// Check if the method should be ignored.
		if isRegistrarMethodIgnored(method.Name) {
			continue
		}

		// Check if the method has no arguments.
		// method.Type.In(0) will return the struct the method is on, not the first argument.
		if method.Type.NumIn() < 2 {
			continue
		}

		// Check if the first method argument is not *disgord.MessageCreate.
		if method.Type.In(1) != typeMessageCreate {
			continue
		}

		values = append(values, value)
		methods = append(methods, method)
	}

	return values, methods, nil
}
