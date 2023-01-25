package voicefx

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

func (r *recordingManager) stop(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	stop, ok := r.recordingStopTriggers[i.GuildID]
	if !ok {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "i'm not recording",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	<-stop

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "ok",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	return
}

func (r *recordingManager) record(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	if _, recording := r.recordingStopTriggers[i.GuildID]; recording {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "i'm already recording",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}

	stop := make(chan struct{})
	r.recordingStopTriggers[i.GuildID] = stop
	go r.handleVoice(vc.OpusRecv, stop)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "ok",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	return
}
