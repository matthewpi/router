# Router

[![GoDoc: Reference](https://img.shields.io/badge/godoc-reference-blue?style=flat-square)](https://godoc.org/go.matthewp.io/router)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square    )](LICENSE.md)

## Features

- Automatic command routing using function names
- Dynamic command arguments
- Argument validation
- Help message generator

## Usage

```go
package main

import (
	"fmt"
	"github.com/andersfylling/disgord"
	"go.matthewp.io/router"
	"os"
)

func main() {
	// Get the bot token using an environment variable.
	// DO NOT PREFIX THE TOKEN WITH "Bot ", disgord ALREADY DOES!
	var token = os.Getenv("BOT_TOKEN")
	if token == "" {
		panic("missing $BOT_TOKEN")
		return
	}

	client := disgord.New(disgord.Config{
		BotToken: token,
		Logger:   disgord.DefaultLogger(true),
	})

	r, err := router.NewRouter(client, ".", &commands{s: client})
	if err != nil {
		panic(err)
		return
	}

	client.On(disgord.EvtMessageCreate, func(s disgord.Session, e *disgord.MessageCreate) {
		err := r.Handle(e)
		if err != nil {
			_, err := s.SendMsg(e.Ctx, e.Message.ChannelID, "<@"+e.Message.Author.ID.String()+">, "+err.Error())
			if err != nil {
				fmt.Printf("failed to send error message to user: %v\n", err)
				return
			}
		}
	})
}

type commands struct {
	s disgord.Session
}

// Any exported functions on the struct will be registered as commands
// unless they are apart of the Registrar interface or if they do not
// accept a *disgord.MessageCreate as the first argument.
func (c *commands) Help(_ *disgord.MessageCreate) error {
	return nil
}

// This method will not be registered as command, because reflection does not support
// unexported methods when you are outside of the package it is defined in.
func (c *commands) test(_ *disgord.MessageCreate) error {
	return nil
}

// This method will also not be registered as a command but is still exported,
// however it will be registered as a command if it's first argument is
// *disgord.MessageCreate
func (c *commands) ThisIsNotACommand() {
	// Do something!
}

// Registrar interface method, it will be ignored.
func (c *commands) Descriptions() map[string]string {
	return map[string]string{
		"help": "Prints this help message",
	}
}

// Registrar interface method, it will be ignored.
func (c *commands) Arguments() map[string][]string {
	return map[string][]string{}
}
```

## Additional Information

### Interfaces

#### Parseable
Allows a custom argument type to get the string argument to handle it's own parsing.

```go
type Parseable interface {
	Parse(string) error
}
```
###### Example (refer to [`args/user.go`](args/user.go))


#### ManualParseable
Allows a custom argument type to get the entire argument string after any preceding arguments,
useful for getting long user inputs.

```go
type ManualParseable interface {
	ParseContent(string) error
}
```
###### Example (refer to [`args/raw.go`](args/raw.go))


#### Formatter
Allows a custom argument type to have a custom usage string.

```go
type Formatter interface {
	Format(string) string
}
```
###### Example (refer to [`args/raw.go`](args/raw.go) or [`args/user.go`](args/user.go))
