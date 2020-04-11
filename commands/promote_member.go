package commands

import (
	"errors"
	"github.com/spf13/cobra"
	vo "github.com/214alphadev/community-bl/value_objects"
	"os"
	"s33d-backend/app"
	"s33d-backend/config"
)

func NewPromoteMemberCommand(appFactory app.ApplicationFactoryType) *cobra.Command {

	return &cobra.Command{
		Use:   "member:promote",
		Short: "promote a member to admin",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing email address")
			}
			_, err := vo.NewEmailAddress(args[0])
			return err
		},
		Run: func(cmd *cobra.Command, args []string) {

			emailAddress, err := vo.NewEmailAddress(args[0])
			if err != nil {
				panic(err)
			}

			conf, err := config.NewConfiguration(os.Getenv)
			if err != nil {
				panic(err)
			}

			a, err := appFactory(conf)
			if err != nil {
				panic(err)
			}

			if err := a.Community.Promote(emailAddress); err != nil {
				panic(err)
			}

		},
	}
}
