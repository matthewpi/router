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

// The only problem with this test is that the "m" array will need to be
// updated whenever a method is added to the Registrar interface.
func Test_getIgnoredRegistrarMethods(t *testing.T) {
	m := []string{
		"Descriptions",
		"Arguments",
	}

	a := assert.New(t)

	methods := getIgnoredRegistrarMethods()
	a.Equal(len(m), len(methods))

	// We cannot just check if the arrays are equal because if they are not in the same order the test will fail.
	for _, method := range methods {
		ok := false
		for _, method2 := range m {
			if method == method2 {
				ok = true
			}
		}

		a.True(ok)
	}
}
