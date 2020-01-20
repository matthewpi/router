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

type Command struct {
	name        string
	Description string

	value  reflect.Value
	method reflect.Method

	arguments []argumentValueFn
	usage     string

	rawArgumentsIndex int
}

// Name returns the command's name.
func (c *Command) Name() string {
	return c.name
}

// Usage returns the command's usage.
func (c *Command) Usage() string {
	return c.usage
}

func (c *Command) isValidArgumentLength(length int) bool {
	// The Raw Arguments Index allows us to receive multiple spaced arguments as
	// one argument,  meaning that you cannot just directly check if the length
	// of arguments from the command and the length of the message's arguments match.
	// c.rawArgumentsIndex == -1 means there are no raw arguments in the command signature.
	if c.rawArgumentsIndex == -1 {
		if length != len(c.arguments) {
			return false
		}
	} else if length < len(c.arguments) {
		return false
	}

	return true
}
