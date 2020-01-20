//
// Copyright (c) 2020 Matthew Penner <me@matthewp.io>
//
// This repository is licensed under the MIT License.
// https://github.com/matthewpi/router/blob/master/LICENSE.md
//

package router

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRouter_Handle(t *testing.T) {
	// Waiting on the patience to write the test for this.
}

func Test_getLabelAndArgument(t *testing.T) {
	t.Run("NoArgument", func(t *testing.T) {
		a := assert.New(t)

		label, argument := getLabelAndArgument(prefix + "label")

		a.Equal("label", label)
		a.Len(argument, 0)
	})

	t.Run("WithArgument", func(t *testing.T) {
		a := assert.New(t)

		label, argument := getLabelAndArgument(prefix + "label an_argument")

		a.Equal("label", label, "wrong label")
		a.Equal("an_argument", argument)
	})
}

func Test_getArguments(t *testing.T) {
	t.Run("NoArguments", func(t *testing.T) {
		a := assert.New(t)

		arguments := getArguments("")

		a.NotNil(arguments, "arguments array is nil")
		a.Len(arguments, 0, "wrong amount of arguments in array")
	})

	t.Run("SingleArgument", func(t *testing.T) {
		a := assert.New(t)

		arguments := getArguments("a_single_argument")

		a.NotNil(arguments, "arguments array is nil")
		a.Len(arguments, 1, "wrong amount of arguments in array")
		a.Equal(arguments[0], "a_single_argument", "first argument does not match")
	})

	t.Run("MultipleArgument", func(t *testing.T) {
		a := assert.New(t)

		arguments := getArguments("first_argument second_argument third_argument wow")

		a.NotNil(arguments, "arguments array is nil")
		a.Len(arguments, 4, "wrong amount of arguments in array")
		a.Equal(arguments[0], "first_argument", "first argument does not match")
		a.Equal(arguments[1], "second_argument", "second argument does not match")
		a.Equal(arguments[2], "third_argument", "third argument does not match")
		a.Equal(arguments[3], "wow", "fourth argument does not match")
	})
}

func Test_getArgumentValues(t *testing.T) {
	// This is going to be a pain in the ass :/
}
