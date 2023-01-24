package ublockfx

import (
	"context"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/fx"

	"github.com/ajvpot/vanderbot/pkg/discordfx"
)

type NewCommandParams struct {
	fx.In
}
type NewCommandResult struct {
	fx.Out
	Command *discordfx.ApplicationCommandWithHandler `group:"commands"`
}

func ptr[T any](t T) *T {
	return &t
}

// NewCommand initializes a uBlock command.
func NewCommand(p NewCommandParams) NewCommandResult {
	return NewCommandResult{Command: &discordfx.ApplicationCommandWithHandler{
		Command: discordgo.ApplicationCommand{
			Name:        "ublock",
			Description: "Generate a report of resources blocked by uBlock Origin.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "URL",
					Description: "URL to check for blocked resources.",
					Required:    true,
				},
			},
		},
		Handler: func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Doing the thing",
				},
			})

			time.AfterFunc(time.Second, func() {
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: ptr(optionMap["URL"].StringValue()),
				})
			})
		},
	}}
}

var Module = fx.Options(fx.Provide(NewCommand))
