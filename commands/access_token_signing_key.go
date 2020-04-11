package commands

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/spf13/cobra"
)

var GenerateAccessTokenSigningKeyCommand = &cobra.Command{
	Use:   "generate-access-token-signing-key",
	Short: "generate an access token signing key",
	Run: func(cmd *cobra.Command, args []string) {

		key := make([]byte, 1024)
		_, err := rand.Read(key)
		if err != nil {
			panic(err)
		}

		fmt.Println(base64.URLEncoding.EncodeToString(key))

	},
}
