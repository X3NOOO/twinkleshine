package commands

import "github.com/bwmarrin/discordgo"

type Command struct {
	Name        string
	Description string
	Handler     func(s *discordgo.Session, i *discordgo.InteractionCreate) error
	Options     []*discordgo.ApplicationCommandOption
	Type        discordgo.ApplicationCommandType
}

type CommandContext struct {
	AI interface {
		Query(text string) (string, error)
	}
}

func (c *CommandContext) GetCommands() []Command {
	return []Command{
		{
			Name:        "about",
			Description: "About the bot",
			Handler:     c.AboutHandler,
		},
		{
			Name:        "ask",
			Description: "Ask a question",
			Handler:     c.AskCLIHandler,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "question",
					Description: "The question you want to ask",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
		{
			Name:    "Respond to the question",
			Type:    discordgo.MessageApplicationCommand,
			Handler: c.AskGUIHandler,
		},
		{
			Name:    "Add to the persistent knowledge",
			Type:    discordgo.MessageApplicationCommand,
			Handler: c.RememberGUIHandler,
		},
		{
			Name:        "remember",
			Description: "Add to the persistent knowledge",
			Handler:     c.RememberCLIHandler,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "file",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Description: "Upload a file to remember",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "file",
							Description: "File to remember",
							Type:        discordgo.ApplicationCommandOptionAttachment,
							Required:    true,
						},
					},
				}, {
					Name:        "text",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Description: "Upload text to remember",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "text",
							Description: "Text to remember",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
			},
		},
	}

}
