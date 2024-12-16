package discord

import (
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/X3NOOO/twinkleshine/ai"
	"github.com/X3NOOO/twinkleshine/commands"
	"github.com/X3NOOO/twinkleshine/commands/utils"
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

func (b *bot) getOnMessageCreateHandler(ctx *commands.CommandContext) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	logger_prefix := "[KNOWLEDGE WATCHDOG] "
	log := log.New(log.Writer(), logger_prefix, log.Flags())
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.Bot {
			return
		}
		author, err := s.GuildMember(m.GuildID, m.Author.ID)
		if err != nil {
			log.Printf("Cannot get member %s [%s]: %v\n", m.Author.Username, m.Author.ID, err)
			return
		}

		hasRole := false
		for _, roleID := range author.Roles {
			if roleID == ctx.AI.Options.Config.Discord.LearnMessagesRoleID {
				hasRole = true
				break
			}
		}
		if !hasRole {
			return
		}

		msg, err := s.ChannelMessage(m.ChannelID, m.ID)
		if err != nil {
			log.Printf("Cannot get message: %v\n", err)
			return
		}

		fullMsg, err := utils.ParseReplies(s, msg)
		if err != nil {
			log.Printf("Cannot parse replies: %v\n", err)
			return
		}

		hasher := sha256.New()
		_, err = hasher.Write([]byte(fullMsg))
		if err != nil {
			log.Printf("Cannot hash message: %v\n", err)
			return
		}
		hash := fmt.Sprintf("%x", hasher.Sum(nil))

		exists, err := ctx.AI.Exists("file.hash", []any{hash})
		if err != nil {
			log.Printf("Cannot check if message exists: %v\n", err)
			return
		}
		if exists {
			log.Println("Message already exists")
			return
		}

		err = ctx.AI.Remember(fullMsg, map[string]interface{}{
			"file": map[string]interface{}{
				"name":    m.Author.Username + "'s message",
				"url":     fmt.Sprintf("https://discord.com/channels/%s/%s/%s", m.GuildID, m.ChannelID, m.ID),
				"hash":    hash,
				"addedBy": "knowledge watchdog",
			},
		})
		if err != nil {
			log.Printf("Cannot remember message: %v\n", err)
			return
		}

		log.Printf("Remembered message: %s\n", fullMsg)
	}
}

func NewBot(token string, configPath string) (*bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Printf("Invalid bot parameters: %v\n", err)
		return nil, err
	}

	ai, err := ai.NewAI(configPath)
	if err != nil {
		log.Printf("Cannot create AI: %v\n", err)
		return nil, err
	}

	ctx := &commands.CommandContext{
		AI: *ai,
	}

	commands := ctx.GetCommands(
		ai.Options.Config.Discord.Security.StaffRoleID,
		ai.Options.Config.Discord.Security.CooldownSeconds,
	)

	bot := &bot{
		s:        session,
		commands: commands,
	}

	onMessageCreate := bot.getOnMessageCreateHandler(ctx)

	bot.s.AddHandler(bot.onReady)
	bot.s.AddHandler(bot.onInteractionCreate)
	bot.s.AddHandler(onMessageCreate)

	return bot, nil
}
