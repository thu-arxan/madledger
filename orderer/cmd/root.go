package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:  "orderer",
		Long: "This is the cli of MadLedger orderer.",
	}
	homeDir, _ = os.Getwd()
	// ordererHome, _ = util.MakeFileAbs(".orderer", homeDir)
)

func init() {
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

// Execute exec the command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func setLog() error {
	// log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	return nil
}

func registerStop() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	return
}
