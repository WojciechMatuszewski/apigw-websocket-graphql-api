package gql

import (
	"encoding/json"

	"github.com/99designs/gqlgen/graphql"
)

const (
	GQL_CONNECTION_INIT  = "connection_init"
	GQL_START            = "start"
	GQL_CONNECTION_ACK   = "connection_ack"
	GQL_CONNECTION_ERROR = "connection_error"
	GQL_DATA             = "data"
	GQL_ERROR            = "error"
)

type SubscriptionData struct {
	Type    string            `json:"type,omitempty"`
	ID      string            `json:"id,omitempty"`
	Payload graphql.RawParams `json:"payload,omitempty"`
}

type SubscriptionResponse struct {
	Type    string            `json:"type,omitempty"`
	ID      string            `json:"id,omitempty"`
	Payload *graphql.Response `json:"payload,omitempty"`
}

func NewACKPayload(ID string) ([]byte, error) {
	return json.Marshal(SubscriptionResponse{Type: GQL_CONNECTION_ACK, ID: ID})
}

func NewConnectionErrorPayload(ID string) ([]byte, error) {
	return json.Marshal(SubscriptionResponse{Type: GQL_CONNECTION_ERROR, ID: ID})
}

func NewDataPayload(ID string, payload *graphql.Response) ([]byte, error) {
	return json.Marshal(SubscriptionResponse{Type: GQL_DATA, ID: ID, Payload: payload})
}
