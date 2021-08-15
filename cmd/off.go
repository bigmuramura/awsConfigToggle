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

// offCmd represents the off command
var offCmd = &cobra.Command{
	Use:   "off",
	Short: "Disabled AWS Config",
	Long:  `Disabled the recorder status of AWS Config for all regions.`,

	Run: func(cmd *cobra.Command, args []string) {
		res, err := disabledAWSConfig()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(offCmd)

}

func disabledAWSConfig() (string, error) {
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

		input := &configservice.StopConfigurationRecorderInput{
			ConfigurationRecorderName: aws.String(RECORDERNAME),
		}
		_, err := svc.StopConfigurationRecorder(input)
		if err != nil {
			return "Failed.", err
		}
	}
	bar.Finish()
	return "Suceed.", nil

}
