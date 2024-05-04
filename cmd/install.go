package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thestuckster/blink/internal"
	"log"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install an addon",
	Long:  `example usage go install https://github.com/Auctionator/Auctionator`,
	Run: func(cmd *cobra.Command, args []string) {
		install(args)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func install(args []string) {
	config := internal.LoadConfig()
	repoName := internal.SplitProjectNameFromUrl(args[0])
	log.Println(repoName)

	internal.FetchLatestRelease(repoName, &config)
}
