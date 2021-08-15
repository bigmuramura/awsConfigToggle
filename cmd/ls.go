package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/bigmuramura/awsConfigToggle/mypkg"
	"github.com/cheggaaa/pb/v3"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type ConfigStatus struct {
	regionName    string
	recorderState bool
}

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Shows status of AWS Config",
	Long:  `Shows the recorder status of AWS Config for all regions.`,

	Run: func(cmd *cobra.Command, args []string) {
		allRegions := mypkg.FetchAllRegions()
		configStatusList := make([]ConfigStatus, 0)

		// Progress Bar
		count := len(allRegions)
		bar := pb.Simple.Start(count)
		bar.SetMaxWidth(80)

		sess := session.Must(session.NewSession())
		for _, region := range allRegions {
			bar.Increment()
			svc := configservice.New(
				sess,
				aws.NewConfig().WithRegion(region))

			input := &configservice.DescribeConfigurationRecorderStatusInput{}

			result, err := svc.DescribeConfigurationRecorderStatus(input)
			if err != nil {
				fmt.Println(err)
			}
			configStatusList = append(configStatusList, ConfigStatus{region, *result.ConfigurationRecordersStatus[0].Recording})
		}

		// リージョン名でソート
		sort.SliceStable(configStatusList, func(i, j int) bool {
			return configStatusList[i].regionName < configStatusList[j].regionName
		})
		bar.Finish()

		// テーブル表示
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Region", "AWS Config"})
		for _, v := range configStatusList {
			state := "Unknown"
			if v.recorderState {
				state = color.HiBlueString("On")
			} else {
				state = color.HiRedString("Off")
			}
			table.Append([]string{v.regionName, state})
		}
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)

}
