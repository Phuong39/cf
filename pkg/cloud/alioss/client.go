package alioss

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"githubu.com/teamssix/cf/pkg/cloud"
	"githubu.com/teamssix/cf/pkg/util"
	"githubu.com/teamssix/cf/pkg/util/cmdutil"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OSSCollector struct {
	Conf   cloud.Credential
	Client *oss.Client
}

func CreateOSSEndpoint(region string) string {
	return fmt.Sprintf("oss-%s.aliyuncs.com", region)
}

func (o *OSSCollector) OSSClient(region string) *OSSCollector {
	config := cmdutil.GetAliCredential()
	if config.AccessKeyId == "" {
		log.Warnln("需要先配置访问凭证 (Access Key need to be configured first)")
		os.Exit(0)
		return nil
	} else {
		if config.STSToken == "" {
			client, err := oss.New(CreateOSSEndpoint(region), config.AccessKeyId, config.AccessKeySecret)
			util.HandleErr(err)
			if err == nil {
				log.Traceln("OSS Client 连接成功 (OSS Client connection successful)")
			}
			o.Client = client
		} else {
			client, err := oss.New(CreateOSSEndpoint(region), config.AccessKeyId, config.AccessKeySecret)
			util.HandleErr(err)
			if err == nil {
				log.Traceln("OSS Client 连接成功 (OSS Client connection successful)")
			}
			client.Config.SecurityToken = strings.TrimSpace(config.STSToken)
			o.Client = client
		}
		return o
	}
}

func formatFileSize(fileSize int64) (size string) {
	if fileSize < 1024 {
		return fmt.Sprintf("%.2f B", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2f KB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2f MB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2f GB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2f TB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else {
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}
