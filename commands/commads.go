package commands

import "github.com/bwmarrin/discordgo"

type Command struct {
	Name        string
	Description string
	Handler     func(s *discordgo.Session, i *discordgo.InteractionCreate) error
	Options     []*discordgo.ApplicationCommandOption
	Type        discordgo.ApplicationCommandType
}
