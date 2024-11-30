package remember

import (
	"errors"
	"time"

	"github.com/bwmarrin/discordgo"
)

func RememberGUIHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Uploading...",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		return errors.New("failed to send initial response: " + err.Error())
	}

	time.Sleep(10 * time.Second)

	msg := "[TODO] Uploaded successfully!"

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
	if err != nil {
		return errors.New("failed to edit response: " + err.Error())
	}

	return nil
}
