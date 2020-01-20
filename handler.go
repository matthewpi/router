//
// Copyright (c) 2020 Matthew Penner <me@matthewp.io>
//
// This repository is licensed under the MIT License.
// https://github.com/matthewpi/router/blob/master/LICENSE.md
//

package router

import (
	"github.com/andersfylling/disgord"
	"reflect"
	"strings"
	"unicode"
)

// Handle handles an incoming *disgord.MessageCreate event.
func (r *Router) Handle(e *disgord.MessageCreate) error {
	message := e.Message.Content
	label, argument := getLabelAndArgument(message)

	// Find the matching command using the label.
	command := r.GetCommandByName(label)
	if command == nil {
		return &ErrUnknownCommand{
			Command: label,
		}
	}

	// Get the []string of arguments.
	arguments := getArguments(argument)

	// Get the argument values for the reflection method call.
	argumentValues, err := getArgumentValues(r.Prefix, command, arguments)
	if err != nil {
		return err
	}

	// Call the command handler.
	if err := callWith(command.value, e, argumentValues...); err != nil {
		return &ErrCommandExecution{
			Command: command,
			err:     err,
		}
	}

	return nil
}

func getLabelAndArgument(message string) (string, string) {
	var label string
	var argument string
	for index, char := range message {
		// Check if this is the last character in the array.
		if index+1 == len(message) {
			label = strings.ToLower(message[1 : index+1])
			break
		}

		// Check if the character is a space.
		if unicode.IsSpace(char) {
			label = strings.ToLower(message[1:index])
			argument = message[index+1:]
			break
		}
	}

	return label, argument
}

func getArguments(argument string) []string {
	var arguments []string
	if strings.Contains(argument, " ") {
		arguments = strings.Split(argument, " ")
	} else if len(argument) > 0 {
		arguments = []string{argument}
	} else {
		arguments = []string{}
	}

	return arguments
}

func getArgumentValues(prefix string, command *Command, arguments []string) ([]reflect.Value, error) {
	commandArgumentsLength := len(command.arguments)
	if commandArgumentsLength < 1 {
		return []reflect.Value{}, nil
	}

	if !command.isValidArgumentLength(len(arguments)) {
		return nil, &ErrMissingArguments{
			Prefix:  prefix,
			Command: command.name,
			Usage:   command.usage,
		}
	}

	argumentValues := make([]reflect.Value, len(command.arguments))
	for i := 0; i < len(arguments); i++ {
		var input string
		rawArgument := command.rawArgumentsIndex == i
		if rawArgument {
			input = strings.Join(arguments[i:], " ")
		} else {
			input = arguments[i]
		}

		v, err := command.arguments[i](input)
		if err != nil {
			return nil, &ErrInvalidUsage{
				Prefix:     prefix,
				Command:    command.name,
				Usage:      command.usage,
				ArgumentID: i,
			}
		}

		argumentValues[i] = v

		if rawArgument {
			break
		}
	}

	return argumentValues, nil
}
