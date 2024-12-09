package utils

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func SendErrorEmbed(msg string, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	log.Println(msg)

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
	}

	return err
}

func SendSuccessEmbed(msg string, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	log.Println(msg)

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
	}

	return err
}
