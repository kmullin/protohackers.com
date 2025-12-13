package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
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
		Short: fmt.Sprintf("%v (# %v)", cases.Title(language.Und).String(name), problem),
		RunE:  runE,

		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			level, err := zerolog.ParseLevel(strings.ToLower(viper.GetString("log-level")))
			if err == nil {
				zerolog.SetGlobalLevel(level)
			}

			if viper.GetBool("text") {
				log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
			}
		},
	}

	rootCmd.SetContext(ctx)

	rootCmd.PersistentFlags().BoolP("text", "t", false, "use text logging")
	rootCmd.PersistentFlags().StringP("log-level", "", "debug", "log level")
	rootCmd.PersistentFlags().StringP("addr", "a", ":8080", "listening address")
	viper.BindPFlags(rootCmd.PersistentFlags())
	viper.MustBindEnv("addr", "ADDRESS")
	viper.MustBindEnv("log-level", "LOG_LEVEL")
	return rootCmd, stop
}
