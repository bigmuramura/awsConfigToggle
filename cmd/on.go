package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/bigmuramura/awsConfigToggle/mypkg"
	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

// onCmd represents the on command
var onCmd = &cobra.Command{
	Use:   "on",
	Short: "Enabled AWS Config",
	Long:  `Enabled the recorder status of AWS Config for all regions.`,

	Run: func(cmd *cobra.Command, args []string) {
		res, err := enbabledAWSConfig()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
	},
}

func init() {
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

	// Progress Bar
	count := len(allRegions)
	bar := pb.Simple.Start(count)
	bar.SetMaxWidth(80)

	// 並列処理を開始
	eg := errgroup.Group{}
	for _, region := range allRegions {
		region := region
		eg.Go(func() error {
			bar.Increment()
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
	bar.Finish()
	return "Suceed.", nil

}
