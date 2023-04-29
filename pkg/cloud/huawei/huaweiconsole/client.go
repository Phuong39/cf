package huaweiconsole

import (
	"errors"
	"fmt"
	"github.com/teamssix/cf/pkg/util/errutil"
	"os"
	"strings"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	iam "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3"
	iamModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	iamRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/region"
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

			// 捕获 iam.NewIamClient 的异常信息
			defer func() {
				r := recover()
				if errStr, ok := r.(string); ok {
					if strings.Contains(errStr, "Incorrect IAM authentication information") {
						errutil.HandleErr(errors.New(errStr))
					} else {
						fmt.Println(errStr)
						os.Exit(0)
					}
				}
			}()

			client := iam.NewIamClient(
				iam.IamClientBuilder().
					WithRegion(iamRegion.ValueOf("cn-east-3")).
					WithCredential(auth).
					Build())

			showPermanentAccessKeyRequestContent := &iamModel.ShowPermanentAccessKeyRequest{}
			showPermanentAccessKeyRequestContent.AccessKey = huaweiConfig.AccessKeyId
			showPermanentAccessKeyRequestResponse, err := client.ShowPermanentAccessKey(showPermanentAccessKeyRequestContent)
			if err != nil {
				errutil.HandleErr(err)
			} else if showPermanentAccessKeyRequestResponse.Credential.Status == "active" {
				log.Traceln("IAM Client 连接成功 (IAM Client connection successful)")
			}
			return client
		} else {
			// 使用 STS Token 连接
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
				errutil.HandleErr(err)
			} else if showPermanentAccessKeyRequestResponse.Credential.Status == "active" {
				log.Traceln("IAM Client 连接成功 (IAM Client connection successful)")
			}
			return client
		}
	}
}

//func ECSClient() *ecs.EcsClient {
//	huaweiConfig := cmdutil.GetConfig("huawei")
//	if huaweiConfig.AccessKeyId == "" {
//		log.Warnln("需要先配置访问密钥 (Access Key need to be configured first)")
//		os.Exit(0)
//		return nil
//	} else {
//		// 判断是否已经配置了STS Token
//		if huaweiConfig.STSToken == "" {
//			auth := basic.NewCredentialsBuilder().
//				WithAk(huaweiConfig.AccessKeyId).
//				WithSk(huaweiConfig.AccessKeySecret).
//				Build()
//
//			client := ecs.NewEcsClient(
//				ecs.EcsClientBuilder().
//					WithRegion(ecsRegion.ValueOf("cn-east-3")).
//					WithCredential(auth).
//					Build())
//			listServersDetailsRequestContent := &ecsModel.ListServersDetailsRequest{}
//			_, err := client.ListServersDetails(listServersDetailsRequestContent)
//			if err != nil {
//				log.Traceln(err)
//			} else {
//				log.Traceln("ECS Client 连接成功 (ECS Client connection successful)")
//			}
//			return client
//		} else {
//			auth := basic.NewCredentialsBuilder().
//				WithAk(huaweiConfig.AccessKeyId).
//				WithSk(huaweiConfig.AccessKeySecret).
//				WithSecurityToken(huaweiConfig.STSToken).
//				Build()
//
//			client := ecs.NewEcsClient(
//				ecs.EcsClientBuilder().
//					WithRegion(ecsRegion.ValueOf("cn-east-3")).
//					WithCredential(auth).
//					Build())
//			listServersDetailsRequestContent := &ecsModel.ListServersDetailsRequest{}
//			_, err := client.ListServersDetails(listServersDetailsRequestContent)
//			if err != nil {
//				log.Traceln(err)
//			} else {
//				log.Traceln("ECS Client 连接成功 (ECS Client connection successful)")
//			}
//			return client
//		}
//	}
//}
