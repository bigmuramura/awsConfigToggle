package mypkg

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

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
