package remember

import (
	"crypto/sha256"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"net/http"
	"time"
)

func generateHash(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func rememberFile(s *discordgo.Session, i *discordgo.InteractionCreate) error {
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

	hash := generateHash(string(content))
	msg := fmt.Sprintf("[TODO] File processed! Hash: %s", hash)

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})

	return err
}

func rememberText(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	options := i.ApplicationCommandData().Options
	text := options[0].Options[0].StringValue()

	time.Sleep(10 * time.Second)

	hash := generateHash(text)
	msg := fmt.Sprintf("[TODO] Text processed! Hash: %s", hash)
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})

	return err
}

func RememberCLIHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
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
		rememberFile(s, i)
	case "text":
		rememberText(s, i)
	}

	return err
}
