package connection

import (
	"context"

	"apigw-graphqlv2/src/pkg/gql"

	"github.com/99designs/gqlgen/graphql"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi/apigatewaymanagementapiiface"
)

type Sender struct {
	api apigatewaymanagementapiiface.ClientAPI
}

func NewSender(api apigatewaymanagementapiiface.ClientAPI) Sender {
	return Sender{api: api}
}

func (s Sender) SendPayload(ctx context.Context,
	connectionID string, ID string, response *graphql.Response) error {
	buf, err := gql.NewDataPayload(ID, response)
	if err != nil {
		return nil
	}

	return s.SendData(ctx, connectionID, buf)

}

func (s Sender) SendACK(ctx context.Context, connectionID string, gqlID string) error {
	buf, err := gql.NewACKPayload(gqlID)
	if err != nil {
		return err
	}

	return s.SendData(ctx, connectionID, buf)
}

func (s Sender) SendConnectionError(ctx context.Context, connectionID string, gqlID string) error {
	buf, err := gql.NewConnectionErrorPayload(gqlID)
	if err != nil {
		return err
	}

	return s.SendData(ctx, connectionID, buf)
}

func (s Sender) SendData(ctx context.Context, connectionID string, data []byte) error {
	req := s.api.PostToConnectionRequest(&apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(connectionID),
		Data:         data,
	})

	_, err := req.Send(ctx)
	return err
}
