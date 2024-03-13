package chat

func MessageHandle(m *Message) *Message {
	switch m.MessageType {
	case string(Text):
		return m
	case string(Emoji):
		return m
	case string(Image):
		// Handle image messages
		return m
	case string(Video):
		return m
	default:
		return m
	}
}
