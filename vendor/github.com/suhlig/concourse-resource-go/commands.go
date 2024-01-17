package concourse

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// NewRootCommand returns a Cobra command with three subcommands (check, get and put), suitable for a Concourse resource.
// Each command corresponds to the respective /opt/resource/{check,in,out} scripts.
//
// Validation is performed on request and response. If you add struct tags to conrete Source (S) and Version (V) types,
// they will be checked, too. Check the [validator] package for details.
//
// [validator]: https://pkg.go.dev/github.com/go-playground/validator
func NewRootCommand[S any, V any, P any](resource Resource[S, V, P]) *cobra.Command {
	var rootCommand = &cobra.Command{
		SilenceUsage: true,
		Short:        resource.Name(),
	}

	rootCommand.AddCommand(checkCommand(resource))
	rootCommand.AddCommand(getCommand(resource))
	rootCommand.AddCommand(putCommand(resource))

	return rootCommand
}

func checkCommand[S any, V any, P any](resource Resource[S, V, P]) *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "Fetches the latest version of the resource and emit its version",
		RunE: func(cmd *cobra.Command, args []string) error {
			var request CheckRequest[S, V]
			err := json.NewDecoder(cmd.InOrStdin()).Decode(&request)

			if err != nil {
				return fmt.Errorf("unable to decode request: %w", err)
			}

			err = request.Validate()

			if err != nil {
				return fmt.Errorf("request is invalid: %w", err)
			}

			var response CheckResponse[V]
			err = resource.Check(cmd.Context(), request, &response, cmd.ErrOrStderr())

			if err != nil {
				return fmt.Errorf("check failed: %w", err)
			}

			err = response.Validate()

			if err != nil {
				return fmt.Errorf("response is invalid: %w", err)
			}

			return json.NewEncoder(cmd.OutOrStdout()).Encode(response)
		},
	}
}

func getCommand[S any, V any, P any](resource Resource[S, V, P]) *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Fetches the requested version of the resource and places its state in the input directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var request GetRequest[S, V, P]
			err := json.NewDecoder(cmd.InOrStdin()).Decode(&request)

			if err != nil {
				return fmt.Errorf("unable to decode request: %w", err)
			}

			err = request.Validate()

			if err != nil {
				return fmt.Errorf("request is invalid: %w", err)
			}

			var response Response[V]
			err = resource.Get(cmd.Context(), request, &response, cmd.ErrOrStderr(), args[0])

			if err != nil {
				return fmt.Errorf("get failed: %w", err)
			}

			err = response.Validate()

			if err != nil {
				return fmt.Errorf("response is invalid: %w", err)
			}

			return json.NewEncoder(cmd.OutOrStdout()).Encode(response)
		},
	}
}

func putCommand[S any, V any, P any](resource Resource[S, V, P]) *cobra.Command {
	return &cobra.Command{
		Use:   "put",
		Short: "Puts a new version of the resource from the state in the output directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var request PutRequest[S, P]
			err := json.NewDecoder(cmd.InOrStdin()).Decode(&request)

			if err != nil {
				return fmt.Errorf("unable to decode request: %w", err)
			}

			err = request.Validate()

			if err != nil {
				return fmt.Errorf("request is invalid: %w", err)
			}

			var response Response[V]

			err = resource.Put(cmd.Context(), request, &response, cmd.ErrOrStderr(), args[0])

			if err != nil {
				return fmt.Errorf("put failed: %w", err)
			}

			err = response.Validate()

			if err != nil {
				return fmt.Errorf("response is invalid: %w", err)
			}

			return json.NewEncoder(cmd.OutOrStdout()).Encode(response)
		},
	}
}
