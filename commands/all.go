package commands

import (
	"github.com/X3NOOO/twinkleshine/commands/about"
	"github.com/X3NOOO/twinkleshine/commands/ask"
	"github.com/X3NOOO/twinkleshine/commands/remember"
	"github.com/bwmarrin/discordgo"
)

var ALL_COMMANDS = []Command{
	{
		Name:        "about",
		Description: "About the bot",
		Handler:     about.AboutHandler,
	},
	{
		Name:        "ask",
		Description: "Ask a question",
		Handler:     ask.AskHandler,
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
		Name:    "Add to the persistent knowledge",
		Type:    discordgo.MessageApplicationCommand,
		Handler: remember.RememberGUIHandler,
	},
	{
		Name:        "remember",
		Description: "Add to the persistent knowledge",
		Handler:     remember.RememberCLIHandler,
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
