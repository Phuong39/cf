package awsec2

import (
	"os"

	"github.com/teamssix/cf/pkg/util/errutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/util/cmdutil"
)

func EC2Client(region string) *ec2.EC2 {
	config := cmdutil.GetConfig("aws")
	if config.AccessKeyId == "" {
		log.Warnln("需要先配置访问密钥 (Access Key need to be configured first)")
		os.Exit(0)
		return nil
	} else {
		if region == "all" {
			region = "us-east-1"
		}
		cfg := &aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewStaticCredentials(config.AccessKeyId, config.AccessKeySecret, config.STSToken),
		}
		sess := session.Must(session.NewSession(cfg))
		svc := ec2.New(sess)
		log.Traceln("EC2 Client 连接成功 (EC2 Client connection successful)")
		return svc
	}
}

func GetEC2Regions() []*ec2.Region {
	svc := EC2Client("all")
	result, err := svc.DescribeRegions(nil)
	errutil.HandleErr(err)
	return result.Regions
}
