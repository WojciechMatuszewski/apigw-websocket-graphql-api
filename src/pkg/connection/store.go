package connection

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

func (s Store) SaveConnection(ctx context.Context, conn Connection) error {
	connItem := conn.AsItem()
	dbItem, err := dynamodbattribute.MarshalMap(&connItem)
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

func (s Store) GetConnectionsByEventName(ctx context.Context, eventName string) ([]Connection, error) {
	pk := newPK(eventName)

	keyCond := expression.Key("pk").Equal(expression.Value(pk))
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

	connsItems := make([]connectionItem, len(out.Items))
	err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &connsItems)
	if err != nil {
		return nil, err
	}

	conns := make([]Connection, len(connsItems))
	for i, connItem := range connsItems {
		conns[i] = connItem.AsConnection()
	}

	return conns, nil
}

func (s Store) GetConnectionsByID(ctx context.Context, connectionID string) ([]connectionItem, error) {
	keyCond := expression.Key("connectionID").Equal(expression.Value(connectionID))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, err
	}

	req := s.db.QueryRequest(&dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		IndexName:                 aws.String("ByConnectionID"),
		KeyConditionExpression:    expr.KeyCondition(),
		TableName:                 aws.String(s.table),
	})

	out, err := req.Send(ctx)
	if err != nil {
		return nil, err
	}

	connsItems := make([]connectionItem, len(out.Items))
	err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &connsItems)
	if err != nil {
		return nil, err
	}

	return connsItems, nil
}

func (s Store) DeleteConnectionsByID(ctx context.Context, connectionID string) error {
	connItems, err := s.GetConnectionsByID(ctx, connectionID)
	if err != nil {
		return err
	}

	wReqs := make([]dynamodb.WriteRequest, len(connItems))
	for i, connItem := range connItems {
		wReqs[i] = dynamodb.WriteRequest{
			DeleteRequest: &dynamodb.DeleteRequest{Key: connItem.Key.AsAttrs()},
		}
	}

	req := s.db.BatchWriteItemRequest(&dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]dynamodb.WriteRequest{
			s.table: wReqs,
		},
	})

	out, err := req.Send(ctx)
	if err != nil {
		return err
	}

	fmt.Println("unprocessed", out.UnprocessedItems)
	return nil
}
