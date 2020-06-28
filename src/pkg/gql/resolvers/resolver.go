package resolvers

//go:generate go run github.com/99designs/gqlgen

import (
	"apigw-graphqlv2/src/pkg/message"
	"apigw-graphqlv2/src/pkg/pubsub"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	MessageStore  message.Store
	PubSubService pubsub.PubSubService
}
