package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func New(name string, problem int, runE func(cmd *cobra.Command, args []string) error) (*cobra.Command, context.CancelFunc) {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	rootCmd := &cobra.Command{
		Use:   fmt.Sprintf("problem%v", problem),
		Short: fmt.Sprintf("%v (%v)", cases.Title(language.Und).String(name), problem),
		RunE:  runE,

		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if viper.GetBool("text") {
				log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
			}
		},
	}

	rootCmd.SetContext(ctx)

	rootCmd.PersistentFlags().BoolP("text", "t", false, "use text logging")
	viper.BindPFlags(rootCmd.PersistentFlags())
	return rootCmd, stop
}
