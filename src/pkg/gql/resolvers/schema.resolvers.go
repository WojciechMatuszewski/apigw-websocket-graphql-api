package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"fmt"

	"apigw-graphqlv2/src/pkg/gql/generated"
	"apigw-graphqlv2/src/pkg/gql/model"
	"apigw-graphqlv2/src/pkg/message"
)

func (r *mutationResolver) Message(ctx context.Context, input model.MessageInput) (*message.Message, error) {
	msg, err := r.MessageStore.SaveMessage(ctx, input)
	if err != nil {
		return nil, err
	}

	buf, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	err = r.PubSubService.Publish(ctx, "MESSAGE_ADDED", buf)
	if err != nil {
		return nil, err
	}

	err = r.PubSubService.Publish(ctx, "OTHER_MESSAGE_ADDED", buf)
	if err != nil {
		return nil, err
	}

	return &msg, err
}

func (r *queryResolver) Foo(ctx context.Context) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *subscriptionResolver) Message(ctx context.Context) (<-chan *message.Message, error) {
	data, err := r.PubSubService.Subscribe(ctx, "MESSAGE_ADDED")
	if err != nil {
		return nil, err
	}

	messageChan := make(chan *message.Message, 1)
	if data == nil {
		messageChan <- &message.Message{}
		return messageChan, nil
	}

	var msg message.Message
	err = json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	messageChan <- &msg
	return messageChan, nil
}

func (r *subscriptionResolver) OtherMessage(ctx context.Context) (<-chan *message.Message, error) {
	data, err := r.PubSubService.Subscribe(ctx, "OTHER_MESSAGE_ADDED")
	if err != nil {
		return nil, err
	}

	messageChan := make(chan *message.Message, 1)
	if data == nil {
		messageChan <- &message.Message{}
		return messageChan, nil
	}

	var msg message.Message
	err = json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	messageChan <- &msg
	return messageChan, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
