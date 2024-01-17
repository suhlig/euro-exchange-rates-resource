// Package concourse provides a type-safe interface for Concourse resources.
package concourse

import (
	"context"
	"errors"
	"io"

	"github.com/go-playground/validator/v10"
)

type Resource[S any, V any, P any] interface {
	// Name returns the human-readable name of the resource
	Name() string

	// Check is invoked to detect new versions of the resource.
	//
	// It is given the configured source and current version, and must append new versions to the response slice, in
	// chronological order (oldest first, including the requested version if it's still valid).
	//
	// [Check]: https://concourse-ci.org/implementing-resource-types.html#resource-check
	Check(ctx context.Context, request CheckRequest[S, V], response *CheckResponse[V], log io.Writer) error

	// Get is invoked to fetch the resource and place it in the given directory.
	//
	// If is given the configured source, exact version of the resource to fetch, and configured parameters. It is
	// also passed the destination directory where the state of the resource is to be placed.
	//
	// The function must add the fetched version to the response, and may add metadata.
	//
	// [In]: https://concourse-ci.org/implementing-resource-types.html#resource-in
	Get(ctx context.Context, request GetRequest[S, V, P], response *Response[V], log io.Writer, destination string) error

	// Put is invoked to store the resource as it is given in the passed directory.
	//
	// If is given the configured source and configured parameters. The function must add the resulting version
	// to the response, and may add metadata.
	//
	// [Out]: https://concourse-ci.org/implementing-resource-types.html#resource-out
	Put(ctx context.Context, request PutRequest[S, P], response *Response[V], log io.Writer, source string) error
}

type CheckRequest[S any, V any] struct {
	Source  S `json:"source" validate:"required"`
	Version V `json:"version" validate:"omitempty"`
}

func (r CheckRequest[S, V]) Validate() error {
	return validator.New(validator.WithRequiredStructEnabled()).Struct(r)
}

type GetRequest[S any, V any, P any] struct {
	Source  S `json:"source" validate:"required"`
	Version V `json:"version" validate:"required"`
	Params  P `json:"params"`
}

func (r GetRequest[S, V, P]) Validate() error {
	return validator.New(validator.WithRequiredStructEnabled()).Struct(r)
}

func (r PutRequest[S, P]) Validate() error {
	return validator.New(validator.WithRequiredStructEnabled()).Struct(r)
}

type PutRequest[S any, P any] struct {
	Source S `json:"source" validate:"required"`
	Params P `json:"params" validate:"dive"`
}

type CheckResponse[V any] []V

func (r CheckResponse[V]) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	var validationErrors error

	for _, v := range r {
		err := validate.Struct(v)

		if err != nil {
			validationErrors = errors.Join(validationErrors, err)
		}
	}

	return validationErrors
}

type Response[V any] struct {
	Version  V               `json:"version" validate:"required"`
	Metadata []NameValuePair `json:"metadata,omitempty"` // TODO optional, but if given, name must not be empty
}

func (r Response[V]) Validate() error {
	return validator.New(validator.WithRequiredStructEnabled()).Struct(r)
}

type NameValuePair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
