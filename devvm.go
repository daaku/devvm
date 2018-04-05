// Command devvm starts a stopped instance and prints its Public DNS name. This
// exists because it can be compiled and used without needing a slew of aws cli
// tools :)
package main // import "github.com/daaku/devvm"

import (
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	instanceID := flag.String("instance", "", "instance id")
	port := flag.Int("port", 22, "port to wait for")
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
	var dnsName string
	for {
		out, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
			InstanceIds: aws.StringSlice([]string{*instanceID}),
		})
		if err != nil {
			panic(err)
		}
		dnsName = *out.Reservations[0].Instances[0].PublicDnsName
		if dnsName == "" {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		fmt.Println(dnsName)
		break
	}
	for {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", dnsName, *port))
		if err != nil {
			continue
		}
		conn.Close()
		break
	}
}
