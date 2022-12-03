package alioss

import (
	"fmt"
	"os"
	"strings"

	"github.com/teamssix/cf/pkg/util/errutil"

	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util/cmdutil"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OSSCollector struct {
	Conf   cloud.Config
	Client *oss.Client
}

func CreateOSSEndpoint(region string) string {
	return fmt.Sprintf("oss-%s.aliyuncs.com", region)
}

func (o *OSSCollector) OSSClient(region string) *OSSCollector {
	config := cmdutil.GetConfig("alibaba")
	if config.AccessKeyId == "" {
		log.Warnln("需要先配置访问密钥 (Access Key need to be configured first)")
		os.Exit(0)
		return nil
	} else {
		if config.STSToken == "" {
			client, err := oss.New(CreateOSSEndpoint(region), config.AccessKeyId, config.AccessKeySecret)
			errutil.HandleErr(err)
			if err == nil {
				log.Traceln("OSS Client 连接成功 (OSS Client connection successful)")
			}
			o.Client = client
		} else {
			client, err := oss.New(CreateOSSEndpoint(region), config.AccessKeyId, config.AccessKeySecret)
			errutil.HandleErr(err)
			if err == nil {
				log.Traceln("OSS Client 连接成功 (OSS Client connection successful)")
			}
			client.Config.SecurityToken = strings.TrimSpace(config.STSToken)
			o.Client = client
		}
		return o
	}
}
