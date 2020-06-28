package pubsub

import (
	"context"

	"apigw-graphqlv2/src/pkg/connection"
)

type PubSubService struct {
	connectionStore connection.Store
	store           Store
}

func NewService(connectionStore connection.Store, store Store) PubSubService {
	return PubSubService{
		connectionStore: connectionStore,
		store:           store,
	}
}

func (pb PubSubService) Subscribe(ctx context.Context, eventName string) ([]byte, error) {
	metadata := metadataFromContext(ctx)
	if metadata == nil {
		return pb.store.Get(ctx, eventName)
	}

	conn := connection.Connection{
		GqlID:  metadata.GqlID,
		ID:     metadata.ID,
		Params: metadata.Params,
		Event:  eventName,
	}
	err := pb.connectionStore.SaveConnection(ctx, conn)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (pb PubSubService) Publish(ctx context.Context, eventName string, data []byte) error {
	return pb.store.Save(ctx, eventName, data)
}
