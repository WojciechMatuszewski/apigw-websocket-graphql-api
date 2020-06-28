package gql

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/executor"
)

func Execute(ctx context.Context, params graphql.RawParams,
	schema graphql.ExecutableSchema) (*graphql.Response, error) {

	ex := executor.New(schema)
	ctx = graphql.StartOperationTrace(ctx)

	rc, err := ex.CreateOperationContext(ctx, &params)
	if err != nil {
		return nil, err
	}

	opRes, ctx2 := ex.DispatchOperation(ctx, rc)
	out := opRes(ctx2)
	return out, nil
}
