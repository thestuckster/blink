package cmd

import (
	"fmt"
	"github.com/thestuckster/blink/internal"
	"log"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("update called")
		singleFlag, err := cmd.Flags().GetBool("single")
		if err != nil {
			panic(err)
		}

		if singleFlag {
			updateSingle(args[0])
		} else {
			updateAll()
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	updateCmd.Flags().BoolP("single", "s", false, "Update a single specified add on")
}

func updateSingle(repo string) {
	config := internal.LoadConfig()
	_, addOnDetails := internal.FindAddOnDetails(repo, &config)

	if addOnDetails == nil {
		log.Println("No details found for add on: " + repo)
		return
	}

	log.Println("Backing up existing files...")
	backedUpFiles := backUpExistingFiles(addOnDetails)

	updated, err := internal.Update(addOnDetails.Url, addOnDetails, &config)
	errHandled := false
	if err != nil {
		log.Println("Error during update. Restoring from backup.")
		internal.RestoreBackUps(backedUpFiles)
		errHandled = true
	}

	if updated == false && errHandled == false {
		internal.RestoreBackUps(backedUpFiles)
	} else {
		log.Println(repo + " successfully updated.")
	}
}

func updateAll() {

}

func backUpExistingFiles(addOnDetails *internal.AddOn) []string {
	backedUpFiles := make([]string, 0)
	for _, file := range addOnDetails.Folders {
		newFileName := internal.BackUpExistingFile(file)
		backedUpFiles = append(backedUpFiles, newFileName)
	}

	return backedUpFiles
}
