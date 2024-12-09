package commands

import (
	"strings"
	"time"

	"github.com/X3NOOO/twinkleshine/commands/utils"

	"github.com/bwmarrin/discordgo"
)

const maxMessageLength = 2000

func sendChunked(s *discordgo.Session, i *discordgo.InteractionCreate, message string) error {
	var response []string
	var chunk string

	lines := strings.Split(message, "\n")
	for _, line := range lines {
		if len(chunk)+len(line)+1 > int(maxMessageLength*0.8) && chunk != "" {
			response = append(response, chunk)
			chunk = line
		} else {
			if chunk != "" {
				chunk += "\n"
			}
			chunk += line
		}
	}
	if chunk != "" {
		response = append(response, chunk)
	}

	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &response[0],
	})
	if err != nil {
		return err
	}

	for _, msg := range response[1:] {
		time.Sleep(1 * time.Second)
		_, err = s.ChannelMessageSend(i.ChannelID, msg)
		if err != nil {
			return err
		}
	}

	return nil
}

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
		return utils.SendErrorEmbed(msg, true, s, i)
	}

	return sendChunked(s, i, reply)
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
		return utils.SendErrorEmbed(msg, true, s, i)
	}

	return sendChunked(s, i, reply)
}
