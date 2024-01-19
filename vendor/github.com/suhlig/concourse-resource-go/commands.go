package concourse

import (
	"github.com/spf13/cobra"
)

// NewRootCommand returns a Cobra command with three subcommands (check, get and put), suitable for a Concourse resource.
// Each command corresponds to the respective /opt/resource/{check,in,out} scripts.
//
// Validation is performed on request and response. If you add struct tags to conrete Source (S) and Version (V) types,
// they will be checked, too. Check the [validator] package for details.
//
// [validator]: https://pkg.go.dev/github.com/go-playground/validator
func NewRootCommand[S any, V any, P any](resource Resource[S, V, P], name string) *cobra.Command {
	var rootCommand = &cobra.Command{
		SilenceUsage: true,
		Short:        name,
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
		RunE: func(cmd *cobra.Command, _ []string) error {
			return CheckWithValidation(cmd.Context(), resource, cmd.InOrStdin(), cmd.OutOrStdout(), cmd.ErrOrStderr())
		},
	}
}

func getCommand[S any, V any, P any](resource Resource[S, V, P]) *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Fetches the requested version of the resource and places its state in the input directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return GetWithValidation(cmd.Context(), resource, cmd.InOrStdin(), cmd.OutOrStdout(), cmd.ErrOrStderr(), args[0])
		},
	}
}

func putCommand[S any, V any, P any](resource Resource[S, V, P]) *cobra.Command {
	return &cobra.Command{
		Use:   "put",
		Short: "Puts a new version of the resource from the state in the output directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return PutWithValidation(cmd.Context(), resource, cmd.InOrStdin(), cmd.OutOrStdout(), cmd.ErrOrStderr(), args[0])
		},
	}
}
