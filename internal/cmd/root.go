package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func New(name string, problem int, runE func(cmd *cobra.Command, args []string) error) *cobra.Command {
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
			log.Logger = log.With().Caller().Logger()

			zerolog.TimeFieldFormat = time.DateTime
			// set it similary to log.Lshortfile
			zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
				return filepath.Base(file) + ":" + strconv.Itoa(line)
			}

			ctx, _ := signal.NotifyContext(
				context.Background(),
				os.Interrupt,
				syscall.SIGTERM,
			)

			cmd.SetContext(ctx)
		},
	}

	rootCmd.PersistentFlags().BoolP("text", "t", false, "use text logging")
	rootCmd.PersistentFlags().StringP("log-level", "", "debug", "log level")
	rootCmd.PersistentFlags().StringP("addr", "a", ":8080", "listening address")
	viper.BindPFlags(rootCmd.PersistentFlags())
	viper.MustBindEnv("addr", "ADDRESS")
	viper.MustBindEnv("log-level", "LOG_LEVEL")
	return rootCmd
}
