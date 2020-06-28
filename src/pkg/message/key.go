package message

import (
	"fmt"
)

type Key struct {
	PK string `json:"pk"`
	SK string `json:"sk"`
}

func newKey(createdAt string) Key {
	return Key{
		PK: fmt.Sprintf("CHAT#MESSAGE"),
		SK: createdAt,
	}
}

func (k Key) AsCreatedAt() string {
	return k.SK
}
