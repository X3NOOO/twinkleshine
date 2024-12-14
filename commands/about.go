package commands

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func (c *CommandContext) AboutHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	logger_prefix := fmt.Sprintf("[ABOUT] %s [%s] ", i.Member.User.Username, i.Member.User.ID)
	log := log.New(log.Writer(), logger_prefix, log.Flags())
	log.Println("About command is called")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Author: &discordgo.MessageEmbedAuthor{
						Name:    "BudCare Galaxy",
						IconURL: "https://raw.githubusercontent.com/X3NOOO/twinkleshine/refs/heads/master/assets/emblem.png",
					},
					Description: "TODO",
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Source",
							Value:  "[GitHub](https://github.com/X3NOOO/twinkleshine)",
							Inline: true,
						},
					},
					Color: 0x5c29e4,
					Type:  discordgo.EmbedTypeRich,
				},
			},
		},
	})
	if err != nil {
		log.Printf("Cannot respond to interaction: %v\n", err)
	}

	log.Println("About command is done")

	return err
}
