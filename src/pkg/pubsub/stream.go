package pubsub

import (
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
)

func EventFromStreamRecord(record events.DynamoDBEventRecord) (string, error) {
	var key Key
	err := unmarshalStreamImage(record.Change.Keys, &key)
	if err != nil {
		return "", err
	}

	if !strings.Contains(key.PK, "PUBSUB#DATA") {
		return "", nil
	}

	return key.AsEvent(), nil
}

func unmarshalStreamImage(attribute map[string]events.DynamoDBAttributeValue, out interface{}) error {
	attrsMap := make(map[string]dynamodb.AttributeValue, len(attribute))
	for k, v := range attribute {
		bytes, err := v.MarshalJSON()
		if err != nil {
			return err
		}
		var dbAttr dynamodb.AttributeValue
		err = json.Unmarshal(bytes, &dbAttr)
		if err != nil {
			return err
		}
		attrsMap[k] = dbAttr
	}
	return dynamodbattribute.UnmarshalMap(attrsMap, out)
}
