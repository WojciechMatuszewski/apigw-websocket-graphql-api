package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"apigw-graphqlv2/src/pkg/connection"
	"apigw-graphqlv2/src/pkg/gql"
	"apigw-graphqlv2/src/pkg/gql/generated"
	"apigw-graphqlv2/src/pkg/gql/resolvers"
	"apigw-graphqlv2/src/pkg/message"
	"apigw-graphqlv2/src/pkg/pubsub"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context,
	event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connID := event.RequestContext.ConnectionID

	var data gql.SubscriptionData
	err := json.Unmarshal([]byte(event.Body), &data)
	if err != nil {
		return proxyResponse(http.StatusInternalServerError), err
	}

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return proxyResponse(http.StatusInternalServerError), err
	}

	apigwCfg := cfg
	apigwCfg.EndpointResolver = aws.ResolveWithEndpointURL(os.Getenv("APIGW_ADDR"))
	apigw := apigatewaymanagementapi.New(apigwCfg)
	connectionSender := connection.NewSender(apigw)

	db := dynamodb.New(cfg)
	messageStore := message.NewStore(db, os.Getenv("DATA_TABLE_NAME"))
	connectionStore := connection.NewStore(db, os.Getenv("DATA_TABLE_NAME"))
	pubSubStore := pubsub.NewStore(db, os.Getenv("DATA_TABLE_NAME"))
	pubSubService := pubsub.NewService(connectionStore, pubSubStore)

	schema := generated.NewExecutableSchema(generated.Config{
		Resolvers: &resolvers.Resolver{
			MessageStore:  messageStore,
			PubSubService: pubSubService,
		},
	})

	if data.Type == gql.GQL_CONNECTION_INIT {
		err = connectionSender.SendACK(ctx, connID, data.ID)
		if err != nil {
			return proxyResponse(http.StatusInternalServerError), err
		}
		return proxyResponse(http.StatusOK), nil
	}

	if data.Type != gql.GQL_START {
		return proxyResponse(http.StatusOK), nil
	}

	fmt.Println(data.Type, connID)

	paramsBuf, err := json.Marshal(data.Payload)
	if err != nil {
		fmt.Println(err.Error(), "paramsBuf")
		err = connectionSender.SendConnectionError(ctx, connID, data.ID)
		if err != nil {
			fmt.Println(err.Error(), "after paramsBuf")
			return proxyResponse(http.StatusInternalServerError), err
		}
		return proxyResponse(http.StatusOK), nil
	}

	metadata := pubsub.CtxMetada{
		GqlID:  data.ID,
		ID:     connID,
		Params: string(paramsBuf),
	}
	ctx = pubsub.WithContextMetadata(ctx, metadata)
	_, err = gql.Execute(ctx, data.Payload, schema)
	if err != nil {
		fmt.Println(err.Error(), "execute")
		err = connectionSender.SendConnectionError(ctx, connID, data.ID)
		if err != nil {
			fmt.Println(err.Error(), "after execute")
			return proxyResponse(http.StatusInternalServerError), err
		}
		return proxyResponse(http.StatusOK), nil
	}

	return proxyResponse(http.StatusOK), err
}

func main() {
	lambda.Start(handler)
}

func proxyResponse(statusCode int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       http.StatusText(statusCode),
		Headers: map[string]string{
			"Sec-WebSocket-Protocol": "graphql-ws",
		},
	}
}
