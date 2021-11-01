package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/bigmuramura/awsConfigToggle/mypkg"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

// onCmd represents the on command
var onCmd = &cobra.Command{
	Use:   "on",
	Short: "Enabled AWS Config",
	Long:  `Enabled the recorder status of AWS Config for all regions.`,

	Run: func(cmd *cobra.Command, args []string) {
		_, err := enbabledAWSConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if listFlag {
			_, err := mypkg.ShowConfigStatus()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	onCmd.Flags().BoolVarP(&listFlag, "verbose", "v", false, "Status display after execution.")
	rootCmd.AddCommand(onCmd)

}

func enbabledAWSConfig() (string, error) {

	const RECORDERNAME = "default"
	allRegions, err := mypkg.FetchAllRegions()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sess := session.Must(session.NewSession())

	// 並列処理を開始
	eg := errgroup.Group{}
	for _, region := range allRegions {
		region := region
		eg.Go(func() error {
			svc := configservice.New(
				sess,
				aws.NewConfig().WithRegion(region))

			input := &configservice.StartConfigurationRecorderInput{
				ConfigurationRecorderName: aws.String(RECORDERNAME),
			}
			_, err := svc.StartConfigurationRecorder(input)
			if err != nil {
				return err
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return "Failed", err
	}
	return "Suceed.", nil

}
