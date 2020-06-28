package pubsub

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/expression"
)

type Store struct {
	db    dynamodbiface.ClientAPI
	table string
}

func NewStore(db dynamodbiface.ClientAPI, table string) Store {
	return Store{
		db:    db,
		table: table,
	}
}

func (s Store) Get(ctx context.Context, eventName string) ([]byte, error) {
	keys := newKey(eventName)

	keyCond := expression.KeyAnd(expression.Key("pk").Equal(
		expression.Value(keys.PK)), expression.Key("sk").Equal(expression.Value(keys.SK)))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, err
	}

	req := s.db.QueryRequest(&dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		TableName:                 aws.String(s.table),
	})

	out, err := req.Send(ctx)
	if err != nil {
		return nil, err
	}

	if len(out.Items) != 1 {
		panic(fmt.Errorf("items length is different than 1"))
	}

	var item Item
	outItem := out.Items[0]
	err = dynamodbattribute.UnmarshalMap(outItem, &item)
	if err != nil {
		return nil, err
	}

	return []byte(item.Data), nil
}

func (s Store) Save(ctx context.Context, eventName string, data []byte) error {
	publishItem := Item{
		EventName: eventName,
		Data:      string(data),
	}.asItem()

	dbItem, err := dynamodbattribute.MarshalMap(&publishItem)
	if err != nil {
		return err
	}

	req := s.db.PutItemRequest(&dynamodb.PutItemInput{
		Item:      dbItem,
		TableName: aws.String(s.table),
	})

	_, err = req.Send(ctx)
	return err
}
