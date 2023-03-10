package voicefx

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func (p *voiceManager) leaveme(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	vc, ok := s.VoiceConnections[i.GuildID]
	if !ok {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "i'm not in your channel",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}
	vc.Disconnect()
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "ok",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

}

func (p *voiceManager) joinme(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	g, err := s.State.Guild(i.GuildID)
	if err != nil {
		return
	}

	userChannel := findChannel(g, i)
	if userChannel == "" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("are you in a voice channel?"),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	cv, _ := s.ChannelVoiceJoin(i.GuildID, userChannel, true, false)
	cv.AddHandler(func(vc *discordgo.VoiceConnection, vs *discordgo.VoiceSpeakingUpdate) {
		p.Log.Debug("voiceSpeakingUpdate", zap.Reflect("gid", vc.GuildID), zap.Reflect("cid", vc.ChannelID), zap.Reflect("speaking", vs))
	})

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("ok"),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	return
}

func findChannel(g *discordgo.Guild, i *discordgo.InteractionCreate) string {
	for _, vs := range g.VoiceStates {
		if vs.UserID == i.Member.User.ID {
			return vs.ChannelID
		}
	}
	return ""
}
