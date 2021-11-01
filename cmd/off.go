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

// Flags
var listFlag bool

// offCmd represents the off command
var offCmd = &cobra.Command{
	Use:   "off",
	Short: "Disabled AWS Config",
	Long:  `Disabled the recorder status of AWS Config for all regions.`,

	Run: func(cmd *cobra.Command, args []string) {
		_, err := disabledAWSConfig()
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
	// フラグ設定
	offCmd.Flags().BoolVarP(&listFlag, "verbose", "v", false, "Status display after execution.")
	rootCmd.AddCommand(offCmd)

}

func disabledAWSConfig() (string, error) {
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

			input := &configservice.StopConfigurationRecorderInput{
				ConfigurationRecorderName: aws.String(RECORDERNAME),
			}
			_, err := svc.StopConfigurationRecorder(input)
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
