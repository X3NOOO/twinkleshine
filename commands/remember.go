package commands

import (
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/X3NOOO/twinkleshine/commands/remember"
)

func (c *CommandContext) RememberCLIHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to respond to interaction: %v", err)
	}

	options := i.ApplicationCommandData().Options
	subcommand := options[0].Name

	switch subcommand {
	case "file":
		remember.RememberFile(s, i)
	case "text":
		remember.RememberText(s, i)
	}

	return err
}

func (c *CommandContext) RememberGUIHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
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
