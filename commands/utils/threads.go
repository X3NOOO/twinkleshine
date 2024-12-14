package utils

import "github.com/bwmarrin/discordgo"

func ParseReplies(s *discordgo.Session, m *discordgo.MessageCreate) (string, error) {
	fullMsg := m.Author.Username + ": \"" + m.Content + "\""
	currentM, err := s.ChannelMessage(m.ChannelID, m.ID)
	if err != nil {
		return "", err
	}

	for currentM.Type == discordgo.MessageTypeReply {
		fullMsg = currentM.Author.Username + ": \"" + currentM.ReferencedMessage.Content + "\"\n" + fullMsg

		currentM, err = s.ChannelMessage(m.ChannelID, currentM.ReferencedMessage.ID)
		if err != nil {
			break
		}
	}

	return fullMsg, nil
}
