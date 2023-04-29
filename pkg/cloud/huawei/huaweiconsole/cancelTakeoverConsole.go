package huaweiconsole

import (
	iamModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/util/database"
	"github.com/teamssix/cf/pkg/util/errutil"
)

func DeleteUser(userId string, userName string) {
	keystoneDeleteUserRequestContent := &iamModel.KeystoneDeleteUserRequest{}
	keystoneDeleteUserRequestContent.UserId = userId
	_, err := IAMClient().KeystoneDeleteUser(keystoneDeleteUserRequestContent)
	errutil.HandleErrNoExit(err)
	if err == nil {
		log.Debugf("删除 %s 用户成功 (Delete %s user successfully)", userName, userName)
	}
}

func CancelTakeoverConsole() {
	TakeoverConsoleCache := database.SelectTakeoverConsoleCache("huawei")
	if len(TakeoverConsoleCache) == 0 {
		log.Infoln("未接管过控制台，无需取消 (No takeover of the console, no need to cancel)")
	} else {
		userId := TakeoverConsoleCache[0].PrimaryAccountID
		userName := TakeoverConsoleCache[0].UserName
		DeleteUser(userId, userName)
		database.DeleteTakeoverConsoleCache("huawei")
		log.Infof("成功删除 %s 用户，已取消控制台接管 (Successful deletion of %s user, console takeover cancelled)", userName, userName)
	}
}
