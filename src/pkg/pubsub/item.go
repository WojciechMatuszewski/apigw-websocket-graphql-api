package pubsub

type Item struct {
	EventName string `json:"eventName"`
	Data      string `json:"data"`
}

func (m Item) asItem() messageItem {
	return messageItem{
		Key:  newKey(m.EventName),
		Data: m.Data,
	}
}

type messageItem struct {
	Key
	Data string `json:"data"`
}
