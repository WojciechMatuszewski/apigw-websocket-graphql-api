package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context,
	event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       http.StatusText(http.StatusOK),
		Headers: map[string]string{
			"Sec-WebSocket-Protocol": "graphql-ws",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
