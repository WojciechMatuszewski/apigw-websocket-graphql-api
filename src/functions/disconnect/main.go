package main

import (
	"context"
	"fmt"
	"os"

	"apigw-graphqlv2/src/pkg/connection"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func handler(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) error {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return err
	}
	db := dynamodb.New(cfg)
	connectionStore := connection.NewStore(db, os.Getenv("DATA_TABLE_NAME"))

	err = connectionStore.DeleteConnectionsByID(ctx, event.RequestContext.ConnectionID)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
