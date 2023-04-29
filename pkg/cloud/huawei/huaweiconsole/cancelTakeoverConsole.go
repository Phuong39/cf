package huaweiconsole

import (
	iamModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/util/database"
)

func DeleteUser(userId string) {
	keystoneDeleteUserRequestContent := &iamModel.KeystoneDeleteUserRequest{}
	keystoneDeleteUserRequestContent.UserId = userId
	keystoneDeleteUserRequestResponse, err := IAMClient().KeystoneDeleteUser(keystoneDeleteUserRequestContent)
	if err == nil {
		log.Debugln(keystoneDeleteUserRequestResponse)
	} else {
		log.Traceln(err)
	}
}

func CancelTakeoverConsole(userName string) {
	TakeoverConsoleCache := database.SelectTakeoverConsoleCache("huawei")
	if len(TakeoverConsoleCache) == 0 {
		log.Infoln("未接管过控制台，无需取消 (No takeover of the console, no need to cancel)")
	} else {
		userId := TakeoverConsoleCache[0].PrimaryAccountID
		DeleteUser(userId)
	}
}
