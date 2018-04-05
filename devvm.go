// Command devvm starts a stopped instance and prints its Public DNS name. This
// exists because it can be compiled and used without needing a slew of aws cli
// tools :)
package main // import "github.com/daaku/devvm"

import (
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	instanceID := flag.String("instance", "", "instance id")
	flag.Parse()
	session, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		panic(err)
	}
	svc := ec2.New(session)
	_, err = svc.StartInstances(&ec2.StartInstancesInput{
		InstanceIds: aws.StringSlice([]string{*instanceID}),
	})
	if err != nil {
		panic(err)
	}
	out, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: aws.StringSlice([]string{*instanceID}),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(*out.Reservations[0].Instances[0].PublicDnsName)
}
