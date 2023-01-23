package ublockfx

import (
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
			Description: "check page with ublock",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "url",
					Description: "url",
					Required:    true,
				},
			},
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
					Content: "Hey there! Congratulations, you just executed your first slash command",
				},
			})
			time.AfterFunc(time.Second, func() {
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: ptr(optionMap["url"].StringValue()),
				})
			})
		},
	}}
}

var Module = fx.Options(fx.Provide(NewCommand))
