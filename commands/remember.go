package commands

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/X3NOOO/twinkleshine/ai"
)

func rememberFile(s *discordgo.Session, i *discordgo.InteractionCreate, ai ai.TwinkleshineAI) error {
	options := i.ApplicationCommandData().Options

	fileOption := options[0].Options[0]
	attachment := i.ApplicationCommandData().Resolved.Attachments[fileOption.Value.(string)]

	resp, err := http.Get(attachment.URL)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read file content: %v", err)
	}

	messageLink := fmt.Sprintf("https://discord.com/channels/%s/%s/%s", i.GuildID, i.ChannelID, i.ApplicationCommandData().TargetID)
	err = ai.RememberFile(content, map[string]any{
		"file": map[string]any{
			"name": attachment.Filename,
			"url":  messageLink,
		},
	})
	if err != nil {
		msg := fmt.Sprintf("Failed to remember file: %v", err)
		err = fmt.Errorf("failed to remember file: %v", err)

		_, derr := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &msg,
		})
		if derr != nil {
			return fmt.Errorf("failed to send error (%v) response: %v", err, derr)
		}
	}

	msg := "Done!"
	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})

	return err
}

func rememberText(s *discordgo.Session, i *discordgo.InteractionCreate, ai ai.TwinkleshineAI) error {
	options := i.ApplicationCommandData().Options
	text := options[0].Options[0].StringValue()

	err := ai.Remember(text, nil)
	if err != nil {
		msg := fmt.Sprintf("Failed to remember text: %v", err)
		err = fmt.Errorf("failed to remember text: %v", err)

		_, derr := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &msg,
		})
		if derr != nil {
			return fmt.Errorf("failed to send error (%v) response: %v", err, derr)
		}

		return err
	}

	msg := "Done!"
	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})

	return err
}

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
		err = rememberFile(s, i, c.AI)
	case "text":
		err = rememberText(s, i, c.AI)
	default:
		msg := "Unknown subcommand"
		err := fmt.Errorf("unknown subcommand: %v", subcommand)
		_, derr := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &msg,
		})
		if derr != nil {
			return fmt.Errorf("failed to send error (%v) response: %v", err, derr)
		}
	}

	return err
}

func (c *CommandContext) RememberGUIHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Uploading attachments...",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send initial response: %v", err)
	}

	attachments := i.ApplicationCommandData().Resolved.Messages[i.ApplicationCommandData().TargetID].Attachments
	if len(attachments) == 0 {
		msg := "No attachments found."
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &msg,
		})
		return err
	}

	errs := make(chan error, len(attachments))

	for _, att := range attachments {
		go func(att *discordgo.MessageAttachment) {
			log.Printf("Processing attachment: %s", att.Filename)

			resp, err := http.Get(att.URL)
			if err != nil {
				log.Printf("Failed to download %s: %v", att.Filename, err)
				errs <- fmt.Errorf("failed to download %s: %v", att.Filename, err)
				return
			}
			defer resp.Body.Close()

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Failed to read data %s: %v", att.Filename, err)
				errs <- fmt.Errorf("failed to read data %s: %v", att.Filename, err)
				return
			}

			messageLink := fmt.Sprintf("https://discord.com/channels/%s/%s/%s", i.GuildID, i.ChannelID, i.ApplicationCommandData().TargetID)
			err = c.AI.RememberFile(data, map[string]any{
				"file": map[string]any{
					"name": att.Filename,
					"url":  messageLink,
				},
			})
			if err != nil {
				log.Printf("Failed to remember file %s: %v", att.Filename, err)
				errs <- fmt.Errorf("failed to remember file %s: %v", att.Filename, err)
				return
			}

			log.Printf("Successfully remembered file: %s", att.Filename)

			errs <- nil
		}(att)
	}

	var failedUploads []string

	for i := 0; i < len(attachments); i++ {
		if err := <-errs; err != nil {
			failedUploads = append(failedUploads, err.Error())
		}
	}

	if len(failedUploads) > 0 {
		msg := "Some attachments failed to upload:\n" + strings.Join(failedUploads, "\n")
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &msg,
		})
		return err
	}

	msg := "Done!"
	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
	if err != nil {
		return fmt.Errorf("failed to edit response: %v", err)
	}

	return nil
}
