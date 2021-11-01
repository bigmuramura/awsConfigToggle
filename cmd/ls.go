package cmd

import (
	"fmt"
	"os"

	"github.com/bigmuramura/awsConfigToggle/mypkg"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Shows status of AWS Config",
	Long:  `Shows the recorder status of AWS Config for all regions.`,

	Run: func(cmd *cobra.Command, args []string) {
		_, err := mypkg.ShowConfigStatus()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)

}
