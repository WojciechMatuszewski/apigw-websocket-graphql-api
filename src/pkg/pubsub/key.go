package pubsub

import (
	"fmt"
	"strings"
)

type Key struct {
	PK string `json:"pk"`
	SK string `json:"sk"`
}

func newKey(eventName string) Key {
	return Key{
		PK: fmt.Sprintf("PUBSUB#DATA"),
		SK: fmt.Sprintf("EVENT#%v", eventName),
	}
}

func (k Key) AsEvent() string {
	return strings.ReplaceAll(k.SK, "EVENT#", "")
}
