package cmd

import (
	"github.com/thestuckster/blink/internal"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Delete an addon",
	Long:  `example usage: blink remove WeakAuras/WeakAuras2`,
	Run: func(cmd *cobra.Command, args []string) {
		remove(args[0])
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func remove(addOnRepoName string) {
	config := internal.LoadConfig()
	deleteAddOnFiles(addOnRepoName, config)
	config.RemoveAnAddOn(addOnRepoName)
}

func deleteAddOnFiles(addOnRepoName string, config internal.Config) {
	for _, addOn := range config.AddOns {
		if addOn.Repo == addOnRepoName {
			for _, folder := range addOn.Folders {
				err := os.RemoveAll(folder)
				if err != nil {
					log.Panic(err)
				}
			}
		}
	}
}
