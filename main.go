package main

import (
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"s33d-backend/app"
	"s33d-backend/commands"
)

func main() {

	_ = godotenv.Load()

	var rootCmd = &cobra.Command{}

	rootCmd.AddCommand(commands.NewServeCommand(app.ApplicationFactory))
	rootCmd.AddCommand(commands.GenerateAccessTokenSigningKeyCommand)
	rootCmd.AddCommand(commands.NewPromoteMemberCommand(app.ApplicationFactory))

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}

}
