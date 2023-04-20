package huaweiiam

import (
	"os"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	ecs "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2"
	ecsModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/model"
	ecsRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/region"
	iam "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3"
	iamModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	iamRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/region"
	rds "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3"
	rdsModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3/model"
	rdsRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3/region"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/util/cmdutil"
)

func IAMClient() *iam.IamClient {
	huaweiConfig := cmdutil.GetConfig("huawei")
	if huaweiConfig.AccessKeyId == "" {
		log.Warnln("需要先配置访问密钥 (Access Key need to be configured first)")
		os.Exit(0)
		return nil
	} else {
		// 判断是否已经配置了STS Token
		if huaweiConfig.STSToken == "" {
			auth := global.NewCredentialsBuilder().
				WithAk(huaweiConfig.AccessKeyId).
				WithSk(huaweiConfig.AccessKeySecret).
				Build()

			client := iam.NewIamClient(
				iam.IamClientBuilder().
					WithRegion(iamRegion.ValueOf("cn-east-3")).
					WithCredential(auth).
					Build())

			showPermanentAccessKeyRequestContent := &iamModel.ShowPermanentAccessKeyRequest{}
			showPermanentAccessKeyRequestContent.AccessKey = huaweiConfig.AccessKeyId
			showPermanentAccessKeyRequestResponse, err := client.ShowPermanentAccessKey(showPermanentAccessKeyRequestContent)
			if err != nil {
				log.Traceln(err)
			} else if showPermanentAccessKeyRequestResponse.Credential.Status == "active" {
				log.Traceln("IAM Client 连接成功 (IAM Client connection successful)")
			}
			return client
		} else {
			auth := global.NewCredentialsBuilder().
				WithAk(huaweiConfig.AccessKeyId).
				WithSk(huaweiConfig.AccessKeySecret).
				WithSecurityToken(huaweiConfig.STSToken).
				Build()

			client := iam.NewIamClient(
				iam.IamClientBuilder().
					WithRegion(iamRegion.ValueOf("cn-east-3")).
					WithCredential(auth).
					Build())

			showPermanentAccessKeyRequestContent := &iamModel.ShowPermanentAccessKeyRequest{}
			showPermanentAccessKeyRequestContent.AccessKey = huaweiConfig.AccessKeyId
			showPermanentAccessKeyRequestResponse, err := client.ShowPermanentAccessKey(showPermanentAccessKeyRequestContent)
			if err != nil {
				log.Traceln(err)
			} else if showPermanentAccessKeyRequestResponse.Credential.Status == "active" {
				log.Traceln("IAM Client 连接成功 (IAM Client connection successful)")
			}
			return client
		}
	}
}

func ECSClient() *ecs.EcsClient {
	huaweiConfig := cmdutil.GetConfig("huawei")
	if huaweiConfig.AccessKeyId == "" {
		log.Warnln("需要先配置访问密钥 (Access Key need to be configured first)")
		os.Exit(0)
		return nil
	} else {
		// 判断是否已经配置了STS Token
		if huaweiConfig.STSToken == "" {
			auth := basic.NewCredentialsBuilder().
				WithAk(huaweiConfig.AccessKeyId).
				WithSk(huaweiConfig.AccessKeySecret).
				Build()

			client := ecs.NewEcsClient(
				ecs.EcsClientBuilder().
					WithRegion(ecsRegion.ValueOf("cn-east-3")).
					WithCredential(auth).
					Build())
			listServersDetailsRequestContent := &ecsModel.ListServersDetailsRequest{}
			_, err := client.ListServersDetails(listServersDetailsRequestContent)
			if err != nil {
				log.Traceln(err)
			} else {
				log.Traceln("ECS Client 连接成功 (ECS Client connection successful)")
			}
			return client
		} else {
			auth := basic.NewCredentialsBuilder().
				WithAk(huaweiConfig.AccessKeyId).
				WithSk(huaweiConfig.AccessKeySecret).
				WithSecurityToken(huaweiConfig.STSToken).
				Build()

			client := ecs.NewEcsClient(
				ecs.EcsClientBuilder().
					WithRegion(ecsRegion.ValueOf("cn-east-3")).
					WithCredential(auth).
					Build())
			listServersDetailsRequestContent := &ecsModel.ListServersDetailsRequest{}
			_, err := client.ListServersDetails(listServersDetailsRequestContent)
			if err != nil {
				log.Traceln(err)
			} else {
				log.Traceln("ECS Client 连接成功 (ECS Client connection successful)")
			}
			return client
		}
	}
}

func RDSClient() *rds.RdsClient {
	huaweiConfig := cmdutil.GetConfig("huawei")
	if huaweiConfig.AccessKeyId == "" {
		log.Warnln("需要先配置访问密钥 (Access Key need to be configured first)")
		os.Exit(0)
		return nil
	} else {
		// 判断是否已经配置了STS Token
		if huaweiConfig.STSToken == "" {
			auth := basic.NewCredentialsBuilder().
				WithAk(huaweiConfig.AccessKeyId).
				WithSk(huaweiConfig.AccessKeySecret).
				Build()

			client := rds.NewRdsClient(
				rds.RdsClientBuilder().
					WithRegion(rdsRegion.ValueOf("cn-east-3")).
					WithCredential(auth).
					Build())
			listInstancesRequestContent := &rdsModel.ListInstancesRequest{}
			_, err := client.ListInstances(listInstancesRequestContent)
			if err != nil {
				log.Traceln(err)
			} else {
				log.Traceln("RDS Client 连接成功 (RDS Client connection successful)")
			}
			return client
		} else {
			auth := basic.NewCredentialsBuilder().
				WithAk(huaweiConfig.AccessKeyId).
				WithSk(huaweiConfig.AccessKeySecret).
				WithSecurityToken(huaweiConfig.STSToken).
				Build()

			client := rds.NewRdsClient(
				rds.RdsClientBuilder().
					WithRegion(rdsRegion.ValueOf("cn-east-3")).
					WithCredential(auth).
					Build())
			listInstancesRequestContent := &rdsModel.ListInstancesRequest{}
			_, err := client.ListInstances(listInstancesRequestContent)
			if err != nil {
				log.Traceln(err)
			} else {
				log.Traceln("RDS Client 连接成功 (RDS Client connection successful)")
			}
			return client
		}
	}
}
