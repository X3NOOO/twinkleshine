package utils

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func SendErrorEmbed(msg string, edit bool, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	log.Printf("Sending error embed to %s [%s]: %s\n", i.Interaction.Member.User.Username, i.Interaction.Member.User.ID, msg)

	if edit {
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "An error occurred",
					Description: msg,
					Color:       0xcc0033,
				},
			},
		})
		if err != nil {
			err = fmt.Errorf("cannot edit interaction response: %v", err)
			log.Println(err)

			return err
		}
	} else {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "An error occurred",
						Description: msg,
						Color:       0xcc0033,
					},
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			err = fmt.Errorf("cannot respond to interaction: %v", err)
			log.Println(err)

			return err
		}
	}

	return nil
}

func SendSuccessEmbed(msg string, edit bool, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	log.Printf("Sending success embed to %s [%s]: %s\n", i.Interaction.Member.User.Username, i.Interaction.Member.User.ID, msg)

	if edit {
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "Success!",
					Description: msg,
					Color:       0x00cc11,
				},
			},
		})
		if err != nil {
			err = fmt.Errorf("cannot edit interaction response: %v", err)
			log.Println(err)

			return err
		}
	} else {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Success!",
						Description: msg,
						Color:       0x00cc11,
					},
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			err = fmt.Errorf("cannot respond to interaction: %v", err)
			log.Println(err)

			return err
		}
	}

	return nil
}
