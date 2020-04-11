package commands

import (
	"github.com/spf13/cobra"
	"os"
	. "s33d-backend/app"
	"s33d-backend/config"
)

func NewServeCommand(appFactory ApplicationFactoryType) *cobra.Command {

	return &cobra.Command{
		Use:   "serve",
		Short: "start the s33d backend",
		Run: func(cmd *cobra.Command, args []string) {

			conf, err := config.NewConfiguration(os.Getenv)
			if err != nil {
				panic(err)
			}

			a, err := appFactory(conf)
			if err != nil {
				panic(err)
			}

			if err := a.Start(); err != nil {
				panic(err)
			}

		},
	}
}
