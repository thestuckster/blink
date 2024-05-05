package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/thestuckster/blink/internal"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all add ons installed with Blink",
	Long:  `example usage: blink list`,
	Run: func(cmd *cobra.Command, args []string) {
		listAddOns()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listAddOns() {
	config := internal.LoadConfig()
	for _, installed := range config.AddOns {
		fmt.Println("-----")
		fmt.Printf("%s\n %s\n Url: %s\n-----", installed.Repo, installed.Version, installed.Url)
	}
}
