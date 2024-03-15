package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/tranHieuDev23/GoLoad/internal/configs"
	"github.com/tranHieuDev23/GoLoad/internal/wiring"
)

var (
	version    string
	commitHash string
)

const (
	flagConfigFilePath = "config-file-path"
)

func server() *cobra.Command {
	command := &cobra.Command{
		Use: "server",
		RunE: func(cmd *cobra.Command, args []string) error {
			configFilePath, err := cmd.Flags().GetString(flagConfigFilePath)
			if err != nil {
				return err
			}

			app, cleanup, err := wiring.InitializeServer(configs.ConfigFilePath(configFilePath))
			if err != nil {
				return err
			}

			defer cleanup()

			return app.Start()
		},
	}

	command.Flags().String(flagConfigFilePath, "", "If provided, will use the provided config file.")

	return command
}

func main() {
	rootCommand := &cobra.Command{
		Version: fmt.Sprintf("%s-%s", version, commitHash),
	}
	rootCommand.AddCommand(
		server(),
	)

	if err := rootCommand.Execute(); err != nil {
		log.Panic(err)
	}
}
