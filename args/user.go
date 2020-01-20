//
// Copyright (c) 2020 Matthew Penner <me@matthewp.io>
//
// This repository is licensed under the MIT License.
// https://github.com/matthewpi/router/blob/master/LICENSE.md
//

package args

import (
	"github.com/andersfylling/disgord"
	"regexp"
)

// userRegex is regex for getting a user mention.
var userRegex = regexp.MustCompile(`<@!?(\d+)>`)

// UserMention represents a User Mention argument.
type UserMention string

func (m *UserMention) Parse(arg string) error {
	return getFirstResult(userRegex, "user mention", arg, (*string)(m))
}

func (m *UserMention) Format(field string) string {
	return "<" + field + ": @user>"
}

// Snowflake returns a disgord.Snowflake for the User Mention.
func (m *UserMention) Snowflake() disgord.Snowflake {
	return disgord.ParseSnowflakeString(string(*m))
}
