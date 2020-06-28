package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"apigw-graphqlv2/src/pkg/connection"
	"apigw-graphqlv2/src/pkg/gql"
	"apigw-graphqlv2/src/pkg/gql/generated"
	"apigw-graphqlv2/src/pkg/gql/resolvers"
	"apigw-graphqlv2/src/pkg/message"
	"apigw-graphqlv2/src/pkg/pubsub"

	"github.com/99designs/gqlgen/graphql"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/aws/aws-sdk-go-v2/aws/external"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.DynamoDBEvent) error {
	record := event.Records[0]
	if events.DynamoDBOperationType(record.EventName) == events.DynamoDBOperationTypeRemove {
		return nil
	}

	eventName, err := pubsub.EventFromStreamRecord(record)
	fmt.Println(eventName)
	if err != nil {
		return err
	}

	if eventName == "" {
		return nil
	}

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return err
	}
	apigwCfg := cfg
	apigwCfg.EndpointResolver = aws.ResolveWithEndpointURL(os.Getenv("APIGW_ADDR"))
	db := dynamodb.New(cfg)
	apigw := apigatewaymanagementapi.New(apigwCfg)

	connectionStore := connection.NewStore(db, os.Getenv("DATA_TABLE_NAME"))
	connectionSender := connection.NewSender(apigw)

	messageStore := message.NewStore(db, os.Getenv("DATA_TABLE_NAME"))
	pubSubStore := pubsub.NewStore(db, os.Getenv("DATA_TABLE_NAME"))
	pubSubService := pubsub.NewService(connectionStore, pubSubStore)

	schema := generated.NewExecutableSchema(generated.Config{
		Resolvers: &resolvers.Resolver{
			MessageStore:  messageStore,
			PubSubService: pubSubService,
		},
	})

	connections, err := connectionStore.GetConnectionsByEventName(ctx, eventName)
	if err != nil {
		return err
	}

	for _, conn := range connections {
		var params graphql.RawParams
		err = json.Unmarshal([]byte(conn.Params), &params)
		if err != nil {
			return err
		}

		response, err := gql.Execute(ctx, params, schema)
		if err != nil {
			return err
		}

		err = connectionSender.SendPayload(ctx, conn.ID, conn.GqlID, response)
		if err != nil {
			fmt.Println("sending the payload failed", err.Error())
		}
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
