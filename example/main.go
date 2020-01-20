//
// Copyright (c) 2020 Matthew Penner <me@matthewp.io>
//
// This repository is licensed under the MIT License.
// https://github.com/matthewpi/router/blob/master/LICENSE.md
//

package main

import (
	"context"
	"github.com/andersfylling/disgord"
	"go.matthewp.io/router"
	"log"
	"os"
)

var r *router.Router

func main() {
	var token = os.Getenv("BOT_TOKEN")
	if token == "" {
		panic("missing $BOT_TOKEN")
		return
	}

	client := disgord.New(disgord.Config{
		BotToken: token,
		Logger:   disgord.DefaultLogger(true),
	})

	var err error
	r, err = router.NewRouter(client, ".", &commands{s: client})
	if err != nil {
		log.Panicf("failed to create a new router: %v", err)
		return
	}

	client.On(disgord.EvtReady, clientReady)
	client.On(disgord.EvtMessageCreate, clientMessageCreate)

	_ = client.StayConnectedUntilInterrupted(context.Background())
}

func clientReady(s disgord.Session, _ *disgord.Ready) {
	err := s.UpdateStatus(&disgord.UpdateStatusPayload{
		Game: &disgord.Activity{
			Name: "Testing stuff :)",
		},
		Status: "dnd",
	})
	if err != nil {
		s.Logger().Error("failed to update status: " + err.Error())
	}
}

func clientMessageCreate(s disgord.Session, e *disgord.MessageCreate) {
	// I doubt this would happen but, lets keep it :)
	if e == nil {
		s.Logger().Debug("EvtMessageCreate: e == nil")
		return
	}

	// Ignore messages sent by bots.
	if e.Message.Author.Bot {
		s.Logger().Debug("EvtMessageCreate: author is a bot")
		return
	}

	// Ignore embedded messages.
	if len(e.Message.Embeds) > 0 {
		s.Logger().Debug("EvtMessageCreate: message is embedded")
		return
	}

	message := e.Message.Content
	if string(message[0]) != r.Prefix {
		s.Logger().Debug("EvtMessageCreate: message does not start with the prefix")
		return
	}

	s.Logger().Debug("EvtMessageCreate: Received '" + message + "'")

	// Command Router
	err := r.Handle(e)
	if err != nil {
		s.Logger().Debug("error while executing command: " + err.Error())

		_, err := s.SendMsg(e.Ctx, e.Message.ChannelID, "<@"+e.Message.Author.ID.String()+">, "+err.Error())
		if err != nil {
			s.Logger().Error("failed to send message to user: " + err.Error())
		}
		return
	}
	s.Logger().Debug("executed command successfully")
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

func (c *commands) ThisIsNotACommand2(_ interface{}) {
	// Do something!
}

func (c *commands) Descriptions() map[string]string {
	return map[string]string{
		"help": "Prints this help message",
	}
}

func (c *commands) Arguments() map[string][]string {
	return map[string][]string{}
}
