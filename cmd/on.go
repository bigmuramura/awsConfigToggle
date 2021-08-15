package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/bigmuramura/awsConfigToggle/mypkg"
	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
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
	allRegions := mypkg.FetchAllRegions()
	sess := session.Must(session.NewSession())

	// Progress Bar
	count := len(allRegions)
	bar := pb.Simple.Start(count)
	bar.SetMaxWidth(80)

	for _, region := range allRegions {
		bar.Increment()
		svc := configservice.New(
			sess,
			aws.NewConfig().WithRegion(region))

		input := &configservice.StartConfigurationRecorderInput{
			ConfigurationRecorderName: aws.String(RECORDERNAME),
		}
		_, err := svc.StartConfigurationRecorder(input)
		if err != nil {
			return "Failed.", err
		}
	}
	bar.Finish()
	return "Suceed.", nil

}
