//
// Copyright (c) 2020 Matthew Penner <me@matthewp.io>
//
// This repository is licensed under the MIT License.
// https://github.com/matthewpi/router/blob/master/LICENSE.md
//

// Package args ...
package args // import "go.matthewp.io/router/args"

import (
	"errors"
	"regexp"
)

// getFirstResult gets the first regexp result.
func getFirstResult(reg *regexp.Regexp, itemName string, input string, output *string) error {
	matches := reg.FindStringSubmatch(input)
	if len(matches) < 2 {
		return errors.New("router: invalid '" + itemName + "'")
	}

	*output = matches[1]
	return nil
}
