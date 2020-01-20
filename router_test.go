//
// Copyright (c) 2020 Matthew Penner <me@matthewp.io>
//
// This repository is licensed under the MIT License.
// https://github.com/matthewpi/router/blob/master/LICENSE.md
//

package router

import (
	"github.com/andersfylling/disgord"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	prefix = "."
	cmds   = &commands{}
)

type commands struct{}

func (c *commands) Yay(_ *disgord.MessageCreate) error {
	return nil
}

func (c *commands) Descriptions() map[string]string {
	return map[string]string{}
}

func (c *commands) Arguments() map[string][]string {
	return map[string][]string{}
}

func newRouter() (*Router, error) {
	return NewRouter(&disgord.Client{}, prefix, cmds)
}

func TestNewRouter(t *testing.T) {
	t.Run("AllArgs", func(t *testing.T) {
		a := assert.New(t)

		// Create a router with only passing a logger to prevent panics
		router, err := newRouter()
		a.NoError(err)
		a.NotNil(router)

		a.NotNil(router.Client)

		a.Equal(prefix, router.Prefix)
		a.NotNil(router.registrar)

		a.NotNil(router.Commands, 0)
	})

	t.Run("MissingClient", func(t *testing.T) {
		a := assert.New(t)

		router, err := NewRouter(nil, prefix, cmds)
		if a.Error(err) {
			a.Equal(ErrMissingClient, err)
		}
		a.Nil(router)
	})

	t.Run("InvalidPrefix", func(t *testing.T) {
		a := assert.New(t)

		router, err := NewRouter(&disgord.Client{}, "", cmds)
		if a.Error(err) {
			a.Equal(ErrInvalidPrefix, err)
		}
		a.Nil(router)
	})

	t.Run("MissingRegistrar", func(t *testing.T) {
		a := assert.New(t)

		router, err := NewRouter(&disgord.Client{}, prefix, nil)
		if a.Error(err) {
			a.Equal(ErrMissingRegistrar, err)
		}
		a.Nil(router)
	})
}

func TestRouter_GetCommandByName(t *testing.T) {
	t.Run("ValidCommand", func(t *testing.T) {
		a := assert.New(t)

		router, err := newRouter()
		a.NoError(err)
		a.NotNil(router)

		a.NotNil(router.GetCommandByName("yay"))
	})

	t.Run("MissingCommand", func(t *testing.T) {
		a := assert.New(t)

		router, err := newRouter()
		a.NoError(err)
		a.NotNil(router)

		a.Nil(router.GetCommandByName("not_a_registered_command"))
	})
}
