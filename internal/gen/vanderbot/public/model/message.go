//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"encoding/json"
)

type Message struct {
	Blob      json.RawMessage
	CreatedAt string
	EditedAt  *string
	MessageID string
	ChannelID string
	GuildID   *string
	IsDelete  bool
}
