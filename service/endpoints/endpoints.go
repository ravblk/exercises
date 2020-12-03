package endpoints

import (
	"context"

	service "ravblk/exercises/service/services"
	"ravblk/exercises/service/services/brackets"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	Fix      endpoint.Endpoint
	Validate endpoint.Endpoint
}

func MakeEndpoints(s service.Brackets) Endpoints {
	return Endpoints{
		Fix:      makeFixEndpoint(s),
		Validate: makeValidateEndpoint(s),
	}
}

func makeFixEndpoint(s service.Brackets) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*brackets.Brackets)
		return s.Fix(ctx, req)
	}
}

func makeValidateEndpoint(s service.Brackets) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*brackets.Brackets)
		return s.Validate(ctx, req)
	}
}
