package utils

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type cooldownLookup struct {
	UserID string
	Time   int64
}

type Security struct {
	StaffRoleID     string
	SlowmodeSeconds int64

	cooldowns []cooldownLookup
	mu        sync.Mutex
}

func (sc *Security) GetCooldown(useID string) *cooldownLookup {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	for _, lookup := range sc.cooldowns {
		if lookup.UserID == useID {
			return &lookup
		}
	}

	return nil
}

func (sc *Security) SetCooldown(lookup cooldownLookup) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	for i, l := range sc.cooldowns {
		if l.UserID == lookup.UserID {
			sc.cooldowns[i] = lookup
			return
		}
	}

	sc.cooldowns = append(sc.cooldowns, lookup)
}

func (sc *Security) Guard(handler func(s *discordgo.Session, i *discordgo.InteractionCreate) error) func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
		member, err := s.GuildMember(i.GuildID, i.Member.User.ID)
		if err != nil {
			return err
		}

		for _, role := range member.Roles {
			if role == sc.StaffRoleID {
				return handler(s, i)
			}
		}

		SendErrorEmbed("You are not authorized to use this command", false, s, i)

		return fmt.Errorf("user %s [%s] tried to use a privileged command", i.Member.User.Username, i.Member.User.ID)
	}
}

func (sc *Security) Timeout(handler func(s *discordgo.Session, i *discordgo.InteractionCreate) error) func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
		cooldown := sc.GetCooldown(i.Member.User.ID)
		if cooldown != nil {
			if cooldown.Time+sc.SlowmodeSeconds > time.Now().Unix() {
				msg := fmt.Sprintf("Please wait %d seconds before using this command again", cooldown.Time+sc.SlowmodeSeconds-time.Now().Unix())
				SendErrorEmbed(msg, false, s, i)

				return fmt.Errorf("%s's [%s] cooldown has not expired", i.Member.User.Username, i.Member.User.ID)
			}
		}

		sc.SetCooldown(cooldownLookup{
			UserID: i.Member.User.ID,
			Time:   time.Now().Unix(),
		})

		return handler(s, i)
	}
}
