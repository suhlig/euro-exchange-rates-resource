package concourse

import (
	"context"
	"io"
)

type Resource[S any, V any, P any] interface {
	Check(ctx context.Context, request Request[S], response *[]V, log io.Writer) error
	Get(ctx context.Context, request GetRequest[S, V, P], response *Response[V], log io.Writer, destination string) error
	Put(ctx context.Context, request PutRequest[S, P], response *Response[V], log io.Writer, source string) error
}

type Request[S any] struct {
	Source S `json:"source" validate:"required"`
}

type Response[V any] struct {
	Version  V               `json:"version"`
	Metadata []NameValuePair `json:"metadata,omitempty"`
}

type NameValuePair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type GetRequest[S any, V any, P any] struct {
	Request[S]
	Version V `json:"version"`
	Params  P `json:"params"`
}

type PutRequest[S any, P any] struct {
	Request[S]
	Params P `json:"params"`
}
