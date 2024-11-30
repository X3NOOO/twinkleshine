package main

import (
	"log"

	"github.com/X3NOOO/twinkleshine/commands"
	"github.com/bwmarrin/discordgo"
)

type bot struct {
	s        *discordgo.Session
	commands []commands.Command
}

func (b *bot) Run() error {
	err := b.s.Open()
	if err != nil {
		log.Printf("Cannot open the session: %v\n", err)
		return err
	}

	log.Println("Bot is running")
	return nil
}

func (b *bot) Stop() error {
	return b.s.Close()
}

func (b *bot) onReady(s *discordgo.Session, event *discordgo.Ready) {
	log.Println("Bot is ready")

	existingCmds, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		log.Printf("Cannot fetch existing commands: %v\n", err)
		return
	}

	// Delete commands that don't have handlers
	for _, existing := range existingCmds {
		hasHandler := false
		for _, cmd := range b.commands {
			if existing.Name == cmd.Name {
				hasHandler = true
				break
			}
		}
		if !hasHandler {
			err := s.ApplicationCommandDelete(s.State.User.ID, "", existing.ID)
			if err != nil {
				log.Printf("Cannot delete command %s: %v\n", existing.Name, err)
			} else {
				log.Printf("Deleted unused command: %s\n", existing.Name)
			}
		}
	}

	for _, cmd := range b.commands {
		exists := false
		for _, existing := range existingCmds {
			if existing.Name == cmd.Name &&
				existing.Description == cmd.Description &&
				compareOptions(existing.Options, cmd.Options) {
				exists = true
				break
			}
		}

		if !exists {
			log.Printf("Registering command: %+v\n", cmd)
			_, err := s.ApplicationCommandCreate(s.State.User.ID, "", &discordgo.ApplicationCommand{
				Name:        cmd.Name,
				Description: cmd.Description,
				Options:     cmd.Options,
				Type:        cmd.Type,
			})

			if err != nil {
				log.Printf("Cannot create command %s: %v\n", cmd.Name, err)
			}
		}
	}
}

func compareOptions(a, b []*discordgo.ApplicationCommandOption) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Name != b[i].Name ||
			a[i].Description != b[i].Description ||
			a[i].Type != b[i].Type ||
			a[i].Required != b[i].Required {
			return false
		}
	}
	return true
}

func (b *bot) onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		for _, cmd := range b.commands {
			if i.ApplicationCommandData().Name == cmd.Name {
				err := cmd.Handler(s, i)
				if err != nil {
					log.Printf("Cannot handle command: %v\n", err)
				}
				return
			}
		}
	}
}

func NewBot(token string, commands []commands.Command) (*bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Printf("Invalid bot parameters: %v\n", err)
		return nil, err
	}

	bot := &bot{
		s:        session,
		commands: commands,
	}

	bot.s.AddHandler(bot.onReady)
	bot.s.AddHandler(bot.onInteractionCreate)

	return bot, nil
}
