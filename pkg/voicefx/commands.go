package voicefx

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func makeHandleVoiceLeave(p Params) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		p.Log.Debug("voiceleave")
		vc, ok := s.VoiceConnections[i.GuildID]
		if !ok {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "i'm not in your channel",
				},
			})
			return
		}
		vc.Disconnect()
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "ok",
			},
		})
	}
}

func makeHandleVoiceJoinMe(p Params) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		g, err := s.State.Guild(i.GuildID)
		if err != nil {
			return
		}

		found := false
		for _, vs := range g.VoiceStates {
			if vs.UserID == i.Member.User.ID {
				cv, _ := s.ChannelVoiceJoin(i.GuildID, vs.ChannelID, true, false)
				cv.AddHandler(func(vc *discordgo.VoiceConnection, vs *discordgo.VoiceSpeakingUpdate) {
					p.Log.Debug("voiceSpeakingUpdate", zap.Reflect("gid", vc.GuildID), zap.Reflect("cid", vc.ChannelID), zap.Reflect("speaking", vs))
				})
				found = true
			}
		}

		if found {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("ok"),
				},
			})
			return
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("are you in a voice channel?"),
			},
		})
	}
}
