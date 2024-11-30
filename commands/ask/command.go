package ask

import (
	"github.com/bwmarrin/discordgo"
)

func AskHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		return err
	}

	text := i.ApplicationCommandData().Options[0].StringValue()

	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: "[TODO] Reply: " + text,
	})

	return err
}
