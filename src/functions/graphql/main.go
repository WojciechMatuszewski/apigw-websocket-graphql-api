package main

import (
	"log"
	"net/http"
	"os"

	"apigw-graphqlv2/src/pkg/connection"
	"apigw-graphqlv2/src/pkg/gql/generated"
	"apigw-graphqlv2/src/pkg/gql/resolvers"
	"apigw-graphqlv2/src/pkg/message"
	"apigw-graphqlv2/src/pkg/pubsub"

	"github.com/rs/cors"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/apex/gateway"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic(err.Error())
	}

	db := dynamodb.New(cfg)
	messageStore := message.NewStore(db, os.Getenv("DATA_TABLE_NAME"))
	connectionStore := connection.NewStore(db, os.Getenv("DATA_TABLE_NAME"))
	pubSubStore := pubsub.NewStore(db, os.Getenv("DATA_TABLE_NAME"))
	pubSubService := pubsub.NewService(connectionStore, pubSubStore)

	schema := generated.NewExecutableSchema(generated.Config{
		Resolvers: &resolvers.Resolver{MessageStore: messageStore, PubSubService: pubSubService},
	})

	svr := handler.NewDefaultServer(schema)
	svr.AddTransport(transport.POST{})

	mux := http.NewServeMux()
	mux.Handle("/graphql", svr)

	h := cors.Default().Handler(mux)

	log.Fatal(gateway.ListenAndServe(":3000", h))
}
