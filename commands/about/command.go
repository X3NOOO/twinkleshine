package about

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func AboutHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
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

	return err
}
