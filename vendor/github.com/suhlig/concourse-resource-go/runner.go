package concourse

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

func CheckWithValidation[S any, V any, P any](ctx context.Context, resource Resource[S, V, P], stdin io.Reader, stdout, stderr io.Writer) error {
	var request CheckRequest[S, V]
	err := json.NewDecoder(stdin).Decode(&request)

	if err != nil {
		return fmt.Errorf("unable to decode request: %w", err)
	}

	err = request.Validate()

	if err != nil {
		return fmt.Errorf("validation of check request failed: %w", err)
	}

	response, err := resource.Check(ctx, request, stderr)

	if err != nil {
		return fmt.Errorf("check failed: %w", err)
	}

	err = response.Validate()

	if err != nil {
		return fmt.Errorf("validation of check response failed: %w", err)
	}

	return json.NewEncoder(stdout).Encode(response)
}

func GetWithValidation[S any, V any, P any](ctx context.Context, resource Resource[S, V, P], stdin io.Reader, stdout, stderr io.Writer, destination string) error {
	var request GetRequest[S, V, P]
	err := json.NewDecoder(stdin).Decode(&request)

	if err != nil {
		return fmt.Errorf("unable to decode request: %w", err)
	}

	err = request.Validate()

	if err != nil {
		return fmt.Errorf("validation of get request failed: %w", err)
	}

	response, err := resource.Get(ctx, request, stderr, destination)

	if err != nil {
		return fmt.Errorf("get failed: %w", err)
	}

	err = response.Validate()

	if err != nil {
		return fmt.Errorf("validation of get response failed: %w", err)
	}

	return json.NewEncoder(stdout).Encode(response)

}

func PutWithValidation[S any, V any, P any](ctx context.Context, resource Resource[S, V, P], stdin io.Reader, stdout, stderr io.Writer, source string) error {
	var request PutRequest[S, P]
	err := json.NewDecoder(stdin).Decode(&request)

	if err != nil {
		return fmt.Errorf("unable to decode request: %w", err)
	}

	err = request.Validate()

	if err != nil {
		return fmt.Errorf("validation of put request failed: %w", err)
	}

	response, err := resource.Put(ctx, request, stderr, source)

	if err != nil {
		return fmt.Errorf("put failed: %w", err)
	}

	err = response.Validate()

	if err != nil {
		return fmt.Errorf("validation of put response failed: %w", err)
	}

	return json.NewEncoder(stdout).Encode(response)
}
