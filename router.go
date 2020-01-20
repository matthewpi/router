//
// Copyright (c) 2020 Matthew Penner <me@matthewp.io>
//
// This repository is licensed under the MIT License.
// https://github.com/matthewpi/router/blob/master/LICENSE.md
//

// Package router ...
package router // import "go.matthewp.io/router"

import (
	"errors"
	"fmt"
	"github.com/andersfylling/disgord"
	"reflect"
	"strings"
)

var (
	// ErrMissingClient .
	ErrMissingClient = errors.New("router: missing client")
	// ErrInvalidPrefix .
	ErrInvalidPrefix = errors.New("router: invalid prefix")
	// ErrMissingRegistrar .
	ErrMissingRegistrar = errors.New("router: missing registrar")
	// ErrInterfaceIsNotAPointer represents an Interface Is Not A Pointer error.
	ErrInterfaceIsNotAPointer = errors.New("router: interface is not a pointer")
	// ErrCommandIsNil represents a Command Is Nil error.
	ErrCommandIsNil = errors.New("router: command is nil")
	// ErrMethodHasNoErrorReturn .
	ErrMethodHasNoErrorReturn = errors.New("router: method has no error return")
	// ErrMethodHasNoArguments .
	ErrMethodHasNoArguments = errors.New("router: method has no arguments")
	// ErrMissingMessageCreateArgument .
	ErrMissingMessageCreateArgument = errors.New("router: missing *disgord.MessageCreate as the first method argument")

	// nilV is used to represent a nil reflect#Value.
	nilV = reflect.Value{}

	typeMessageCreate    = reflect.TypeOf((*disgord.MessageCreate)(nil))
	typeIError           = reflect.TypeOf((*error)(nil)).Elem()
	typeIParseable       = reflect.TypeOf((*Parseable)(nil)).Elem()
	typeIManualParseable = reflect.TypeOf((*ManualParseable)(nil)).Elem()
	typeIFormatter       = reflect.TypeOf((*Formatter)(nil)).Elem()
)

// Router .
type Router struct {
	*disgord.Client

	Prefix    string
	registrar Registrar

	Commands []*Command
}

// NewRouter .
func NewRouter(client *disgord.Client, prefix string, i Registrar) (*Router, error) {
	if client == nil {
		return nil, ErrMissingClient
	}

	if len(prefix) != 1 {
		return nil, ErrInvalidPrefix
	}

	if i == nil {
		return nil, ErrMissingRegistrar
	}

	r := &Router{
		Client: client,

		Prefix:    prefix,
		registrar: i,
	}

	if err := r.registerCommands(); err != nil {
		return nil, err
	}
	return r, nil
}

// GetCommandByName attempts to get a *Command by matching it's name.
func (r *Router) GetCommandByName(name string) *Command {
	for _, c := range r.Commands {
		if c.name == name {
			return c
		}
	}

	return nil
}

// registerCommands registers the commands on the registrar.
func (r *Router) registerCommands() error {
	values, methods, err := getRegistrarMethods(r.registrar)
	if err != nil {
		return err
	}

	commands := make([]*Command, 0)
	for i := 0; i < len(values); i++ {
		command, err := r.getCommand(values[i], methods[i])
		if err != nil {
			return err
		}

		if command == nil {
			return ErrCommandIsNil
		}

		commands = append(commands, command)
	}

	r.Commands = commands
	return nil
}

func (r *Router) getCommand(value reflect.Value, method reflect.Method) (*Command, error) {
	// Check if the method does not return anything to prevent a panic.
	if value.Type().NumOut() < 1 {
		return nil, ErrMethodHasNoErrorReturn
	}

	// Check if the first return value is not an error
	if err := value.Type().Out(0); err == nil || !err.Implements(typeIError) {
		return nil, ErrMethodHasNoErrorReturn
	}

	// The first argument will always be the struct value, so we ignore it
	args := method.Type.NumIn()
	if args < 2 {
		return nil, ErrMethodHasNoArguments
	}

	// Check first argument type
	if method.Type.In(1) != typeMessageCreate {
		return nil, ErrMissingMessageCreateArgument
	}

	// Create a new command
	command := &Command{
		name: strings.ToLower(method.Name),

		value:  value,
		method: method,

		arguments: make([]argumentValueFn, 0, args),

		rawArgumentsIndex: -1,
	}
	command.Description = r.registrar.Descriptions()[command.name]

	// Handle method arguments
	if args > 2 {
		methodArgs, ok := r.registrar.Arguments()[command.name]
		if !ok {
			return nil, fmt.Errorf("router: %s takes arguments and does not have a usage", method.Name)
		}

		if args-2 != len(methodArgs) {
			return nil, fmt.Errorf("router: %s's usage does not have all the arguments present", method.Name)
		}

		var usageBuilder strings.Builder

		for i := 2; i < args; i++ {
			t := method.Type.In(i)

			if t.Implements(typeIManualParseable) {
				command.rawArgumentsIndex = i - 2
			}

			argValue, err := getArgumentValueFn(t)
			if err != nil {
				return nil, fmt.Errorf("router: error parsing argument %s: %v", t.String(), err)
			}

			var usage string
			if t.Implements(typeIFormatter) {
				mt, ok := t.MethodByName("Format")
				if !ok {
					panic("router: type IFormatter does not implement Format")
				}

				v := reflect.New(t.Elem())

				ret := mt.Func.Call([]reflect.Value{
					v, reflect.ValueOf(methodArgs[i-2]),
				})

				usage = ret[0].Interface().(string)
			} else {
				usage = "<" + methodArgs[i-2] + ": " + t.String() + ">"
			}

			command.arguments = append(command.arguments, argValue)
			usageBuilder.WriteString(" " + usage)
		}

		command.usage = usageBuilder.String()
	}

	return command, nil
}
