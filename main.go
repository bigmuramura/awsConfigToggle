package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	// Region
	REGION := "ap-northeast-1"
	sess := session.Must(session.NewSession())
	svc := ec2.New(
		sess,
		aws.NewConfig().WithRegion(REGION))

	resultRegions, err := svc.DescribeRegions(nil)
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	var regionNames []string
	for _, regions := range resultRegions.Regions {
		regionName := *regions.RegionName
		regionNames = append(regionNames, regionName)
	}
	listStatus(regionNames)
	configOn()
}

func listStatus(regions []string) {
	for _, region := range regions {
		sess := session.Must(session.NewSession())
		svc := configservice.New(
			sess,
			aws.NewConfig().WithRegion(region))

		input := &configservice.DescribeConfigurationRecorderStatusInput{}

		result, err := svc.DescribeConfigurationRecorderStatus(input)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(region, *result.ConfigurationRecordersStatus[0].Recording)
	}
}

func configOn() {
	region := "ap-northeast-1"
	sess := session.Must(session.NewSession())
	svc := configservice.New(
		sess,
		aws.NewConfig().WithRegion(region))

	input := &configservice.StartConfigurationRecorderInput{
		ConfigurationRecorderName: aws.String("default"),
	}

	result, err := svc.StartConfigurationRecorder(input)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)

}
