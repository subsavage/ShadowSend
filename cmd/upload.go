package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload [file path]",
	Short: "Encrypt and upload a file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide the path to the file.")
			return
		}
		path := args[0]
		fmt.Println("Uploading file:", path)
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}
