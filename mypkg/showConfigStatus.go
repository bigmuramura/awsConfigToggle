package mypkg

import (
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/sync/errgroup"
)

type ConfigStatus struct {
	regionName    string
	recorderState bool
}

func ShowConfigStatus() (string, error) {
	allRegions, err := FetchAllRegions()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	configStatusList := make([]ConfigStatus, 0)

	sess := session.Must(session.NewSession())

	// 並列処理開始
	eg := errgroup.Group{}
	mutex := sync.Mutex{}
	for _, region := range allRegions {
		region := region

		eg.Go(func() error {
			svc := configservice.New(
				sess,
				aws.NewConfig().WithRegion(region))

			input := &configservice.DescribeConfigurationRecorderStatusInput{}

			result, err := svc.DescribeConfigurationRecorderStatus(input)
			if err != nil {
				return err
			}
			mutex.Lock()
			configStatusList = append(configStatusList, ConfigStatus{region, *result.ConfigurationRecordersStatus[0].Recording})
			mutex.Unlock()
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		fmt.Println(err)
		return "Failed.", err
	}

	// リージョン名でソート
	sort.SliceStable(configStatusList, func(i, j int) bool {
		return configStatusList[i].regionName < configStatusList[j].regionName
	})

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

	return "Suceed.", nil
}

func FetchAllRegions() ([]string, error) {
	REGION := "ap-northeast-1"
	sess := session.Must(session.NewSession())
	svc := ec2.New(
		sess,
		aws.NewConfig().WithRegion(REGION))

	var regionNames []string
	resultRegions, err := svc.DescribeRegions(nil)
	if err != nil {
		return regionNames, err
	}

	for _, regions := range resultRegions.Regions {
		regionName := *regions.RegionName
		regionNames = append(regionNames, regionName)
	}
	return regionNames, nil
}
