package message

type Message struct {
	ID        string `json:"id"`
	Message   string `json:"message"`
	CreatedAt string `json:"createdAt"`
}

func (m Message) AsItem() messageItem {
	return messageItem{
		Key:     newKey(m.CreatedAt),
		ID:      m.ID,
		Message: m.Message,
	}
}

type messageItem struct {
	Key
	ID      string `json:"id"`
	Message string `json:"message"`
}

func (mi messageItem) AsMessage() Message {
	return Message{
		ID:        mi.ID,
		Message:   mi.Message,
		CreatedAt: mi.AsCreatedAt(),
	}
}
