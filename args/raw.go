//
// Copyright (c) 2020 Matthew Penner <me@matthewp.io>
//
// This repository is licensed under the MIT License.
// https://github.com/matthewpi/router/blob/master/LICENSE.md
//

package args

// RawArguments represents a Raw Arguments argument.
type RawArguments struct {
	Argument string
}

func (r *RawArguments) ParseContent(args string) error {
	r.Argument = args
	return nil
}

func (r *RawArguments) Format(field string) string {
	return "[" + field + ": string...]"
}

func (r *RawArguments) String() string {
	return r.Argument
}
