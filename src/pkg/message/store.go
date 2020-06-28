package message

import (
	"context"
	"time"

	"apigw-graphqlv2/src/pkg/gql/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
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

func (s Store) SaveMessage(ctx context.Context, input model.MessageInput) (Message, error) {
	msg := Message{
		ID:        uuid.Must(uuid.NewRandom()).String(),
		Message:   input.Message,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	msgItem := msg.AsItem()

	dbItem, err := dynamodbattribute.MarshalMap(&msgItem)
	if err != nil {
		return Message{}, err
	}

	req := s.db.PutItemRequest(&dynamodb.PutItemInput{
		Item:      dbItem,
		TableName: aws.String(s.table),
	})

	_, err = req.Send(ctx)
	if err != nil {
		return Message{}, err
	}

	return msg, nil
}
