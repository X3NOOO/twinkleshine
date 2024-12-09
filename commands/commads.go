package commands

import (
	"github.com/X3NOOO/twinkleshine/commands/utils"

	"github.com/X3NOOO/twinkleshine/ai"
	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Name        string
	Description string
	Handler     func(s *discordgo.Session, i *discordgo.InteractionCreate) error
	Options     []*discordgo.ApplicationCommandOption
	Type        discordgo.ApplicationCommandType
}

type CommandContext struct {
	AI ai.TwinkleshineAI
}

func (c *CommandContext) GetCommands(staffRole string, slowmodeSeconds int64) []Command {
	security := utils.Security{
		StaffRoleID:     staffRole,
		SlowmodeSeconds: slowmodeSeconds,
	}

	return []Command{
		{
			Name:        "about",
			Description: "About the bot",
			Handler:     c.AboutHandler,
		},
		{
			Name:        "ask",
			Description: "Ask a question",
			Handler:     security.Timeout(c.AskCLIHandler),
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
			Handler: security.Timeout(c.AskGUIHandler),
		},
		{
			Name:    "Add to the persistent knowledge",
			Type:    discordgo.MessageApplicationCommand,
			Handler: security.Guard(c.RememberGUIHandler),
		},
		{
			Name:        "remember",
			Description: "Add to the persistent knowledge",
			Handler:     security.Guard(c.RememberCLIHandler),
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
				},
				{
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
				{
					Name:        "urls",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Description: "Upload websites to remember",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "urls",
							Description: "Space separated URLs to remember",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
			},
		},
	}

}
