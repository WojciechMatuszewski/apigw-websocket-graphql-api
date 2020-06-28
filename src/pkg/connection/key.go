package connection

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Key struct {
	PK string `json:"pk"`
	SK string `json:"sk"`
}

func NewKey(eventName string, connectionID string) Key {
	return Key{
		PK: newPK(eventName),
		SK: newSK(connectionID),
	}
}

func (k Key) AsEventName() string {
	return strings.ReplaceAll(k.PK, "EVENT#", "")
}

func (k Key) AsConnectionID() string {
	return strings.ReplaceAll(k.SK, "CONNECTION#", "")
}

func (k Key) AsAttrs() map[string]dynamodb.AttributeValue {
	return map[string]dynamodb.AttributeValue{
		"pk": {
			S: aws.String(k.PK),
		},
		"sk": {
			S: aws.String(k.SK),
		},
	}
}

func newPK(eventName string) string {
	return fmt.Sprintf("EVENT#%v", eventName)
}

func newSK(connectionID string) string {
	return fmt.Sprintf("CONNECTION#%v", connectionID)
}
