package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

func main() {
	fmt.Println("test")

	REGION := "ap-northeast-1"
	sess := session.Must(session.NewSession())
	svc := costexplorer.New(
		sess,
		aws.NewConfig().WithRegion(REGION))

}
