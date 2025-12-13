package main

import (
	"os"

	"github.com/kmullin/protohackers.com/internal/cmd"
	"github.com/kmullin/protohackers.com/internal/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd, stop := cmd.New("line reversal", 7, func(cmd *cobra.Command, args []string) error {
		s := NewServer(cmd.Context(), log.Logger)
		server.UDP(s)
		return nil
	})
	defer stop()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
