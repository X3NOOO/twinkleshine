package commands

import (
	"github.com/bwmarrin/discordgo"
)

func (c *CommandContext) AskCLIHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		return err
	}

	text := i.ApplicationCommandData().Options[0].StringValue()

	reply, err := c.AI.Query(text)
	if err != nil {
		return err
	}

	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: reply,
	})

	return err
}

func (c *CommandContext) AskGUIHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		return err
	}

	text := i.ApplicationCommandData().Resolved.Messages[i.ApplicationCommandData().TargetID].Content

	reply, err := c.AI.Query(text)
	if err != nil {
		return err
	}

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &reply,
	})

	return err
}
