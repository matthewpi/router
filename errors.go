//
// Copyright (c) 2020 Matthew Penner <me@matthewp.io>
//
// This repository is licensed under the MIT License.
// https://github.com/matthewpi/router/blob/master/LICENSE.md
//

package router

// ErrUnknownCommand represents an Unknown Command error.
type ErrUnknownCommand struct {
	Command string
}

func (err *ErrUnknownCommand) Error() string {
	return "Unknown Command: `" + err.Command + "`"
}

// ErrMissingArguments represents a Missing Arguments error.
type ErrMissingArguments struct {
	Prefix  string
	Command string
	Usage   string
}

func (err *ErrMissingArguments) Error() string {
	return "Usage: `" + err.Prefix + err.Command + err.Usage + "`"
}

// ErrInvalidUsage represents an Invalid Usage error.
type ErrInvalidUsage struct {
	Prefix     string
	Command    string
	Usage      string
	ArgumentID int
}

func (err *ErrInvalidUsage) Error() string {
	// TODO: Format err.Usage with an emphasis around the invalid argument.
	return "Usage: `" + err.Prefix + err.Command + err.Usage + "`"
}

// ErrCommandExecution represents an unexpected error during a Command Execution.
type ErrCommandExecution struct {
	Command *Command
	err     error
}

func (err *ErrCommandExecution) Error() string {
	return "an unexpected error occurred while running that command. (error=" + err.err.Error() + ")"
}
