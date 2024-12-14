package commands

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/net/html"

	"github.com/X3NOOO/twinkleshine/ai"
	"github.com/X3NOOO/twinkleshine/commands/utils"
)

func rememberFile(s *discordgo.Session, i *discordgo.InteractionCreate, ai ai.TwinkleshineAI, log *log.Logger) error {
	options := i.ApplicationCommandData().Options

	fileOption := options[0].Options[0]
	attachment := i.ApplicationCommandData().Resolved.Attachments[fileOption.Value.(string)]
	log.Println("Processing file:", attachment.Filename)

	exists, err := ai.Exists("file.name", []any{attachment.Filename})
	if err != nil {
		msg := fmt.Sprintf("Failed to check if file exists: %v", err)
		log.Println(msg)
		return utils.SendErrorEmbed(msg, true, s, i)
	}
	if exists {
		msg := fmt.Sprintf("File %s already exists", attachment.Filename)
		log.Println(msg)
		return utils.SendErrorEmbed(msg, true, s, i)
	}

	resp, err := http.Get(attachment.URL)
	if err != nil {
		msg := fmt.Sprintf("Failed to download file: %v", err)
		log.Println(msg)
		return utils.SendErrorEmbed(msg, true, s, i)
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		msg := fmt.Sprintf("Failed to read file content: %v", err)
		log.Println(msg)
		return utils.SendErrorEmbed(msg, true, s, i)
	}

	hasher := sha256.New()
	_, err = hasher.Write(content)
	if err != nil {
		msg := fmt.Sprintf("Failed to hash file content: %v", err)
		log.Println(msg)
		return utils.SendErrorEmbed(msg, true, s, i)
	}
	hash := fmt.Sprintf("%x", hasher.Sum(nil))
	exists, err = ai.Exists("hash", []any{hash})
	if err != nil {
		msg := fmt.Sprintf("Failed to check if hash exists: %v", err)
		log.Println(msg)
		return utils.SendErrorEmbed(msg, true, s, i)
	}
	if exists {
		msg := fmt.Sprintf("File %s (with hash %s) already exists", attachment.Filename, hash)
		log.Println(msg)
		return utils.SendErrorEmbed(msg, true, s, i)
	}

	messageLink := fmt.Sprintf("https://discord.com/channels/%s/%s/%s", i.GuildID, i.ChannelID, i.ApplicationCommandData().TargetID)
	err = ai.RememberFile(content, map[string]any{
		"file": map[string]any{
			"name":    attachment.Filename,
			"url":     messageLink,
			"hash":    hash,
			"addedBy": i.Member.User.ID,
		},
	})
	if err != nil {
		msg := fmt.Sprintf("Failed to remember file: %v", err)
		log.Println(msg)
		return utils.SendErrorEmbed(msg, true, s, i)
	}

	log.Println("File uploaded successfully!")
	return utils.SendSuccessEmbed("File uploaded successfully!", true, s, i)
}

func rememberText(s *discordgo.Session, i *discordgo.InteractionCreate, ai ai.TwinkleshineAI, log *log.Logger) error {
	options := i.ApplicationCommandData().Options
	text := options[0].Options[0].StringValue()

	log.Println("Processing text:", text)

	hasher := sha256.New()
	_, err := hasher.Write([]byte(text))
	if err != nil {
		msg := fmt.Sprintf("Failed to hash text: %v", err)
		log.Println(msg)
		return utils.SendErrorEmbed(msg, true, s, i)
	}
	hash := fmt.Sprintf("%x", hasher.Sum(nil))

	exists, err := ai.Exists("hash", []any{hash})
	if err != nil {
		msg := fmt.Sprintf("Failed to check if hash exists: %v", err)
		log.Println(msg)
		return utils.SendErrorEmbed(msg, true, s, i)
	}
	if exists {
		msg := "Text already exists in the memory"
		log.Println(msg)
		return utils.SendErrorEmbed(msg, true, s, i)
	}

	err = ai.Remember(text, map[string]any{
		"file": map[string]any{
			"name":    i.Member.User.Username + "'s text",
			"hash":    hash,
			"addedBy": i.Member.User.ID,
		},
	})
	if err != nil {
		msg := fmt.Sprintf("Failed to remember text: %v", err)
		log.Println(msg)
		return utils.SendErrorEmbed(msg, true, s, i)
	}

	log.Println("Text remembered successfully!")
	return utils.SendSuccessEmbed("Text remembered successfully!", true, s, i)
}

func rememberUrls(s *discordgo.Session, i *discordgo.InteractionCreate, ai ai.TwinkleshineAI, log *log.Logger) error {
	options := i.ApplicationCommandData().Options
	urls := options[0].Options[0].StringValue()
	log.Println("Processing URLs:", urls)

	urlList := strings.Split(urls, " ")

	errs := make(chan error, len(urlList))

	for _, url := range urlList {
		go func(url string) {
			log.Println("Processing URL:", url)

			exists, err := ai.Exists("file.url", []any{url})
			if err != nil {
				err = fmt.Errorf("failed to check if URL exists: %v", err)
				errs <- err
			}
			if exists {
				err = fmt.Errorf("URL %s already exists", url)
				errs <- err
			}

			resp, err := http.Get(url)
			if err != nil {
				err = fmt.Errorf("failed to download URL: %v", err)
				errs <- err
			}

			content, err := io.ReadAll(resp.Body)
			if err != nil {
				err = fmt.Errorf("failed to read URL content: %v", err)
				errs <- err
			}

			hasher := sha256.New()
			_, err = hasher.Write(content)
			if err != nil {
				err = fmt.Errorf("failed to hash URL content: %v", err)
				errs <- err
			}
			hash := fmt.Sprintf("%x", hasher.Sum(nil))
			exists, err = ai.Exists("hash", []any{hash})
			if err != nil {
				err = fmt.Errorf("failed to check if hash exists: %v", err)
				errs <- err
			}
			if exists {
				err = fmt.Errorf("URL %s (with hash %s) already exists", url, hash)
				errs <- err
			}

			doc, err := html.Parse(strings.NewReader(string(content)))
			if err != nil {
				err = fmt.Errorf("failed to parse HTML: %v", err)
				errs <- err
			}

			var title string
			var maxDepth uint8 = 0xff
			var getTitle func(*html.Node, uint8)
			getTitle = func(n *html.Node, depth uint8) {
				if depth > maxDepth || n == nil {
					return
				}
				if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
					title = n.FirstChild.Data
					return
				}
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					getTitle(c, depth+1)
				}
			}
			getTitle(doc, 0)

			if title == "" {
				parts := strings.Split(url, ".")
				if len(parts) > 1 {
					title = "Website: " + parts[len(parts)-2]
				} else {
					title = url
				}
			}

			err = ai.RememberFile(content, map[string]any{
				"file": map[string]any{
					"name":    title,
					"url":     url,
					"hash":    hash,
					"addedBy": i.Member.User.ID,
				},
			})
			if err != nil {
				err = fmt.Errorf("failed to remember URL: %v", err)
				errs <- err
			}

			errs <- nil
		}(url)
	}

	var failedUploads []string
	for i := 0; i < len(urlList); i++ {
		if err := <-errs; err != nil {
			failedUploads = append(failedUploads, err.Error())
		}
	}

	if len(failedUploads) > 0 {
		msg := "Some URLs failed to upload:\n" + strings.Join(failedUploads, "\n") + "\n\nThe rest of the URLs were uploaded successfully."
		log.Println(msg)
		return utils.SendErrorEmbed(msg, true, s, i)
	}

	log.Println("All URLs remembered successfully!")
	return utils.SendSuccessEmbed("All URLs remembered successfully!", true, s, i)
}

func (c *CommandContext) RememberCLIHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to respond to interaction: %v", err)
	}

	options := i.ApplicationCommandData().Options
	subcommand := options[0].Name

	logger_prefix := fmt.Sprintf("[RememberCLIHandler %s] %s [%s] ", subcommand, i.Member.User.Username, i.Member.User.ID)
	logger := log.New(log.Writer(), logger_prefix, log.Flags())

	switch subcommand {
	case "file":
		err = rememberFile(s, i, c.AI, logger)
	case "text":
		err = rememberText(s, i, c.AI, logger)
	case "urls":
		err = rememberUrls(s, i, c.AI, logger)
	default:
		msg := fmt.Sprintf("Unknown subcommand: %v", subcommand)
		logger.Println(msg)
		return utils.SendErrorEmbed(msg, true, s, i)
	}
	if err != nil {
		log.Println(err)
	}

	return err
}

func (c *CommandContext) RememberGUIHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send initial response: %v", err)
	}

	logger_prefix := fmt.Sprintf("[RememberGUIHandler] %s [%s] ", i.Member.User.Username, i.Member.User.ID)
	log := log.New(log.Writer(), logger_prefix, log.Flags())
	log.Println("Processing remember command")

	attachments := i.ApplicationCommandData().Resolved.Messages[i.ApplicationCommandData().TargetID].Attachments
	if len(attachments) == 0 {
		msg := i.ApplicationCommandData().Resolved.Messages[i.ApplicationCommandData().TargetID]
		messageContent := msg.Content
		if messageContent == "" {
			msg := "No attachments or message content found."
			log.Println(msg)
			return utils.SendErrorEmbed(msg, true, s, i)
		}

		fullMsg, err := utils.ParseReplies(s, msg)
		if err != nil {
			msg := fmt.Sprintf("Failed to parse replies: %v", err)
			log.Println(msg)
			return utils.SendErrorEmbed(msg, true, s, i)
		}

		hasher := sha256.New()
		_, err = hasher.Write([]byte(fullMsg))
		if err != nil {
			msg := fmt.Sprintf("Failed to hash message content: %v", err)
			log.Println(msg)
			return utils.SendErrorEmbed(msg, true, s, i)
		}
		hash := fmt.Sprintf("%x", hasher.Sum(nil))

		exists, err := c.AI.Exists("file.hash", []any{hash})
		if err != nil {
			msg := fmt.Sprintf("Failed to check if hash exists: %v", err)
			log.Println(msg)
			return utils.SendErrorEmbed(msg, true, s, i)
		}
		if exists {
			msg := "Message content already exists in the memory"
			log.Println(msg)
			return utils.SendErrorEmbed(msg, true, s, i)
		}

		err = c.AI.Remember(fullMsg, map[string]any{
			"file": map[string]any{
				"name":    msg.Author.Username + "'s message",
				"url":     fmt.Sprintf("https://discord.com/channels/%s/%s/%s", i.GuildID, i.ChannelID, i.ApplicationCommandData().TargetID),
				"hash":    hash,
				"addedBy": i.Member.User.ID,
			},
		})
		if err != nil {
			msg := fmt.Sprintf("Failed to remember message content: %v", err)
			log.Println(msg)
			return utils.SendErrorEmbed(msg, true, s, i)
		}

		return utils.SendSuccessEmbed("Message content remembered successfully!", true, s, i)
	}

	errs := make(chan error, len(attachments))

	for _, att := range attachments {
		go func(att *discordgo.MessageAttachment) {
			log.Println("Processing attachment:", att.Filename)

			exists, err := c.AI.Exists("file.name", []any{att.Filename})
			if err != nil {
				err = fmt.Errorf("failed to check if %s exists: %v", att.Filename, err)
				errs <- err
				return
			}
			if exists {
				err = fmt.Errorf("file %s already exists", att.Filename)
				errs <- err
				return
			}

			resp, err := http.Get(att.URL)
			if err != nil {
				log.Printf("Failed to download %s: %v\n", att.Filename, err)
				errs <- fmt.Errorf("failed to download %s: %v", att.Filename, err)
				return
			}
			defer resp.Body.Close()

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Failed to read data %s: %v\n", att.Filename, err)
				errs <- fmt.Errorf("failed to read data %s: %v", att.Filename, err)
				return
			}

			hasher := sha256.New()
			hasher.Write(data)
			_, err = hasher.Write(data)
			if err != nil {
				log.Printf("Failed to hash data %s: %v\n", att.Filename, err)
				errs <- fmt.Errorf("failed to hash data %s: %v", att.Filename, err)
				return
			}
			hash := fmt.Sprintf("%x", hasher.Sum(nil))
			exists, err = c.AI.Exists("hash", []any{hash})
			if err != nil {
				log.Printf("Failed to check if hash %s exists: %v\n", hash, err)
				errs <- fmt.Errorf("failed to check if hash %s exists: %v", hash, err)
				return
			}
			if exists {
				err = fmt.Errorf("file %s (with hash %s) already exists", att.Filename, hash)
				errs <- err
				return
			}

			messageLink := fmt.Sprintf("https://discord.com/channels/%s/%s/%s", i.GuildID, i.ChannelID, i.ApplicationCommandData().TargetID)
			err = c.AI.RememberFile(data, map[string]any{
				"file": map[string]any{
					"name": att.Filename,
					"url":  messageLink,
					"hash": hash,
				},
			})
			if err != nil {
				log.Printf("Failed to remember file %s: %v\n", att.Filename, err)
				errs <- fmt.Errorf("failed to remember file %s: %v", att.Filename, err)
				return
			}

			errs <- nil
		}(att)
	}

	var failedUploads []string

	for i := 0; i < len(attachments); i++ {
		if err := <-errs; err != nil {
			failedUploads = append(failedUploads, err.Error())
		}
	}

	if len(failedUploads) > 0 {
		msg := "Some attachments failed to upload:\n" + strings.Join(failedUploads, "\n") + "\n\nThe rest of the attachments were uploaded successfully."
		log.Println(msg)
		return utils.SendErrorEmbed(msg, true, s, i)
	}

	log.Println("All attachments uploaded successfully!")

	return utils.SendSuccessEmbed("All attachments uploaded successfully!", true, s, i)
}
