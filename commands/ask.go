package commands

import (
	"fmt"

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
		msg := "Failed to process the text: " + err.Error()
		err = fmt.Errorf("failed to process the text: %v", err)
		_, derr := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: msg,
		})
		if derr != nil {
			return fmt.Errorf("failed to send error (%v) response: %v", err, derr)
		}

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
		msg := "Failed to process the text: " + err.Error()
		err = fmt.Errorf("failed to process the text: %v", err)
		_, derr := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: msg,
		})
		if derr != nil {
			return fmt.Errorf("failed to send error (%v) response: %v", err, derr)
		}

		return err
	}

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &reply,
	})

	return err
}
